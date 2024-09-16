package config

import (
	"log"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
)

type Config struct {
	Env            string `yaml:"env" env-required:"true"`
	Kafka          `yaml:"kafka" env-required:"true"`
	FirstConsumer  Consumer `yaml:"consumer_1" env-required:"true"`
	SecondConsumer Consumer `yaml:"consumer_2" env-required:"true"`
	HTTPServer     `yaml:"http_server" env-required:"true"`
}

type Kafka struct {
	KafkaBrokers         string `yaml:"kafka_brokers" env-reqiered:"true"`
	ProducerAsks         string `yaml:"producer_asks" env-required:"true"`
	ProduceRetries       int    `yaml:"produce_retries" env-reqiered:"true"`
	MessageTopic         string `yaml:"message_topic" env-reqiured:"true"`
	MessageConsumerGroup string `yaml:"message_consumer_group" env-reqiered:"true"`
}

type Consumer struct {
	KafkaBrokers  string `yaml:"kafka_brokers" env-reqiered:"true"`
	ConsumerGroup string `yaml:"consumer_group" env-required:"true"`
	MessageTopic  string `yaml:"message_topic" env-reqiured:"true"`
}

type HTTPServer struct {
	Host        string        `yaml:"host" env-required:"true"`
	Port        int           `yaml:"port" env-required:"true"`
	Protocol    string        `yaml:"protocol" env-required:"true"`
	Timeout     time.Duration `yaml:"timeout" env-required:"true"`
	IdleTimeout time.Duration `yaml:"idle_timeout" env-required:"true"`
}

func MustLoad() *Config {
	if err := godotenv.Load(".env"); err != nil {
		log.Fatal("failed to load environment file, error: ", err)
	}

	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		log.Fatal("CONFIG_PATH is not set")
	}

	cfg := MustLoadWithPath(configPath)
	return cfg
}

func MustLoadWithPath(configPath string) *Config {
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatal("config file not found: ", err)
	}

	var cfg Config

	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		log.Fatal("failed to read config, error: ", err)
	}

	return &cfg
}
