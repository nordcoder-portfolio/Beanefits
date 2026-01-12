package httpserver

import "time"

type Config struct {
	Addr            string
	BaseURL         string
	ReadTimeout     time.Duration
	WriteTimeout    time.Duration
	IdleTimeout     time.Duration
	ShutdownTimeout time.Duration
	RequestTimeout  time.Duration
}

func DefaultConfig() Config {
	return Config{
		Addr:            ":8080",
		BaseURL:         "",
		ReadTimeout:     5 * time.Second,
		WriteTimeout:    10 * time.Second,
		IdleTimeout:     60 * time.Second,
		ShutdownTimeout: 10 * time.Second,
		RequestTimeout:  15 * time.Second,
	}
}
