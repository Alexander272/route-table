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
		Redis       RedisConfig
		Postgres    PostgresConfig
		Auth        AuthConfig
		Http        HttpConfig
		Limiter     LimiterConfig
		Urgency     UrgencyConfig
		Order       OrderConfig
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
func Init(configDir string) (*Config, error) {
	if err := parseConfigFile(configDir); err != nil {
		return nil, err
	}

	var conf Config
	if err := unmarhal(&conf); err != nil {
		return nil, err
	}
	if err := setFromEnv(&conf); err != nil {
		return nil, err
	}

	return &conf, nil
}

func parseConfigFile(folder string) error {
	viper.AddConfigPath(folder)
	viper.SetConfigName("config")

	return viper.MergeInConfig()
}

func unmarhal(conf *Config) error {
	if err := viper.UnmarshalKey("redis", &conf.Redis); err != nil {
		return err
	}
	if err := viper.UnmarshalKey("postgres", &conf.Postgres); err != nil {
		return err
	}
	if err := viper.UnmarshalKey("http", &conf.Http); err != nil {
		return err
	}
	if err := viper.UnmarshalKey("limiter", &conf.Limiter); err != nil {
		return err
	}
	if err := viper.UnmarshalKey("urgency", &conf.Urgency); err != nil {
		return err
	}
	if err := viper.UnmarshalKey("order", &conf.Order); err != nil {
		return err
	}
	if err := viper.UnmarshalKey("auth", &conf.Auth); err != nil {
		return err
	}

	return nil
}

func setFromEnv(conf *Config) error {
	if err := envconfig.Process("http", &conf.Http); err != nil {
		return err
	}
	if err := envconfig.Process("jwt", &conf.Auth); err != nil {
		return err
	}
	if err := envconfig.Process("redis", &conf.Redis); err != nil {
		return err
	}
	if err := envconfig.Process("postgres", &conf.Postgres); err != nil {
		return err
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
