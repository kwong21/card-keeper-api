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

func setResponse(c *gin.Context, m string, s int) *gin.Context {
	c.JSON(s, gin.H{
		"message": m,
	})
	return c
}
