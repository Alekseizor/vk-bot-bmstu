package config

import (
	"context"
	"fmt"
	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"os"
)

// Config Структура конфигурации;
// Содержит все конфигурационные данные о сервисе;
// автоподгружается при изменении исходного файла
type Config struct {
	VKToken    string
	Redis      RedisConfig
	BITOPToken string
	BITOP      BitopConfig
}

type BitopConfig struct {
	SiteAdress string
	Protocol   string
	PathPath   string
	PathSearch string
}

type RedisConfig struct {
	// from config file
	TTL         int
	DialTimeout int
	ReadTimeout int
	// from env
	Host     string
	Port     int
	User     string
	Password string
}

// NewConfig Создаёт новый объект конфигурации, загружая данные из файла конфигурации
func NewConfig(ctx context.Context) (*Config, error) {
	var err error
	cfg := &Config{}

	configName := "config"
	_ = godotenv.Load()
	if os.Getenv("CONFIG_NAME") != "" {
		configName = os.Getenv("CONFIG_NAME")
	}

	cfg.Redis.Password = os.Getenv("REDIS_PASS")

	if token := os.Getenv("VK_TOKEN"); token == "" {
		return nil, fmt.Errorf("")
	} else {
		cfg.VKToken = token
	}

	if token := os.Getenv("BITOP_TOKEN"); token == "" {
		return nil, fmt.Errorf("")
	} else {
		cfg.BITOPToken = token
	}

	viper.SetConfigName(configName)
	viper.SetConfigType("toml")
	viper.AddConfigPath("config")
	viper.AddConfigPath(".")
	viper.WatchConfig()

	err = viper.ReadInConfig()
	if err != nil {
		return nil, err
	}

	err = viper.Unmarshal(cfg)
	if err != nil {
		return nil, err
	}

	log.Info("config parsed")

	return cfg, nil
}
