package config

import (
	"fmt"

	"github.com/spf13/viper"
)

const (
	EnvLocal = "local"
	EnvProd  = "prod"
)

type SQLite struct {
	Database string `mapstructure:"SQLITE_DB"`
}

type Config struct {
	Env    string `mapstructure:"APP_ENV"`
	SQLite SQLite
}

func NewConfig(folder string) (*Config, error) {
	cfg := new(Config)

	viper.SetDefault("APP_ENV", EnvLocal)
	viper.SetDefault("SQLITE_DB", "sqllite.db")

	viper.SetConfigFile(".env")
	viper.AutomaticEnv() // Automatically override with environment variables

	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	env := viper.GetString("APP_ENV")

	if env == "" {
		return nil, fmt.Errorf("environment not set")
	}

	// if env is not local or prod, return error
	if env != EnvLocal && env != EnvProd {
		return nil, fmt.Errorf("invalid environment: %s", env)
	}

	viper.SetConfigName(env)
	viper.AddConfigPath(folder)
	viper.SetConfigType("yml") // Look for specific type
	//viper.MergeInConfig()
	if err := viper.MergeInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			return nil, fmt.Errorf("config file not found")
		}
		return nil, err
	}

	// Unmarshal the configuration into the Config struct
	if err := viper.Unmarshal(cfg); err != nil {
		return nil, err
	}

	// Unmarshal the configuration into the SQLite struct
	if err := viper.Unmarshal(&cfg.SQLite); err != nil {
		return nil, err
	}

	return cfg, nil
}
