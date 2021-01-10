package controller

import (
	"errors"

	"card-keeper-api/cardservice"
	configs "card-keeper-api/internal/configs"
	logger "card-keeper-api/internal/logging"
	"card-keeper-api/middleware"

	jwtmiddleware "github.com/auth0/go-jwt-middleware"
	"github.com/gin-gonic/gin"
)

var routerLogger = logger.NewLogger()

// InitServer registers the routes for the application
func InitServer(configs configs.Configuration) *gin.Engine {
	router := gin.New()

	router.Use(middleware.LogToFile())
	router.Use(middleware.CorsMiddleware(configs.APIAllowedOrigin()))

	controller := setupController(configs.DBConfigs())

	v1 := router.Group("v1")

	jwtMiddleware := middleware.GetJWTMiddleware(configs.AuthConfiguration())

	v1.GET("/collection", checkJWT(jwtMiddleware), controller.GetCollection)
	v1.POST("/collection", checkJWT(jwtMiddleware), controller.AddToCollection)

	router.GET("/ping", controller.Ping)

	return router
}

func checkJWT(jwtMiddleware *jwtmiddleware.JWTMiddleware) gin.HandlerFunc {
	return func(c *gin.Context) {
		if err := jwtMiddleware.CheckJWT(c.Writer, c.Request); err != nil {
			routerLogger.LogInfo(err.Error())
			c.AbortWithStatus(401)
		}
	}
}

func setupController(configs configs.DBConfiguration) *Controller {
	controller := new(Controller)

	repo, err := initializeRepository(configs)

	if err != nil {
		routerLogger.LogErrorWithFields(
			logger.LogFields{
				"err": err,
			}, "not able to instantiate the requested repo configuration")
		routerLogger.LogFatal("fatal error creating controller")
	}

	s := cardservice.Service{
		Repository: repo,
	}

	controller.Service = &s

	return controller
}

func initializeRepository(dbConfig configs.DBConfiguration) (cardservice.Repository, error) {
	var configuredRepo cardservice.Repository
	var err error

	switch repo := dbConfig.Type; repo {
	case "in-memory":
		configuredRepo, err = cardservice.InMemoryStore()
	case "mongodb":
		configuredRepo, err = cardservice.MongoDB(dbConfig)
	default:
		err = errors.New("unsupported repository")
	}

	return configuredRepo, err
}
