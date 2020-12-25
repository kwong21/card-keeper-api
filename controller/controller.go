package controller

import (
	logger "card-keeper-api/log"
	"card-keeper-api/model"
	"card-keeper-api/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Controller struct
type Controller struct {
	Service *service.Service
}

var controllerLogger = logger.NewLogger()

// AddToCollection accepts POST request for adding card to collection
func (controller *Controller) AddToCollection(c *gin.Context) {
	var newCard model.Card
	error := c.BindJSON(&newCard)

	if error != nil {
		setResponse(c, "invalid data", http.StatusBadRequest)
	} else {
		err := controller.Service.AddCard(newCard)

		if err != nil {
			msg, code := checkErrorAndReturnStatus(err)
			setResponse(c, msg, code)
		} else {
			setResponse(c, "ok", http.StatusAccepted)
		}
	}
	return
}

// Ping returns status of the API
func (controller *Controller) Ping(c *gin.Context) {
	setResponse(c, "pong", http.StatusOK)
}

func setResponse(c *gin.Context, m string, s int) *gin.Context {
	c.JSON(s, gin.H{
		"message": m,
	})
	return c
}

func checkErrorAndReturnStatus(err error) (string, int) {
	switch err := err; err.(type) {
	case *service.DuplicateError:
		return "duplicate item", http.StatusConflict
	default:
		return "internal server error", http.StatusInternalServerError
	}
}
