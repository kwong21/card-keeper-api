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
	"gorm.io/driver/postgres"
	_ "gorm.io/driver/postgres"
	"gorm.io/gorm"
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

	db, err := setupDatabase()

	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	repo := service.Database(db)
	s := service.Service{
		Repository: repo,
	}

	c.Service = &s

	r.POST("/collection", c.AddToCollection)

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

func setupTestControllerAndHTTPRecorder() (*httptest.ResponseRecorder, *gin.Engine, *controller.Controller) {
	w := httptest.NewRecorder()
	r := gin.Default()
	c := new(controller.Controller)

	return w, r, c
}

func setupDatabase() (*gorm.DB, error) {
	dsn := "user=tester dbname=cardkeeper port=26257 sslmode=disable"
	//	dsn := "user=keeperuser dbname=cardkeeper port=26257 sslmode=prefer sslrootcert=certs/ca.crt sslcert=certs/client.keeperuser.crt sslkey=certs/client.keeperuser.key"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	db.Migrator().CreateTable(&model.Card{})

	return db, err
}

func verifyHTTPResponseBody(expected string, actual string, t *testing.T) {
	if actual != expected {
		t.Errorf("Expected message body of %v, but got %v", expected, actual)
	}
}
