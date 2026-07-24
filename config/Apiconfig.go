package config

import (
	"database/sql"
)

type ApiConfig struct {
	DB *sql.DB
	// Rabbit *rabbitmq.RabbitMQ
}
