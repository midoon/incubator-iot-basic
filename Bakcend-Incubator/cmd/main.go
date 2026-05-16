package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"backend-incubator/internal/delivery"
	httphandler "backend-incubator/internal/delivery/http"
	"backend-incubator/internal/infrastructure"
	"backend-incubator/internal/repository"
	"backend-incubator/internal/usecase"
	"backend-incubator/pkg/config"
)

func main() {
	// Load konfigurasi dari .env.
	cfg := config.Load()

	// Inisialisasi koneksi MQTT.
	mqttClient, err := infrastructure.NewMQTTClient(cfg)
	if err != nil {
		log.Fatalf("[MAIN] Failed to connect MQTT: %v", err)
	}
	defer mqttClient.Disconnect(250)

	// Inisialisasi layer-layer clean architecture.
	repo := repository.NewMQTTRepository(mqttClient, cfg)
	uc := usecase.NewIncubatorUsecase(repo)
	handler := httphandler.NewHandler(uc)

	// Mulai subscribe ke MQTT broker (data dari ESP32).
	repository.StartSubscriber(mqttClient, cfg, uc, handler)

	// Setup HTTP server.
	router := delivery.NewRouter(handler)
	server := &http.Server{
		Addr:        ":" + cfg.AppPort,
		Handler:     router,
		ReadTimeout: 10 * time.Second,
		// WriteTimeout sengaja 0 agar SSE tidak di-timeout oleh server.
		WriteTimeout: 0,
		IdleTimeout:  120 * time.Second,
	}

	// Jalankan server di goroutine terpisah.
	go func() {
		log.Printf("[MAIN] Server running on port %s\n", cfg.AppPort)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("[MAIN] Server error: %v", err)
		}
	}()

	// Graceful shutdown — tunggu sinyal OS (Ctrl+C atau SIGTERM dari Docker).
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("[MAIN] Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("[MAIN] Server forced to shutdown: %v", err)
	}

	log.Println("[MAIN] Server stopped.")
}
