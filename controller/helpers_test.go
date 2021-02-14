package controller

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"card-keeper-api/cardservice"
	configs "card-keeper-api/internal/configs"

	"github.com/gin-gonic/gin"
)

var DefaultSupportedCollection = []string{"hockey"}

type getCollectionResponse struct {
	Message string             `json:"message"`
	Cards   []cardservice.Card `json:"cards"`
}

func setupControllerWithInMemoryBackend() *Controller {
	inmemory := configs.Default()

	return setupController(inmemory.DBConfigs())
}

func makeAddCardRequestToHTTPServer(serializedData []byte, engine *gin.Engine) (*httptest.ResponseRecorder, error) {
	req, err := http.NewRequest("POST", "/collection/hockey", bytes.NewBuffer(serializedData))

	if err != nil {
		return nil, err
	}

	httpResponseRecorder := httptest.NewRecorder()
	engine.ServeHTTP(httpResponseRecorder, req)

	return httpResponseRecorder, nil
}

func makeGetCardsRequesttToHTTPServer(engine *gin.Engine) (*httptest.ResponseRecorder, error) {
	req, err := http.NewRequest("GET", "/collection/hockey", nil)

	if err != nil {
		return nil, err
	}

	httpResponseRecorder := httptest.NewRecorder()
	engine.ServeHTTP(httpResponseRecorder, req)

	return httpResponseRecorder, nil
}

func getBearerTokenForTest() string {
	auth := new(AuthResponse)

	url := os.Getenv("AUTH0_JWKS")

	clientID := os.Getenv("AUTH0_CLIENT_ID")
	secretID := os.Getenv("AUTH0_SECRET_ID")
	audience := os.Getenv("AUTH0_AUDIENCE")

	payload := strings.NewReader(
		"{\"client_id\":" + "\"" + clientID + "\"" +
			",\"client_secret\":" + "\"" + secretID + "\"" +
			",\"audience\":" + "\"" + audience + "\"" +
			",\"grant_type\":\"client_credentials\"}")

	req, _ := http.NewRequest("POST", url, payload)

	req.Header.Add("content-type", "application/json")

	res, err := http.DefaultClient.Do(req)

	if err != nil {
		panic(err)
	}

	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)

	err = json.Unmarshal(body, auth)

	if err != nil {
		panic(err)
	}

	return auth.Type + " " + auth.Token
}

func verifyHTTPResponseBody(expected string, actual string, t *testing.T) {
	if actual != expected {
		t.Errorf("Expected message body of %v, but got %v", expected, actual)
	}
}

func getSerializedTestCard() []byte {
	base := cardservice.Base{
		Year:   "2020",
		Make:   "Upper Deck",
		Set:    "Series One",
		Player: "Brock Boeser",
	}

	insert := cardservice.Insert{}

	card := cardservice.Card{
		Base:   base,
		Insert: insert,
	}

	b, _ := json.Marshal(card)

	return b
}

func initalizeRepositoryHelper(dbType string) (cardservice.Repository, error) {
	fixture := configs.DBConfiguration{
		Type: dbType,
	}

	return initializeRepository(fixture)
}
