package controller

import (
	"bytes"
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
	card := model.Card{
		Year:   2020,
		Maker:  "Upper Deck",
		Set:    "Series One",
		Player: "Brock Boeser",
	}

	w, r, c := setupTestControllerAndHTTPRecorder()

	repo := service.InMemoryStore()
	s := service.Service{
		Repository: repo,
	}

	c.Service = &s

	r.POST("/collection", c.Collection)

	b, _ := json.Marshal(card)
	req, err := http.NewRequest("POST", "/collection", bytes.NewBuffer(b))

	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	r.ServeHTTP(w, req)

	// Verify that the POST request succeeded with HTTP 200
	if w.Code != http.StatusOK {
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

// TestAddNewCardError expects error to be returned when data is incorrect
func TestAddNewCardError(t *testing.T) {
	w, r, c := setupTestControllerAndHTTPRecorder()
	r.POST("/collection", c.Collection)

	b, _ := json.Marshal("foo bar")
	req, err := http.NewRequest("POST", "/collection", bytes.NewBuffer(b))

	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected to get HTTP 400, but got %d", w.Code)
	}

	// Verify message body gives `ok`
	expected := `{"message":"invalid data"}`
	verifyHTTPResponseBody(expected, w.Body.String(), t)
}

func setupTestControllerAndHTTPRecorder() (*httptest.ResponseRecorder, *gin.Engine, *Controller) {
	w := httptest.NewRecorder()
	r := gin.Default()
	c := new(Controller)

	return w, r, c
}

func verifyHTTPResponseBody(expected string, actual string, t *testing.T) {
	if actual != expected {
		t.Errorf("Expected message body of %v, but got %v", expected, actual)
	}
}
