package repository

import (
	httphandler "backend-incubator/internal/delivery/http"
	"backend-incubator/internal/domain"
	"backend-incubator/pkg/config"
	"encoding/json"
	"log"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

// esp32Payload adalah struktur JSON yang dikirim ESP32.
type esp32Payload struct {
	Temperature *float64 `json:"temperature,omitempty"`
	Humidity    *float64 `json:"humidity,omitempty"`
	Mode        *string  `json:"mode,omitempty"`
	Lamp        *bool    `json:"lamp,omitempty"`
}

// StartSubscriber mendaftarkan semua subscription MQTT dan
// meneruskan data ke usecase + broadcast ke SSE clients.

func StartSubscriber(
	client mqtt.Client,
	cfg *config.Config,
	uc domain.IncubatorUsecase,
	handler *httphandler.Handler,
) {
	// Handler yang dipanggil tiap ada pesan MQTT masuk.
	messageHandler := func(c mqtt.Client, msg mqtt.Message) {
		topic := msg.Topic()
		payload := msg.Payload()

		log.Printf("[MQTT] Received topic=%s payload=%s\n", topic, string(payload))

		var p esp32Payload
		if err := json.Unmarshal(payload, &p); err != nil {
			log.Printf("[MQTT] Failed to unmarshal payload on topic %s: %v\n", topic, err)
			return
		}

		// Ambil state saat ini, update field yang relevan saja.
		state := uc.GetState()

		switch topic {
		case cfg.TopicTemperature:
			if p.Temperature != nil {
				state.Temperature = *p.Temperature
			}
		case cfg.TopicHumidity:
			if p.Humidity != nil {
				state.Humidity = *p.Humidity
			}
		case cfg.TopicMode:
			if p.Mode != nil {
				state.Mode = *p.Mode
			}
		case cfg.TopicLamp:
			if p.Lamp != nil {
				state.Lamp = *p.Lamp
			}
		default:
			// Jika ESP32 mengirim semua field dalam satu topic.
			if p.Temperature != nil {
				state.Temperature = *p.Temperature
			}
			if p.Humidity != nil {
				state.Humidity = *p.Humidity
			}
			if p.Mode != nil {
				state.Mode = *p.Mode
			}
			if p.Lamp != nil {
				state.Lamp = *p.Lamp
			}
		}

		// Simpan state baru di memory.
		uc.UpdateState(state)

		// Broadcast ke semua SSE client yang terhubung.
		handler.BroadcastState(state)
	}

	topics := []string{
		cfg.TopicTemperature,
		cfg.TopicHumidity,
		cfg.TopicMode,
		cfg.TopicLamp,
	}

	for _, topic := range topics {
		token := client.Subscribe(topic, 1, messageHandler)
		token.Wait()
		if token.Error() != nil {
			log.Printf("[MQTT] Failed to subscribe to %s: %v\n", topic, token.Error())
		} else {
			log.Printf("[MQTT] Subscribed to topic: %s\n", topic)
		}
	}
}
