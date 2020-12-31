package controller

import (
	"card-keeper-api/config"
	logger "card-keeper-api/log"
	"card-keeper-api/middleware"
	"card-keeper-api/service"
	"errors"

	"github.com/gin-gonic/gin"
)

var routerLogger = logger.NewLogger()

// InitServer registers the routes for the application
func InitServer(configs config.Configuration) *gin.Engine {
	router := gin.New()

	router.Use(middleware.LogToFile())
	router.Use(middleware.CorsMiddleware(configs.APIAllowedOrigin()))

	controller := setupController(configs.DBConfigs())

	v1 := router.Group("v1")

	v1.POST("/collection", checkJWT(), controller.AddToCollection)
	router.GET("/ping", controller.Ping)

	return router
}

func checkJWT() gin.HandlerFunc {
	return func(c *gin.Context) {
		jwtMiddleware := middleware.JWTMiddleware()

		if err := jwtMiddleware.CheckJWT(c.Writer, c.Request); err != nil {
			c.AbortWithStatus(401)
		}
	}
}

func setupController(configs config.DBConfiguration) *Controller {
	controller := new(Controller)

	repo, err := initializeRepository(configs)

	if err != nil {
		routerLogger.LogErrorWithFields(
			logger.Fields{
				"err": err,
			}, "not able to instantiate the requested repo configuration")
		routerLogger.LogFatal("fatal error creating controller")
	}

	s := service.Service{
		Repository: repo,
	}

	controller.Service = &s

	return controller
}

func initializeRepository(dbConfig config.DBConfiguration) (service.Repository, error) {
	var configuredRepo service.Repository
	var err error

	switch repo := dbConfig.Type; repo {
	case "in-memory":
		configuredRepo, err = service.InMemoryStore()
	case "mongodb":
		configuredRepo, err = service.MongoDB(dbConfig)
	default:
		err = errors.New("unsupported repository")
	}

	return configuredRepo, err
}
