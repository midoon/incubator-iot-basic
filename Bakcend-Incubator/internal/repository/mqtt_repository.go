package repository

import (
	"backend-incubator/pkg/config"
	"encoding/json"
	"fmt"
	"log"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type mqttRepository struct {
	client mqtt.Client
	cfg    *config.Config
}

func NewMQTTRepository(client mqtt.Client, cfg *config.Config) *mqttRepository {
	return &mqttRepository{
		client: client,
		cfg:    cfg,
	}
}

// PublishMode mem-publish perintah ganti mode ke MQTT broker.
func (r *mqttRepository) PublishMode(mode string) error {
	payload, err := json.Marshal(map[string]string{"mode": mode})
	if err != nil {
		return fmt.Errorf("marshal mode payload: %w", err)
	}

	token := r.client.Publish(r.cfg.TopicCmdMode, 1, false, payload)
	if token.WaitTimeout(5*time.Second) && token.Error() != nil {
		return fmt.Errorf("publish mode: %w", token.Error())
	}

	log.Printf("[MQTT] Published mode=%s to %s\n", mode, r.cfg.TopicCmdMode)
	return nil
}

// PublishLamp mem-publish perintah kontrol lampu ke MQTT broker.
func (r *mqttRepository) PublishLamp(on bool) error {
	payload, err := json.Marshal(map[string]bool{"lamp": on})
	if err != nil {
		return fmt.Errorf("marshal lamp payload: %w", err)
	}

	token := r.client.Publish(r.cfg.TopicCmdLamp, 1, false, payload)
	if token.WaitTimeout(5*time.Second) && token.Error() != nil {
		return fmt.Errorf("publish lamp: %w", token.Error())
	}

	log.Printf("[MQTT] Published lamp=%v to %s\n", on, r.cfg.TopicCmdLamp)
	return nil
}
