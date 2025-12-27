package config

import (
	"os"
	"strconv"
	"strings"
	"time"
)

type Config struct {
	HTTPAddr        string
	HTTPBaseURL     string
	ShutdownTimeout time.Duration

	PostgresURL string

	JWTSecret string
	JWTIssuer string
	JWTTTL    time.Duration

	KafkaBrokers         []string
	KafkaTopic           string
	KafkaWriteTimeout    time.Duration
	KafkaBatchSize       int
	KafkaBatchTimeout    time.Duration
	KafkaAutoPublish     bool
	KafkaPublishInterval time.Duration
}

func MustLoad() Config {
	return Config{
		HTTPAddr:        getenv("HTTP_ADDR", ":8080"),
		HTTPBaseURL:     getenv("HTTP_BASE_URL", ""),
		ShutdownTimeout: mustDuration(getenv("SHUTDOWN_TIMEOUT", "10s")),

		PostgresURL: getenv("POSTGRES_URL", "postgres://postgres:postgres@localhost:5432/beanefits?sslmode=disable"),

		JWTSecret: getenv("JWT_SECRET", "dev-secret-change-me"),
		JWTIssuer: getenv("JWT_ISSUER", "beanefits"),
		JWTTTL:    mustDuration(getenv("JWT_TTL", "24h")),

		KafkaBrokers:         mustCSVStrings(getenv("KAFKA_BROKERS", "localhost:9092")),
		KafkaTopic:           getenv("KAFKA_TOPIC", "ints"),
		KafkaWriteTimeout:    mustDuration(getenv("KAFKA_WRITE_TIMEOUT", "3s")),
		KafkaBatchSize:       mustInt(getenv("KAFKA_BATCH_SIZE", "1")),
		KafkaBatchTimeout:    mustDuration(getenv("KAFKA_BATCH_TIMEOUT", "0s")),
		KafkaAutoPublish:     mustBool(getenv("KAFKA_AUTOPUBLISH", "false")),
		KafkaPublishInterval: mustDuration(getenv("KAFKA_PUBLISH_INTERVAL", "2s")),
	}
}

func getenv(k, def string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return def
}

func mustDuration(s string) time.Duration {
	d, err := time.ParseDuration(s)
	if err != nil {
		panic("invalid duration: " + s)
	}
	return d
}

func mustInt(s string) int {
	v, err := strconv.Atoi(s)
	if err != nil {
		panic("invalid int: " + s)
	}
	return v
}

func mustCSVStrings(s string) []string {
	parts := strings.Split(s, ",")
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p == "" {
			continue
		}
		out = append(out, p)
	}
	if len(out) == 0 {
		panic("invalid csv string list: empty")
	}
	return out
}

func mustBool(s string) bool {
	v, err := strconv.ParseBool(s)
	if err != nil {
		panic("invalid bool: " + s)
	}
	return v
}
