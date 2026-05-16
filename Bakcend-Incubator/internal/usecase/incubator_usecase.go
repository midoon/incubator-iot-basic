package usecase

import (
	"backend-incubator/internal/domain"
	"fmt"
	"sync"
)

type incubatorUsecase struct {
	mu    sync.RWMutex
	state domain.IncubatorState
	repo  domain.IncubatorRepository
}

func NewIncubatorUsecase(repo domain.IncubatorRepository) domain.IncubatorUsecase {
	return &incubatorUsecase{
		state: domain.IncubatorState{
			Temperature: 0,
			Humidity:    0,
			Mode:        "auto",
			Lamp:        false,
		},
	}
}

// GetState mengembalikan state terkini (thread-safe).
func (u *incubatorUsecase) GetState() domain.IncubatorState {
	u.mu.RLock()
	defer u.mu.RUnlock()
	return u.state
}

// UpdateState memperbarui state dari data MQTT yang masuk (thread-safe).
func (u *incubatorUsecase) UpdateState(state domain.IncubatorState) {
	u.mu.Lock()
	defer u.mu.Unlock()
	u.state = state
}

// SetMode memavalidasi dan mengirim command ganti mode ke broker
func (u *incubatorUsecase) SetMode(mode string) error {
	if mode != "auto" && mode != "manual" {
		return fmt.Errorf("invalid mode: %s, must be 'auto' or 'manual'", mode)
	}

	if err := u.repo.PublishMode(mode); err != nil {
		return fmt.Errorf("failed to publish mode: %w", err)
	}

	u.mu.Lock()
	u.state.Mode = mode

	// Jika kembali ke auto, matikan lampu secara lokal.
	if mode == "auto" {
		u.state.Lamp = false
	}

	u.mu.Unlock()
	return nil
}

// SetLamp memvalidasi dan mengirim command kontrol lampu ke broker.
func (u *incubatorUsecase) SetLamp(on bool) error {
	u.mu.RLock()
	currentMode := u.state.Mode
	u.mu.RUnlock()

	if currentMode != "manual" {
		return fmt.Errorf("lamp can only be controlled in manual mode")
	}

	if err := u.repo.PublishLamp(on); err != nil {
		return fmt.Errorf("set lamp: %w", err)
	}

	u.mu.Lock()
	u.state.Lamp = on
	u.mu.Unlock()

	return nil
}
