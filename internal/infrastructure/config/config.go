package config

import (
	"fmt"
	"os"

	"github.com/spf13/viper"
)

type Config struct {
	MongoURI     string `yaml:"mongoUri" json:"mongoUri"`
	DatabaseName string `yaml:"databaseName" json:"databaseName"`
	JWTSecret    string `yaml:"jwtSecret" json:"jwtSecret"`
	HTTPPort     string `yaml:"httpPort" json:"httpPort"`
	GRPCPort     string `yaml:"grpcPort" json:"grpcPort"`
}

// Load reads configuration from environment variables or defaults.
func Load() *Config {
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.SetConfigName("config")

	if err := viper.ReadInConfig(); err != nil {
		fmt.Println("Error reading config file:", err)
		os.Exit(1)
	}
	viper.AutomaticEnv()

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		fmt.Println("Error unmarshalling config:", err)
	}
	return &config
}
