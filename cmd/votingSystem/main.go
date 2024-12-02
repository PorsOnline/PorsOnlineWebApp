package main

import (
	"PorsOnlineWebApp/config"
	"flag"
	"log"
	"os"

	"PorsOnlineWebApp/api/handler/http"
	"PorsOnlineWebApp/pkg/logger"
)

func main() {
	var configPath = flag.String("config", "config.yml", "configuration file path")
	flag.Parse()

	if v := os.Getenv("CONFIG_PATH"); len(v) > 0 {
		*configPath = v
	}
	c := config.MustReadConfig(*configPath)
	err := logger.InitLogger(c)
	if err != nil {
		log.Fatal("can not initialize logger")
	}
	logger.Info("Starting the program", nil)
	err = http.Run(c)
	if err != nil {
		log.Fatal("can not start the programm")
	}
}
