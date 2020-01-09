package utils

import (
	"fmt"
	"github.com/spf13/viper"
	"log"
)

var cfgFile string

type ServerConfig struct {
	Host       string `yaml:"host"`
	DbHost     string `yaml:"db_host"`
	DbName     string `yaml:"db_name"`
	DbUser     string `yaml:"db_user"`
	DbPassword string `yaml:"db_password"`
	AmpqHost string `yaml:"ampq_host"`
}


func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		viper.SetConfigName("conf")
		viper.AddConfigPath("configs")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}

func GetAppConfig() ServerConfig {

	initConfig()
	var appConfig ServerConfig

	err := viper.Unmarshal(&appConfig)
	if err != nil {
		log.Fatalf("unable to decode into struct, %v", err)
	}

	return appConfig
}

