package config

import (
	"errors"
	"fmt"

	"github.com/spf13/viper"
)

type SQLite struct {
	Database string
}

type Config struct {
	Env    string
	SQLite SQLite
}

func NewConfig(folder, filename string) (*Config, error) {
	cfg := new(Config)

	//set default app environment
	viper.SetDefault("env", "local")

	//load from directory
	viper.SetConfigName(filename)
	viper.AddConfigPath(folder)

	//load env
	viper.AutomaticEnv()

	//read
	if err := viper.ReadInConfig(); err != nil {
		var configFileNotFoundError viper.ConfigFileNotFoundError
		if errors.As(err, &configFileNotFoundError) {
			return nil, fmt.Errorf("config file not found")
		} else {
			return nil, err
		}
	}

	//get config of app environment
	env := viper.Get("env")
	viper.SetConfigName(env.(string))

	//merge with default config
	if err := viper.MergeInConfig(); err != nil {
		fmt.Printf("Error reading config file, %s", err)
	}

	//unmarshal
	if err := viper.Unmarshal(cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}
