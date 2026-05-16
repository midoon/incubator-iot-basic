package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	AppPort          string
	MQTTBroker       string
	MQTTClientID     string
	TopicTemperature string
	TopicHumidity    string
	TopicMode        string
	TopicLamp        string
	TopicCmdMode     string
	TopicCmdLamp     string
}

func Load() *Config {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	return &Config{
		AppPort:          getEnv("APP_PORT", "8080"),
		MQTTBroker:       getEnv("MQTT_BROKER", "mqtt://localhost:1883"),
		MQTTClientID:     getEnv("MQTT_CLIENT_ID", "incubator-backend"),
		TopicTemperature: getEnv("MQTT_TOPIC_TEMPERATURE", "incubator/temperature"),
		TopicHumidity:    getEnv("MQTT_TOPIC_HUMIDITY", "incubator/humidity"),
		TopicMode:        getEnv("MQTT_TOPIC_MODE", "incubator/mode"),
		TopicLamp:        getEnv("MQTT_TOPIC_LAMP", "incubator/lamp"),
		TopicCmdMode:     getEnv("MQTT_TOPIC_CMD_MODE", "incubator/cmd/mode"),
		TopicCmdLamp:     getEnv("MQTT_TOPIC_CMD_LAMP", "incubator/cmd/lamp"),
	}
}

func getEnv(key, fallback string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return fallback
}
