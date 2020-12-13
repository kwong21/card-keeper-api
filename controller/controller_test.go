package controller

import (
	"bytes"
	"card-keeper-api/model"
	"card-keeper-api/service"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
)

// TestAddsNewCardToRepo verifies behaviour for adding a new card.
func TestAddsNewCardToRepo(t *testing.T) {
	w, r, c := setupTestControllerAndHTTPRecorder()

	repo := service.InMemoryStore()
	s := service.Service{
		Repository: repo,
	}

	c.Service = &s

	r.POST("/collection", c.AddToCollection)

	b := getSerializedTestCard()
	req, err := http.NewRequest("POST", "/collection", bytes.NewBuffer(b))

	checkError(err, t)

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

// TestAddNewCardError expects error to be returned when data is incorrect
func TestAddNewCardError(t *testing.T) {
	w, r, c := setupTestControllerAndHTTPRecorder()
	r.POST("/collection", c.AddToCollection)

	b, _ := json.Marshal("foo bar")
	req, err := http.NewRequest("POST", "/collection", bytes.NewBuffer(b))

	checkError(err, t)

	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected to get HTTP 400, but got %d", w.Code)
	}

	// Verify message body gives `ok`
	expected := `{"message":"invalid data"}`
	verifyHTTPResponseBody(expected, w.Body.String(), t)
}

// TestUnAuthenticatedAPICall verifies a 401 is received if unauthorized
func TestUnAuthenticatedAPICall(t *testing.T) {
	req, err := http.NewRequest("GET", "/ping", nil)

	checkError(err, t)

	responseCode := testCheckJWTRequests(req)

	if responseCode != http.StatusUnauthorized {
		t.Errorf("Expected to get HTTP 401, but got %d", responseCode)
	}
}

// TestAuthenticatedAPICall verifies a 200 is received if authorized
func TestAuthenticatedAPICall(t *testing.T) {
	req, err := http.NewRequest("GET", "/ping", nil)

	req.Header.Set("authorization", "Bearer eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCIsImtpZCI6IlhDeE53XzE1WFh1ZkItRk5vQ2FVVSJ9.eyJpc3MiOiJodHRwczovL2Rldi1zaGliYXRlay51cy5hdXRoMC5jb20vIiwic3ViIjoiVFVXaEFIQ2hNT0dUU0dYZlVEMEdwRWpOV3Bsc0lNMXBAY2xpZW50cyIsImF1ZCI6Imh0dHBzOi8vY2FyZGtlZXBlci1kZXYvYXBpIiwiaWF0IjoxNjA3ODk4MzQyLCJleHAiOjE2MDc5ODQ3NDIsImF6cCI6IlRVV2hBSENoTU9HVFNHWGZVRDBHcEVqTldwbHNJTTFwIiwiZ3R5IjoiY2xpZW50LWNyZWRlbnRpYWxzIn0.MMaJXDyZJYVP7WvBhLCrWJd6VfkykI-kgQ02Ra65aRYuqKfY2zwqeIam_dZHYAG0JyZyIQl6nE_AHEXIjwpGKNLynuFT9eHPaP3QOLI3FYDS0a8pgjOY0bCvTnRGTSWJn1Z93HIlZoX7-E6KbARpb0t-H-1_CaxbkDAptB7g5eQBklJR9ZpESePZ9t6cSgh0bF1n2CDoeAAXg-VW9sR4jJ0LxVL1EkMFMNyPchjaSgk8HDBcVWkoV1ZwvvdNc__LofeIjSERHSaIBRVCqj85PuUCR2TVOXEcuxP5h01Ehp4oR48fO2jeOGLZMxGVXh62vpicPFp5bfE6w0mN-xCqHA")

	checkError(err, t)

	responseCode := testCheckJWTRequests(req)

	if responseCode != http.StatusOK {
		t.Errorf("Expected to get HTTP 200, but got %d", responseCode)
	}
}

func testCheckJWTRequests(req *http.Request) int {
	w, r, c := setupTestControllerAndHTTPRecorder()

	// Set dev testing Auth0 instance
	os.Setenv("AUTH0_AUDIENCE", "https://cardkeeper-dev/api")
	os.Setenv("AUTH0_ISSUER", "https://dev-shibatek.us.auth0.com/")

	r.GET("/ping", checkJWT(), c.Ping)

	r.ServeHTTP(w, req)

	return w.Code
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

func getSerializedTestCard() []byte {
	card := model.Card{
		Year:   2020,
		Maker:  "Upper Deck",
		Set:    "Series One",
		Player: "Brock Boeser",
	}

	b, _ := json.Marshal(card)

	return b
}

func checkError(err error, t *testing.T) {
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
}
