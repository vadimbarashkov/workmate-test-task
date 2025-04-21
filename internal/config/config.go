package config

import (
	"fmt"
	"os"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/goccy/go-yaml"
)

const (
	EnvDev  = "dev"
	EnvTest = "test"
	EnvProd = "prod"
)

type Server struct {
	Port           int           `yaml:"port" validate:"required,min=1,max=65535"`
	ReadTimeout    time.Duration `yaml:"read_timeout" validate:"gt=0"`
	WriteTimeout   time.Duration `yaml:"write_timeout" validate:"gt=0"`
	IdleTimeout    time.Duration `yaml:"idle_timeout" validate:"gt=0"`
	MaxHeaderBytes int           `yaml:"max_header_bytes" validate:"gte=0"`
}

var defaultServer = Server{
	Port:           8080,
	ReadTimeout:    5 * time.Second,
	WriteTimeout:   10 * time.Second,
	IdleTimeout:    time.Minute,
	MaxHeaderBytes: 1 << 20,
}

func (s *Server) Addr() string {
	return fmt.Sprintf(":%d", s.Port)
}

type Config struct {
	Env    string `yaml:"env" validate:"required,oneof=dev test prod"`
	Server Server `yaml:"server" validate:"required"`
}

var defaultConfig = Config{
	Env:    EnvDev,
	Server: defaultServer,
}

var validate = validator.New()

func (c *Config) validate() error {
	return validate.Struct(c)
}

func Load(path string) (*Config, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("open config file: %w", err)
	}
	defer f.Close()

	cfg := defaultConfig
	if err := yaml.NewDecoder(f).Decode(&cfg); err != nil {
		return nil, fmt.Errorf("decode config file: %w", err)
	}

	if err := cfg.validate(); err != nil {
		return nil, fmt.Errorf("invalid config: %w", err)
	}

	return &cfg, nil
}
