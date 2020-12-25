package it

import (
	"bytes"
	"card-keeper-api/controller"
	"card-keeper-api/model"
	"card-keeper-api/service"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

// TestAddsNewCardToRepo verifies behaviour for adding a new card.
func TestAddsNewCardToRepo(t *testing.T) {
	w, r, c := setupTestControllerAndHTTPRecorder()

	repo, err := service.MongoDB()

	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	s := service.Service{
		Repository: repo,
	}

	c.Service = &s

	r.POST("/collection", c.AddToCollection)

	b := getSerializedTestCard()
	req, err := http.NewRequest("POST", "/collection", bytes.NewBuffer(b))

	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	r.ServeHTTP(w, req)

	// Verify that the POST request succeeded with HTTP 200
	if w.Code != http.StatusAccepted {
		t.Errorf("Expected to get HTTP 200, but got %d", w.Code)
		t.Fail()
	}

	// Verify message body gives `ok`
	expected := `{"message":"ok"}`
	verifyHTTPResponseBody(expected, w.Body.String(), t)

	// Verify that the card is in the Stored
	cards := c.Service.GetAll()
	if len(*cards) != 1 {
		t.Errorf("Expected to get 1 card, but got %d", len(*cards))
		t.Fail()
	}
}

func getSerializedTestCard() []byte {
	base := model.Base{
		Year:   2020,
		Make:   "Upper Deck",
		Set:    "Series One",
		Player: "Brock Boeser",
	}

	insert := model.Insert{}

	card := model.Card{
		Base:   base,
		Insert: insert,
	}

	b, _ := json.Marshal(card)

	return b
}

func setupTestControllerAndHTTPRecorder() (*httptest.ResponseRecorder, *gin.Engine, *controller.Controller) {
	w := httptest.NewRecorder()
	r := gin.Default()
	c := new(controller.Controller)

	return w, r, c
}

func verifyHTTPResponseBody(expected string, actual string, t *testing.T) {
	if actual != expected {
		t.Errorf("Expected message body of %v, but got %v", expected, actual)
	}
}
