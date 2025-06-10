package wati

import (
	"time"
)

// Config representa la configuración del cliente WATI
type Config struct {
	APIEndpoint string
	Token       string
	Timeout     time.Duration
	RetryCount  int
	RateLimit   *RateLimitConfig
	Debug       bool
}

// RateLimitConfig configura los límites de velocidad
type RateLimitConfig struct {
	RequestsPerSecond int
	BurstSize         int
}

// DefaultConfig retorna una configuración por defecto
func DefaultConfig() *Config {
	return &Config{
		Timeout:    30 * time.Second,
		RetryCount: 3,
		RateLimit: &RateLimitConfig{
			RequestsPerSecond: 10,
			BurstSize:         20,
		},
		Debug: false,
	}
}

// ClientOption es una función que modifica la configuración del cliente
type ClientOption func(*Config)

// WithTimeout establece el timeout para las peticiones HTTP
func WithTimeout(timeout time.Duration) ClientOption {
	return func(c *Config) {
		c.Timeout = timeout
	}
}

// WithRetryCount establece el número de reintentos
func WithRetryCount(count int) ClientOption {
	return func(c *Config) {
		c.RetryCount = count
	}
}

// WithRateLimit establece los límites de velocidad
func WithRateLimit(requestsPerSecond int, burstSize int) ClientOption {
	return func(c *Config) {
		c.RateLimit = &RateLimitConfig{
			RequestsPerSecond: requestsPerSecond,
			BurstSize:         burstSize,
		}
	}
}

// WithDebug habilita o deshabilita el modo debug
func WithDebug(debug bool) ClientOption {
	return func(c *Config) {
		c.Debug = debug
	}
}

