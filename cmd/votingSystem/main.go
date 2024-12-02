// func main() {
// 	c := config.MustReadConfig("config.json")
// 	fmt.Println(c)
// 	helper.LoadEnvFile()

// }
package main

import (
	"flag"
	"log"
	"os"

	"github.com/porseOnline/api/handlers/http"
	"github.com/porseOnline/app"
	"github.com/porseOnline/config"
)

var configPath = flag.String("config", "config.json", "service configuration file")

func main() {
	flag.Parse()

	if v := os.Getenv("CONFIG_PATH"); len(v) > 0 {
		*configPath = v
	}

	c := config.MustReadConfig(*configPath)

	appContainer := app.NewMustApp(c)

	log.Fatal(http.Run(appContainer, c.Server))
}
