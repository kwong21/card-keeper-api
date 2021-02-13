package main

import (
	"flag"
	"fmt"

	"card-keeper-api/controller"
	config "card-keeper-api/internal/configs"

	"github.com/gin-gonic/gin"
)

func main() {
	gin.DisableConsoleColor()

	var configs config.Configuration
	var c string

	flag.StringVar(&c, "conf", "", "Configuration file")
	flag.Parse()

	if c == "" {
		configs = config.Default()
	} else {
		configs = config.NewFromFile(c)
	}

	server := controller.InitServer(configs)

	port := fmt.Sprintf(":%s", configs.APIListenPort())
	server.Run(port)
}
