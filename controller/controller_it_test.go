package controller

import (
	configs "card-keeper-api/internal/configs"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/gin-gonic/gin"
)

func configureITTestEnvironment() *gin.Engine {
	r := gin.New()
	c := setupControllerWithDatabaseBackend()

	r.GET("/collection/:collection", c.GetCollection)
	r.POST("/collection/:collection", c.AddToCollection)

	return r
}

func setupControllerWithDatabaseBackend() *Controller {
	dbConfigs := configs.DBConfiguration{
		Type:     "mongodb",
		Host:     "localhost:27017",
		Database: "card-keeper-it",
	}

	return setupController(dbConfigs)
}

// TestAddNewCardToRepoIntegration verifies behaviour for adding a new card.
func TestAddNewCardToRepoIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	engine := configureITTestEnvironment()

	b := getSerializedTestCard()
	recordedPOSTResponse, err := makeAddCardRequestToHTTPServer(b, engine)

	if err != nil {
		t.Error("Failed to make POST request")
		t.Fail()
	}

	// Verify that the POST request succeeded with HTTP 200
	if recordedPOSTResponse.Code != http.StatusAccepted {
		t.Errorf("Expected to get HTTP 202, but got %d", recordedPOSTResponse.Code)
		t.Fail()
	}

	// Verify message body gives `ok`
	expected := `{"message":"ok"}`
	verifyHTTPResponseBody(expected, recordedPOSTResponse.Body.String(), t)

	recordedGETResponse, err := makeGetCardsRequesttToHTTPServer(engine)

	if err != nil {
		t.Error("Failed to make GET request")
		t.Fail()
	}

	if recordedGETResponse.Code != http.StatusOK {
		t.Errorf("Expected to get HTTP 200, but got %d", recordedGETResponse.Code)
		t.Fail()
	}

	response := getCollectionResponse{}

	err = json.Unmarshal(recordedGETResponse.Body.Bytes(), &response)

	if err != nil {
		t.Errorf("Failed to unmarshal get response %s", err)
		t.Fail()
	}

	if len(response.Cards) != 1 {
		t.Errorf("Expected to find one card returned after POST but got %d", len(response.Cards))
		t.Fail()
	}
}

// TestDuplicateCardNotAddedIntegration verifies that duplicates are not added to collection
func TestDuplicateCardNotAddedIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	b := getSerializedTestCard()

	engine := configureITTestEnvironment()

	recordedPOSTResponse, err := makeAddCardRequestToHTTPServer(b, engine)

	if err != nil {
		t.Error("Error making POST request")
		t.Fail()
	}

	// Verify that the POST request failed with 409
	if recordedPOSTResponse.Code != http.StatusConflict {
		t.Errorf("Expected to get HTTP 409, but got %d", recordedPOSTResponse.Code)
		t.Fail()
	}

	// Verify message body gives `ok`
	expected := `{"message":"duplicate item"}`
	verifyHTTPResponseBody(expected, recordedPOSTResponse.Body.String(), t)
}
