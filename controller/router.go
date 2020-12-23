package controller

import (
	"card-keeper-api/model"
	"card-keeper-api/service"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// InitRouter registers the routes for the application
func InitRouter() *gin.Engine {
	controller := setupController()

	router := gin.Default()
	router.Use(corsMiddleware())

	v1 := router.Group("v1")

	v1.POST("/collection", checkJWT(), controller.AddToCollection)

	router.GET("/ping", controller.Ping)

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

func setupController() *Controller {
	controller := new(Controller)

	db, _ := setupDatabase()

	repo := service.Database(db)
	s := service.Service{
		Repository: repo,
	}

	controller.Service = &s

	return controller
}

func setupDatabase() (*gorm.DB, error) {
	dsn := "user=tester dbname=cardkeeper port=26257 sslmode=disable"
	//	dsn := "user=keeperuser dbname=cardkeeper port=26257 sslmode=prefer sslrootcert=certs/ca.crt sslcert=certs/client.keeperuser.crt sslkey=certs/client.keeperuser.key"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	db.Migrator().CreateTable(&model.Card{})

	return db, err
}
