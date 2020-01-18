package config

import (
	"time"

	"github.com/kelseyhightower/envconfig"
)

// Config struct
type Config struct {
	API struct {
		Port           string        `default:":8090"`
		ReadTimeout    time.Duration `default:"5s"`
		WriteTimeout   time.Duration `default:"5s"`
		ProccesingTime time.Duration `default:"10m"`
	}
	DB struct {
		Host string `default:"localhost"`
		Name string `default:"processingdb"`
		User string `default:"postgres"`
		Port int    `default:"54320"`

		PoolSize     int           `default:"15"`
		MaxIdleConns int           `default:"15"`
		ConnLifetime time.Duration `default:"5m"`
		Tmpl         string        `default:"host=%s port=%d dbname=%s user=%s sslmode=disable application_name=%s"`
	}
}

// ParseConfig func
func ParseConfig(app string) (cfg Config, err error) {
	if err := envconfig.Process(app, &cfg); err != nil {
		if err := envconfig.Usage(app, &cfg); err != nil {
			return cfg, err
		}
		return cfg, err
	}
	return cfg, nil
}
