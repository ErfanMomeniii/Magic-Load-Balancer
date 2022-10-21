package config

import (
	"fmt"
	"github.com/mitchellh/mapstructure"
	"github.com/spf13/viper"
	"strings"
	"time"
)

var (
	Default = "config.yaml"
	C       *Config
)

type Config struct {
	HTTPServer HTTPServer `yaml:"http-server"`
	Endpoints  []Endpoint `yaml:"endpoints"`
	Logger     Logger     `yaml:"logger"`
	Tracing    Tracing    `yaml:"tracing"`
}

// Server config values
type Server struct {
	IP          string        `yaml:"ip"`
	WorkingTime time.Duration `yaml:"working_time"`
}

// Service config values
type Service struct {
	Name    string   `yaml:"name"`
	Servers []Server `yaml:"servers"`
}

// Tracing config struct
type Tracing struct {
	Enabled      bool    `yaml:"enabled"`
	AgentHost    string  `yaml:"agent_host"`
	AgentPort    string  `yaml:"agent_port"`
	SamplerRatio float64 `yaml:"sampler_ratio"`
}

// Endpoint config values
type Endpoint struct {
	URL     string  `yaml:"url"`
	Service Service `yaml:"service"`
}

// Logger config values
type Logger struct {
	Level string `yaml:"level"`
}

// HTTPServer config values
type HTTPServer struct {
	Listen            string        `yaml:"listen"`
	ReadTimeout       time.Duration `yaml:"read_Timeout"`
	WriteTimeout      time.Duration `yaml:"write_timeout"`
	ReadHeaderTimeout time.Duration `yaml:"read_header_timeout"`
	IdleTimeout       time.Duration `yaml:"idle_timeout"`
}

func Init(filename string) *Config {
	c := new(Config)
	v := viper.New()
	v.SetConfigType("yaml")
	v.AddConfigPath(".")
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_", "-", "_"))
	v.AutomaticEnv()
	v.SetConfigName(Default)

	if err := v.ReadInConfig(); err != nil {
		panic(fmt.Errorf("error loading configs: %s", err.Error()).(any))
	}

	err := v.Unmarshal(c, func(config *mapstructure.DecoderConfig) {
		config.TagName = "yaml"
		config.DecodeHook = mapstructure.ComposeDecodeHookFunc(
			mapstructure.StringToTimeDurationHookFunc(),
			mapstructure.StringToSliceHookFunc(","),
		)
	})
	if err != nil {
		panic(fmt.Errorf("failed on config `%s` unmarshal: %w", Default, err).(any))
	}

	C = c

	return c
}
