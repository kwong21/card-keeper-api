// +build integration

package controller

import (
	"bytes"
	"card-keeper-api/config"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

var engine *gin.Engine
var controller *Controller

func init() {
	engine, controller = configureITTestEnvironment()

	engine.POST("/collection", controller.AddToCollection)
}

// TestAddNewCardToRepoIT verifies behaviour for adding a new card.
func TestAddNewCardToRepoIT(t *testing.T) {
	b := getSerializedTestCard()
	req, err := http.NewRequest("POST", "/collection", bytes.NewBuffer(b))

	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	httpResponseRecorder := httptest.NewRecorder()
	recordRequest(req, httpResponseRecorder)

	// Verify that the POST request succeeded with HTTP 200
	if httpResponseRecorder.Code != http.StatusAccepted {
		t.Errorf("Expected to get HTTP 202, but got %d", httpResponseRecorder.Code)
		t.Fail()
	}

	// Verify message body gives `ok`
	expected := `{"message":"ok"}`
	verifyHTTPResponseBody(expected, httpResponseRecorder.Body.String(), t)

	// Verify that the card is in the Stored
	cards := controller.Service.GetAll()

	if len(*cards) != 1 {
		t.Errorf("Expected to get 1 card, but got %d", len(*cards))
		t.Fail()
	}

	httpResponseRecorder.Flush()
}

// TestDuplicateCardNotAdded verifies that duplicates are not added to collection
func TestDuplicateCardNotAdded(t *testing.T) {
	b := getSerializedTestCard()
	req, err := http.NewRequest("POST", "/collection", bytes.NewBuffer(b))

	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	httpResponseRecorder := httptest.NewRecorder()
	recordRequest(req, httpResponseRecorder)

	// Verify that the POST request failed with 409
	if httpResponseRecorder.Code != http.StatusConflict {
		t.Errorf("Expected to get HTTP 409, but got %d", httpResponseRecorder.Code)
		t.Fail()
	}

	// Verify message body gives `ok`
	expected := `{"message":"duplicate item"}`
	verifyHTTPResponseBody(expected, httpResponseRecorder.Body.String(), t)

	httpResponseRecorder.Flush()
}

func configureITTestEnvironment() (*gin.Engine, *Controller) {
	r := gin.New()
	c := setupControllerWithDatabaseBackend()

	return r, c
}

func setupControllerWithDatabaseBackend() *Controller {
	dbConfigs := config.DBConfiguration{
		Type:     "mongodb",
		Host:     "localhost:27017",
		Database: "card-keeper-it",
	}

	return setupController(dbConfigs)
}

func recordRequest(req *http.Request, recorder *httptest.ResponseRecorder) {
	engine.ServeHTTP(recorder, req)
}
