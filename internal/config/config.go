package config

import (
	"sync"

	"github.com/Traunin/review-assigner/internal/env"
	_ "github.com/lib/pq"
)

type Config struct {
	dbHost     string
	dbPort     string
	dbUser     string
	dbPassword string
	dbName     string
	port       string
}

var (
	cfg  *Config
	once sync.Once
)

func (c *Config) DBHost() string     { return c.dbHost }
func (c *Config) DBPort() string     { return c.dbPort }
func (c *Config) DBUser() string     { return c.dbUser }
func (c *Config) DBPassword() string { return c.dbPassword }
func (c *Config) DBName() string     { return c.dbName }
func (c *Config) Port() string       { return c.port }

func Load() *Config {
	once.Do(func() {
		cfg = &Config{
			dbHost:     env.Must("DB_HOST"),
			dbPort:     env.Must("DB_PORT"),
			dbUser:     env.Must("DB_USER"),
			dbPassword: env.Must("DB_PASSWORD"),
			dbName:     env.Must("DB_NAME"),
			port:       env.Fallback("SERVER_PORT", "8080"),
		}
	})

	return cfg
}
