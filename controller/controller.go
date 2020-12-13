package controller

import (
	"card-keeper-api/model"
	"card-keeper-api/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Controller struct
type Controller struct {
	Service *service.Service
}

// AddToCollection accepts POST request for adding card to collection
func (controller *Controller) AddToCollection(c *gin.Context) {
	var newCard model.Card
	error := c.BindJSON(&newCard)

	if error != nil {
		setResponse(c, "invalid data", http.StatusBadRequest)
	} else {
		err := controller.Service.AddCard(newCard)

		if err != nil {
			setResponse(c, "internal error", http.StatusInternalServerError)
		}
		setResponse(c, "ok", http.StatusAccepted)
	}
	return
}

// Ping checks status of API server
// @TODO make this more useful?
func (controller *Controller) Ping(c *gin.Context) {
	setResponse(c, "pong", http.StatusOK)
}

func setResponse(c *gin.Context, m string, s int) *gin.Context {
	c.JSON(s, gin.H{
		"message": m,
	})
	return c
}
