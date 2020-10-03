package main

import (
	"log"
	"net/http"
	"os"

	"github.com/spf13/viper"
)

func main() {
	viper.SetConfigFile(os.Args[1])
	router := GetRouter()
	log.Println("Server is ready")
	http.ListenAndServe(":"+viper.GetString("port"), router)
}
