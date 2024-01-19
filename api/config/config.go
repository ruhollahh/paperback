package config

import (
	"github.com/ruhollahh/paperback/internal/app/service"
)

const Version = "1.0.0"

type Config struct {
	Port    int
	Env     string
	Db      service.DBConfig
	Limiter struct {
		Rps     float64
		Burst   int
		Enabled bool
	}
	Smtp struct {
		Host     string
		Port     int
		Username string
		Password string
		Sender   string
	}
	Cors struct {
		TrustedOrigins []string
	}
}
