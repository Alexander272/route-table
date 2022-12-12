package config

import (
	"fmt"
	"os"
	"time"

	"github.com/kelseyhightower/envconfig"
	"github.com/spf13/viper"
)

type (
	Config struct {
		Environment string
		Redis       *RedisConfig
		Postgres    *PostgresConfig
		Auth        *AuthConfig
		Http        *HttpConfig
		Limiter     *LimiterConfig
		Urgency     *UrgencyConfig
		Order       *OrderConfig
	}

	RedisConfig struct {
		Host     string `mapstructure:"Host"`
		Port     string `mapstructure:"Port"`
		DB       int    `mapstructure:"DB"`
		Password string
	}

	PostgresConfig struct {
		Host     string `mapstructure:"host"`
		Port     string `mapstructure:"port"`
		Username string `mapstructure:"username"`
		Password string
		DbName   string `mapstructure:"dbname"`
		SSLMode  string `mapstructure:"sslmode"`
	}

	AuthConfig struct {
		AccessTokenTTL  time.Duration `mapstructure:"accessTokenTTL"`
		RefreshTokenTTL time.Duration `mapstructure:"refreshTokenTTL"`
		LimitAuthTTL    time.Duration `mapstructure:"limitAuthTTL"`
		CountAttempt    int32         `mapstructure:"countAttempt"`
		Secure          bool          `mapstructure:"secure"`
		Domain          string        `mapstructure:"domain"`
		Key             string
	}

	HttpConfig struct {
		Host               string        `mapstructure:"host"`
		Port               string        `mapstructure:"port"`
		ReadTimeout        time.Duration `mapstructure:"readTimeout"`
		WriteTimeout       time.Duration `mapstructure:"writeTimeout"`
		MaxHeaderMegabytes int           `mapstructure:"maxHeaderBytes"`
	}

	LimiterConfig struct {
		RPS   int           `mapstructure:"rps"`
		Burst int           `mapstructure:"burst"`
		TTL   time.Duration `mapstructure:"ttl"`
	}

	UrgencyConfig struct {
		High   time.Duration `mapstructure:"high"`
		Middle time.Duration `mapstructure:"middle"`
	}

	OrderConfig struct {
		Delay time.Duration `mapstructure:"delay"`
		Term  time.Duration `mapstructure:"term"`
	}
)

// Инициализация конфига прогаммы
func Init(configDir string) (conf *Config, err error) {
	if err := parseConfigFile(configDir); err != nil {
		return nil, fmt.Errorf("failed to parse config. err: %w", err)
	}

	conf = &Config{}
	if err := unmarhal(conf); err != nil {
		return nil, fmt.Errorf("failed to unmarhal value. error: %w", err)
	}
	if err := setFromEnv(conf); err != nil {
		return nil, fmt.Errorf("failed to get from env. error: %w", err)
	}

	return conf, nil
}

func parseConfigFile(folder string) error {
	viper.AddConfigPath(folder)
	viper.SetConfigName("config")

	return viper.MergeInConfig()
}

func unmarhal(conf *Config) error {
	conf.Redis = &RedisConfig{}
	if err := viper.UnmarshalKey("redis", conf.Redis); err != nil {
		return fmt.Errorf("failed to unmarhal key redis. error: %w", err)
	}
	conf.Postgres = &PostgresConfig{}
	if err := viper.UnmarshalKey("postgres", conf.Postgres); err != nil {
		return fmt.Errorf("failed to unmarhal key postgres. error: %w", err)
	}
	conf.Http = &HttpConfig{}
	if err := viper.UnmarshalKey("http", conf.Http); err != nil {
		return fmt.Errorf("failed to unmarhal key http. error: %w", err)
	}
	conf.Limiter = &LimiterConfig{}
	if err := viper.UnmarshalKey("limiter", conf.Limiter); err != nil {
		return fmt.Errorf("failed to unmarhal key limiter. error: %w", err)
	}
	conf.Urgency = &UrgencyConfig{}
	if err := viper.UnmarshalKey("urgency", conf.Urgency); err != nil {
		return fmt.Errorf("failed to unmarhal key urgency. error: %w", err)
	}
	conf.Order = &OrderConfig{}
	if err := viper.UnmarshalKey("order", conf.Order); err != nil {
		return fmt.Errorf("failed to unmarhal key order. error: %w", err)
	}
	conf.Auth = &AuthConfig{}
	if err := viper.UnmarshalKey("auth", conf.Auth); err != nil {
		return fmt.Errorf("failed to unmarhal key auth. error: %w", err)
	}

	return nil
}

func setFromEnv(conf *Config) error {
	if err := envconfig.Process("http", conf.Http); err != nil {
		return fmt.Errorf("failed to get http from env. error: %w", err)
	}
	if err := envconfig.Process("jwt", conf.Auth); err != nil {
		return fmt.Errorf("failed to get jwt from env. error: %w", err)
	}
	if err := envconfig.Process("redis", conf.Redis); err != nil {
		return fmt.Errorf("failed to get redis from env. error: %w", err)
	}
	if err := envconfig.Process("postgres", conf.Postgres); err != nil {
		return fmt.Errorf("failed to get postgres from env. error: %w", err)
	}

	conf.Environment = os.Getenv("APP_ENV")

	return nil
}

func (u *UrgencyConfig) ChangeUrgency(high time.Duration, middle time.Duration) error {
	viper.Set("urgency.high", high)
	viper.Set("urgency.middle", middle)
	err := viper.WriteConfig()
	if err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	u.High = high
	u.Middle = middle
	return nil
}
