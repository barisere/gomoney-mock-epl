package config

import (
	"os"
	"strconv"
	"strings"
)

type Environment uint8

const (
	Development Environment = iota
	Testing
	Production
)

var envMap = map[string]Environment{
	"development": Development,
	"testing":     Testing,
	"production":  Production,
}

type Config struct {
	Environment
	HttpBindPort uint64
	MongoURL     string
}

func LoadConfig() (*Config, error) {
	env := strings.TrimSpace(os.Getenv("DEPLOY_ENV"))
	port := strings.TrimSpace(os.Getenv("PORT"))
	mongoURL := strings.TrimSpace(os.Getenv("MONGO_URL"))

	environment, ok := envMap[env]
	if !ok {
		environment = Development
	}

	httpPort, err := strconv.ParseUint(port, 10, 32)
	if err != nil {
		return nil, err
	}

	return &Config{
		Environment:  environment,
		HttpBindPort: httpPort,
		MongoURL:     mongoURL,
	}, nil
}
