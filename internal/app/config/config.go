package config

import (
	"errors"
	"os"
	"strconv"
	"time"
)

type Config struct {
	HTTPAddr     string
	PostgresDSN  string
	InputDir     string
	OutputDir    string
	PollInterval time.Duration
	Workers      int
}

func Load() (Config, error) {
	var c Config

	c.HTTPAddr = envOr("HTTP_ADDR", ":8080")
	c.PostgresDSN = os.Getenv("POSTGRES_DSN")
	if c.PostgresDSN == "" {
		return c, errors.New("POSTGRES_DSN is required")
	}

	c.InputDir = os.Getenv("INPUT_DIR")
	if c.InputDir == "" {
		return c, errors.New("INPUT_DIR is required")
	}

	c.OutputDir = os.Getenv("OUTPUT_DIR")
	if c.OutputDir == "" {
		return c, errors.New("OUTPUT_DIR is required")
	}

	pollSec := envOr("POLL_INTERVAL_SEC", "5")
	n, err := strconv.Atoi(pollSec)
	if err != nil || n <= 0 {
		n = 5
	}
	c.PollInterval = time.Duration(n) * time.Second

	workers := envOr("WORKERS", "4")
	w, err := strconv.Atoi(workers)
	if err != nil || w <= 0 {
		w = 4
	}
	c.Workers = w

	return c, nil
}

func envOr(k, def string) string {
	v := os.Getenv(k)
	if v == "" {
		return def
	}
	return v
}
