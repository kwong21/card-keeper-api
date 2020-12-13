package controller

import "github.com/gin-gonic/gin"

// InitRouter registers the routes for the application
func InitRouter() *gin.Engine {
	router := gin.Default()

	v1 := router.Group("v1")
	controller := new(Controller)

	v1.POST("/collection", checkJWT(), controller.AddToCollection)
	return router
}

func checkJWT() gin.HandlerFunc {
	return func(c *gin.Context) {
		jwtMiddleware := getJWTMiddleware()

		if err := jwtMiddleware.CheckJWT(c.Writer, c.Request); err != nil {
			c.AbortWithStatus(401)
		}
	}
}
