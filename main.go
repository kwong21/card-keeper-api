package main

import (
	"card-keeper-api/config"
	"card-keeper-api/controller"
	"flag"
	"io"
	"os"

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

	initLogging(configs.APILogPath())

	r := controller.InitRouter()

	r.Run(":8080")
}

func initLogging(logpath string) {
	f, _ := os.Create(logpath)
	gin.DefaultWriter = io.MultiWriter(f)
}
