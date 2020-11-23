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

// Collection handles requests to collection endpoint
func (controller *Controller) Collection(c *gin.Context) {
	var newCard model.Card
	error := c.BindJSON(&newCard)

	if error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "invalid data",
		})
	} else {
		err := controller.Service.AddCard(newCard)

		if err != nil {
			// Proper err handling
		}
		c.JSON(http.StatusOK, gin.H{
			"message": "ok",
		})
	}
	return
}
