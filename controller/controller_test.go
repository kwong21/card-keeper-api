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

// AuthReponse is the struct wrapper for the Auth0 auth request
type AuthResponse struct {
	Token string `json:"access_token"`
	Type  string `json:"token_type"`
}

// TestAddsNewCardToRepo verifies behaviour for adding a new card.
func TestAddsNewCardToRepo(t *testing.T) {
	b := getSerializedTestCard()
	s := getTestService()

	req, err := http.NewRequest("POST", "/collection", bytes.NewBuffer(b))

	w := sendTestRequest(req, s)

	if err != nil {
		checkError(err, t)
	}

	// Verify that the POST request succeeded with HTTP 200
	if w.Code != http.StatusAccepted {
		t.Errorf("Expected to get HTTP 200, but got %d", w.Code)
		t.Fail()
	}

	// Verify message body gives `ok`
	expected := `{"message":"ok"}`
	verifyHTTPResponseBody(expected, w.Body.String(), t)

	// Verify that the card is in the Stored
	cards := s.GetAll()
	if len(*cards) != 1 {
		t.Errorf("Expected to get 1 card, but got %d", len(*cards))
		t.Fail()
	}
}

// TestAddNewCardError expects error to be returned when data is incorrect
func TestAddNewCardError(t *testing.T) {
	b, _ := json.Marshal("foo bar")
	s := getTestService()

	req, err := http.NewRequest("POST", "/collection", bytes.NewBuffer(b))

	w := sendTestRequest(req, s)

	checkError(err, t)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected to get HTTP 400, but got %d", w.Code)
	}

	// Verify message body gives `ok`
	expected := `{"message":"invalid data"}`
	verifyHTTPResponseBody(expected, w.Body.String(), t)
}

// TestPing checks that the API is responding
func TestPing(t *testing.T) {
	s := getTestService()

	req, err := http.NewRequest("GET", "/ping", nil)

	w := sendTestRequest(req, s)

	checkError(err, t)

	if w.Code != http.StatusOK {
		t.Errorf("Expected 200 but got %d", w.Code)
	}
}

// TestErrorOnDuplicate verifies an error is returned when a duplicate card is added.
func TestErrorOnDuplicate(t *testing.T) {
	b := getSerializedTestCard()
	s := getTestService()

	req, err := http.NewRequest("POST", "/collection", bytes.NewBuffer(b))

	w := sendTestRequest(req, s)

	if err != nil {
		checkError(err, t)
	}

	if w.Code != http.StatusAccepted {
		t.Errorf("Expected 202 but got %d", w.Code)
	}

	req, err = http.NewRequest("POST", "/collection", bytes.NewBuffer(b))
	w = sendTestRequest(req, s)

	if err != nil {
		checkError(err, t)
	}

	if w.Code != http.StatusConflict {
		t.Errorf("Expected 409, but got %d", w.Code)
	}
}

func sendTestRequest(req *http.Request, service service.Service) *httptest.ResponseRecorder {
	w, r, c := setupTestControllerAndHTTPRecorder()

	r.POST("/collection", c.AddToCollection)
	r.GET("/ping", c.Ping)

	c.Service = &service

	r.ServeHTTP(w, req)

	return w
}

func getTestService() service.Service {
	repo, _ := service.InMemoryStore()
	s := service.Service{
		Repository: repo,
	}

	return s
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

func checkError(err error, t *testing.T) {
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
}

// // TestAuthenticatedAPICall verifies a 200 is received if authorized
// func TestAuthenticatedAPICall(t *testing.T) {
// 	req, err := http.NewRequest("GET", "/ping", nil)

// 	req.Header.Set("authorization", getBearerTokenForTest())

// 	checkError(err, t)

// 	responseCode := testCheckJWTRequests(req)

// 	if responseCode != http.StatusOK {
// 		t.Errorf("Expected to get HTTP 200, but got %d", responseCode)
// 	}
// }

// func testCheckJWTRequests(req *http.Request) int {
// 	w, r, c := setupTestControllerAndHTTPRecorder()

// 	r.GET("/ping", checkJWT(), c.Ping)

// 	r.ServeHTTP(w, req)

// 	return w.Code
// }

// func getBearerTokenForTest() string {
// 	auth := new(AuthResponse)

// 	url := os.Getenv("AUTH0_URL")

// 	clientID := os.Getenv("AUTH0_CLIENT_ID")
// 	secretID := os.Getenv("AUTH0_SECRET_ID")
// 	audience := os.Getenv("AUTH0_AUDIENCE")

// 	payload := strings.NewReader(
// 		"{\"client_id\":" + "\"" + clientID + "\"" +
// 			",\"client_secret\":" + "\"" + secretID + "\"" +
// 			",\"audience\":" + "\"" + audience + "\"" +
// 			",\"grant_type\":\"client_credentials\"}")

// 	req, _ := http.NewRequest("POST", url, payload)

// 	req.Header.Add("content-type", "application/json")

// 	res, err := http.DefaultClient.Do(req)

// 	if err != nil {
// 		panic(err)
// 	}

// 	defer res.Body.Close()
// 	body, _ := ioutil.ReadAll(res.Body)

// 	err = json.Unmarshal(body, auth)

// 	if err != nil {
// 		panic(err)
// 	}

// 	return auth.Type + " " + auth.Token
// }
