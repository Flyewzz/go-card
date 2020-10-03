package main

import (
	"log"
	"os"

	"github.com/spf13/viper"
)

func SetUpConfig() {
	viper.SetConfigFile(os.Args[1])
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Cannot read a config file: %v\n", err)
	}
}
