package config

import (
	"github.com/spf13/viper"
)

type Config struct {
	DBDriver          string `mapstructure:"DB_DRIVER"`
	DBSource          string `mapstructure:"DB_SOURCE"`
	TokenSymmetricKey string `mapstructure:"TOKEN_SYMMETRIC_KEY"`
	HTTPServerAddress string `mapstructure:"HTTP_SERVER_ADDRESS"`
}

func LoadConfig(path string) (config Config, err error) {
	viper.SetConfigName(".env")
	
	// Set the path to look for the configuration file
	viper.AddConfigPath(path)


	// Set default values
	viper.SetDefault("DB_DRIVER", "postgres")
	viper.SetDefault("DB_SOURCE", "postgres://root:secret@localhost:5433/defi?sslmode=disable")
	viper.SetDefault("HTTP_SERVER_ADDRESS", ":8080")
	
	// Enable VIPER to read environment variables
	viper.AutomaticEnv()
	
	// Set the type of the configuration file
	viper.SetConfigType("env")

	// Read the configuration file
	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	// Unmarshal the configuration into the Config struct
	err = viper.Unmarshal(&config)
	return
}
