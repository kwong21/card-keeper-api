package controller

import (
	"card-keeper-api/cardservice"

	"net/http"

	"github.com/gin-gonic/gin"
)

// Controller struct
type Controller struct {
	Service *cardservice.Service
}

// AddToCollection accepts POST request for adding card to collection
func (controller *Controller) AddToCollection(c *gin.Context) {
	var newCard cardservice.Card
	error := c.BindJSON(&newCard)
	collection := c.Param("collection")

	if error != nil {
		setResponse(c, "invalid data", http.StatusBadRequest)
	} else {
		err := controller.Service.AddCardToCollection(newCard, collection)

		if err != nil {
			msg, code := checkErrorAndReturnStatus(err)
			setResponse(c, msg, code)
		} else {
			setResponse(c, "ok", http.StatusAccepted)
		}
	}
	return
}

// GetCollection accepts GET request and get cards in collection
func (controller *Controller) GetCollection(c *gin.Context) {
	collection := c.Param("collection")
	cards, err := controller.Service.GetAllCardsInCollection(collection)

	if err != nil {
		setResponse(c, "error getting cards", http.StatusInternalServerError)
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "ok",
		"cards":   cards,
	})

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
	case *cardservice.DuplicateError:
		return "duplicate item", http.StatusConflict
	default:
		return "internal server error", http.StatusInternalServerError
	}
}
