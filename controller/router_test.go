package controller

import (
	"card-keeper-api/internal/configs"
	"card-keeper-api/middleware"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
)

// TestInitializeRepositoryShouldReturnErrorWhenDBTypeNotSupported
// Checks an error is returned when the DB Type is not configured
func TestInitializeRepositoryShouldReturnErrorWhenDBTypeNotSupported(t *testing.T) {
	_, err := initalizeRepositoryHelper("mysql")

	if err == nil {
		t.Errorf("Expected an error when db type is not configured")
	}
}

// TestInitializeRepositoryShouldNotReturnErrorForSupportedTypes
// Checks that the in-memory repo is returned
func TestInitializeRepositoryShouldNotReturnErrorForSupportedTypes(t *testing.T) {
	supported := []string{"in-memory"}

	for _, v := range supported {
		_, err := initalizeRepositoryHelper(v)

		if err != nil {
			t.Errorf("Expected no errors for supported db type, but got %s", err)
		}
	}
}

// TestAuthenticatedAPICallIntegration verifies a 200 is received if authorized
func TestAuthenticatedAPICallIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	req, err := http.NewRequest("GET", "/ping", nil)

	req.Header.Set("authorization", getBearerTokenForTest())

	if err != nil {
		t.Error("Error while creating HTTP Request", err)
	}

	responseCode := testCheckJWTRequests(req)

	if responseCode != http.StatusOK {
		t.Errorf("Expected to get HTTP 200, but got %d", responseCode)
	}
}

// TestUnauthorizedResponseIntegration checks that 401 is returned when JWT token is missing
func TestUnauthorizedResponseIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	req, err := http.NewRequest("GET", "/ping", nil)

	if err != nil {
		t.Error("error generating request")
		t.Fail()
	}

	code := testCheckJWTRequests(req)

	if code != http.StatusUnauthorized {
		t.Errorf("Expected 401, got %d", code)
	}
}

func testCheckJWTRequests(req *http.Request) int {
	r := gin.New()
	c := setupControllerWithDatabaseBackend()

	httpResponseRecorder := httptest.NewRecorder()

	testAuthConfigs := configs.AuthConfiguration{}

	testAuthConfigs.Audience = os.Getenv("AUTH0_AUDIENCE")
	testAuthConfigs.Issuer = os.Getenv("AUTH0_ISSUER")
	testAuthConfigs.JWKS = os.Getenv("AUTH0_JWKS")

	testJWTMiddleware := middleware.GetJWTMiddleware(testAuthConfigs)

	r.GET("/ping", checkJWT(testJWTMiddleware), c.Ping)

	r.ServeHTTP(httpResponseRecorder, req)

	return httpResponseRecorder.Code
}
