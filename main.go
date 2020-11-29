package main

import (
	"card-keeper-api/controller"
	"io"
	"os"

	"github.com/gin-gonic/gin"
)

func main() {
	gin.DisableConsoleColor()

	f, _ := os.Create("card-keeper-api.log")
	gin.DefaultWriter = io.MultiWriter(f)

	r := controller.InitRouter()

	r.Run(":8080")
}
