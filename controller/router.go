package controller

import "github.com/gin-gonic/gin"

// InitRouter registers the routes for the application
func InitRouter() *gin.Engine {
	router := gin.Default()

	v1 := router.Group("v1")
	controller := new(Controller)

	v1.POST("/colleciton", controller.AddToCollection)
	return router
}
