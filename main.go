package main

import (
	"log"
	"net/http"

	"github.com/spf13/viper"
)

func main() {
	SetUpConfig()
	router := GetRouter()
	log.Println("Server is ready")
	http.ListenAndServe(":"+viper.GetString("port"), router)
}
