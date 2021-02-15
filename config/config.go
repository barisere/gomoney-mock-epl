package config

import (
	"os"
	"strconv"
	"strings"
)

type Config struct {
	HttpBindPort uint
	MongoURL     string
}

func LoadConfig() (*Config, error) {
	port := strings.TrimSpace(os.Getenv("PORT"))
	mongoURL := strings.TrimSpace(os.Getenv("MONGO_URL"))

	var httpPort uint = 8080
	if port != "" {
		if p, err := strconv.ParseUint(port, 10, 32); err != nil {
			return nil, err
		} else {
			httpPort = uint(p)
		}
	}

	if mongoURL == "" {
		mongoURL = "mongodb://localhost:27017/hf?ssl=false"
	}

	return &Config{
		HttpBindPort: httpPort,
		MongoURL:     mongoURL,
	}, nil
}
