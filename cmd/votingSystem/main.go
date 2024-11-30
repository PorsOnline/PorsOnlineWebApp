package main

import (
	"fmt"

	"github.com/porseOnline/config"
	"github.com/porseOnline/pkg/helper"
	db "github.com/porseOnline/pkg/postgres"
)

func main() {
	c := config.MustReadConfig("config.json")
	fmt.Println(c)
	helper.LoadEnvFile()
	db.NewConnection()
}
