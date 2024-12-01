package main

import (
	"PorsOnlineWebApp/app"
	"PorsOnlineWebApp/config"
	"flag"
	"log"
	"os"

	"PorsOnlineWebApp/api/handler/http"
)

func main() {
	var configPath = flag.String("config", "config.yml", "configuration file path")
	flag.Parse()

	if v := os.Getenv("CONFIG_PATH"); len(v) > 0 {
		*configPath = v
	}
	c := config.MustReadConfig(*configPath)
	appContainer := app.NewMustApp(c)
	err := http.Run(appContainer, c)
	if err != nil {
		log.Fatal("can not start the programm")
	}
}
