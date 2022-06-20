package config

import (
	"log"

	"github.com/spf13/viper"
)

var config *viper.Viper

// New is an exported method that takes the environment starts the viper
// (external lib) and returns the configuration struct.
func New(env string) {
	var err error
	config = viper.New()
	config.SetConfigType("yaml")
	config.SetConfigName(env)
	config.AddConfigPath("../config/")
	config.AddConfigPath("config/")
	err = config.ReadInConfig()
	if err != nil {
		log.Fatal("error on parsing configuration file")
	}
}

func GetConfig() *viper.Viper {
	return config
}
