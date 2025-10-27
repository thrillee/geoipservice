package main

import (
	"log"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	ServerPort    string `envconfig:"PORT" default:"8080"`
	RedisAddr     string `envconfig:"REDIS_ADDR" default:"localhost:6379"`
	RedisPassword string `envconfig:"REDIS_PASSWORD" default:""`
	RedisDB       int    `envconfig:"REDIS_DB" default:"0"`

	// AWS Configuration
	AWSRegion          string `envconfig:"AWS_REGION" required:"true"`
	AWSAccessKeyID     string `envconfig:"AWS_ACCESS_KEY_ID" required:"true"`
	AWSSecretAccessKey string `envconfig:"AWS_SECRET_ACCESS_KEY" required:"true"`
	S3Bucket           string `envconfig:"S3_BUCKET" required:"true"`
	S3Key              string `envconfig:"S3_KEY" required:"true"`

	// Local cache for the database file
	LocalDBPath string `envconfig:"LOCAL_DB_PATH" default:"./geoip.mmdb"`
}

func LoadConfig() (*Config, error) {
	// Attempt to load .env file, but don't treat its absence as an error.
	_ = godotenv.Load()

	var cfg Config
	err := envconfig.Process("", &cfg)
	if err != nil {
		return nil, err
	}

	log.Println("Configuration loaded successfully")
	return &cfg, nil
}
