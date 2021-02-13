package controller

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/gin-gonic/gin"
)

// AuthReponse is the struct wrapper for the Auth0 auth request
type AuthResponse struct {
	Token string `json:"access_token"`
	Type  string `json:"token_type"`
}

// TestAddsNewCardToRepo verifies behaviour for adding a new card.
func TestAddsNewCardToRepo(t *testing.T) {
	t.Parallel()

	testEngine := setupTestEnvironment()

	b := getSerializedTestCard()
	recordedPostResponse, err := makeAddCardRequestToHTTPServer(b, testEngine)

	if err != nil {
		t.Error("Got error when making request.")
		t.Fail()
	}

	// Verify that the POST request succeeded with HTTP 200
	if recordedPostResponse.Code != http.StatusAccepted {
		t.Errorf("Expected to get HTTP 200, but got %d", recordedPostResponse.Code)
		t.Fail()
	}

	// Verify message body gives `ok`
	expected := `{"message":"ok"}`
	verifyHTTPResponseBody(expected, recordedPostResponse.Body.String(), t)

	recordedGETResponse, err := makeGetCardsRequesttToHTTPServer(testEngine)

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

// TestErrorOnDuplicate verifies an error is returned when a duplicate card is added.
func TestErrorOnDuplicate(t *testing.T) {
	t.Parallel()

	testEngine := setupTestEnvironment()

	b := getSerializedTestCard()
	firstPostRequest, err := makeAddCardRequestToHTTPServer(b, testEngine)

	if err != nil {
		t.Errorf("error making request %s", err)
		t.Fail()
	}

	if firstPostRequest.Code != http.StatusAccepted {
		t.Errorf("Expected 202, but got %d", firstPostRequest.Code)
	}

	duplicateRequest, err := makeAddCardRequestToHTTPServer(b, testEngine)

	if err != nil {
		t.Errorf("error making request %s", err)
		t.Fail()
	}

	if duplicateRequest.Code != http.StatusConflict {
		t.Errorf("Expected 409, but got %d", duplicateRequest.Code)
	}

}

// TestAddNewCardError expects error to be returned when data is incorrect
func TestAddNewCardError(t *testing.T) {
	t.Parallel()

	testEngine := setupTestEnvironment()

	b, _ := json.Marshal("foo bar")

	recordedPostResponse, err := makeAddCardRequestToHTTPServer(b, testEngine)

	if err != nil {
		t.Errorf("error making request %s", err)
		t.Fail()
	}

	if recordedPostResponse.Code != http.StatusBadRequest {
		t.Errorf("Expected to get HTTP 400, but got %d", recordedPostResponse.Code)
	}

	// Verify message body gives `ok`
	expected := `{"message":"invalid data"}`
	verifyHTTPResponseBody(expected, recordedPostResponse.Body.String(), t)
}

func setupTestEnvironment() *gin.Engine {
	testEngine := gin.New()
	testController := setupControllerWithInMemoryBackend()

	testEngine.GET("/collection/:collection", testController.GetCollection)
	testEngine.POST("/collection/:collection", testController.AddToCollection)

	return testEngine
}
