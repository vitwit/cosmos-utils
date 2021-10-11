package config

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gopkg.in/go-playground/validator.v9"
)

type (
	// Scraper time interval
	Scraper struct {
		Rate string `mapstructure:"rate"`
	}
	// Config
	Config struct {
		LCDEndpoint    string  `mapstructure:"lcd_endpoint"`
		Deamon         string  `mapstructure:"deamon"`
		KeyName        string  `mapstructure:"key_name"`
		AccountAddress string  `mapstructure:"account_address"`
		ChainID        string  `mapstructure:"chain_id"`
		Fees           string  `mapstructure:"fees"`
		Scraper        Scraper `mapstructure:"scraper"`
	}
)

// ReadConfigFromFile to read config details from file using viper
func ReadConfigFromFile() (*Config, error) {
	v := viper.New()
	v.AddConfigPath(".")
	v.AddConfigPath("./config/")
	v.SetConfigName("config")
	if err := v.ReadInConfig(); err != nil {
		log.Fatalf("error while reading config.toml: %v", err)
	}

	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		log.Fatalf("error unmarshaling config.toml to application config: %v", err)
	}

	if err := cfg.Validate(); err != nil {
		log.Fatalf("error occurred in config validation: %v", err)
	}

	return &cfg, nil
}

// Validate config struct
func (c *Config) Validate(e ...string) error {
	v := validator.New()
	if len(e) == 0 {
		return v.Struct(c)
	}
	return v.StructExcept(c, e...)
}
