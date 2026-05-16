package http

import (
	"backend-incubator/internal/domain"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"
)

type sseClient struct {
	ch chan domain.IncubatorState
}

type Handler struct {
	usecase domain.IncubatorUsecase
	mu      sync.Mutex
	clients map[*sseClient]struct{}
}

func NewHandler(uc domain.IncubatorUsecase) *Handler {
	return &Handler{
		usecase: uc,
		clients: make(map[*sseClient]struct{}),
	}
}

// BroadcastState dipanggil setiap kali ada data MQTT baru masuk.
// Mengirim state terbaru ke semua SSE client yang aktif.
func (h *Handler) BroadcastState(state domain.IncubatorState) {
	h.mu.Lock()
	defer h.mu.Unlock()

	for client := range h.clients {
		// Non-blocking send: jika channel penuh, skip client ini.
		select {
		case client.ch <- state:
		default:
		}
	}
}

// HandleSSE adalah endpoint SSE — frontend subscribe ke sini.
// GET /api/stream
func (h *Handler) HandleSSE(w http.ResponseWriter, r *http.Request) {
	// Set header SSE.
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "streaming not supported", http.StatusInternalServerError)
		return
	}

	// Daftarkan client baru.
	client := &sseClient{ch: make(chan domain.IncubatorState, 10)}
	h.mu.Lock()
	h.clients[client] = struct{}{}
	h.mu.Unlock()

	log.Printf("[SSE] Client connected. Total: %d\n", len(h.clients))

	// Kirim state awal langsung saat client connect.
	initialState := h.usecase.GetState()
	sendSSEEvent(w, flusher, initialState)

	// Ticker untuk heartbeat agar koneksi tidak timeout.
	ticker := time.NewTicker(15 * time.Second)
	defer ticker.Stop()

	defer func() {
		h.mu.Lock()
		delete(h.clients, client)
		h.mu.Unlock()
		log.Printf("[SSE] Client disconnected. Total: %d\n", len(h.clients))
	}()

	for {
		select {
		case <-r.Context().Done():
			// Frontend disconnect atau request dibatalkan.
			return

		case state := <-client.ch:
			sendSSEEvent(w, flusher, state)

		case <-ticker.C:
			// Heartbeat — kirim comment agar koneksi tetap hidup.
			fmt.Fprintf(w, ": heartbeat\n\n")
			flusher.Flush()
		}
	}
}

// sendSSEEvent memformat dan mengirim satu event SSE ke client.
func sendSSEEvent(w http.ResponseWriter, flusher http.Flusher, state domain.IncubatorState) {
	data, err := json.Marshal(state)
	if err != nil {
		log.Println("[SSE] Failed to marshal state:", err)
		return
	}
	fmt.Fprintf(w, "data: %s\n\n", data)
	flusher.Flush()
}

// HandleGetState mengembalikan state terkini sebagai JSON.
// GET /api/state
func (h *Handler) HandleGetState(w http.ResponseWriter, r *http.Request) {
	state := h.usecase.GetState()
	writeJSON(w, http.StatusOK, state)
}

// HandleSetMode menerima perintah ganti mode dari frontend.
// POST /api/mode
func (h *Handler) HandleSetMode(w http.ResponseWriter, r *http.Request) {
	var cmd domain.CommandMode
	if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if err := h.usecase.SetMode(cmd.Mode); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	// Broadcast state terbaru ke semua SSE client.
	h.BroadcastState(h.usecase.GetState())

	writeJSON(w, http.StatusOK, map[string]string{"message": "mode updated", "mode": cmd.Mode})
}

// HandleSetLamp menerima perintah kontrol lampu dari frontend.
// POST /api/lamp
func (h *Handler) HandleSetLamp(w http.ResponseWriter, r *http.Request) {
	var cmd domain.CommandLamp
	if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if err := h.usecase.SetLamp(cmd.Lamp); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	// Broadcast state terbaru ke semua SSE client.
	h.BroadcastState(h.usecase.GetState())

	writeJSON(w, http.StatusOK, map[string]interface{}{"message": "lamp updated", "lamp": cmd.Lamp})
}

func writeJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func writeError(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(map[string]string{"error": message})
}
