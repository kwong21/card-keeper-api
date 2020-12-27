package controller

import (
	"bytes"
	"card-keeper-api/config"
	"card-keeper-api/service"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

// TestUnauthorizedResponse checks that 401 is returned when JWT token is missing
func TestUnauthorizedResponse(t *testing.T) {
	config := config.Default()
	r := InitServer(config)
	w := httptest.NewRecorder()
	b, _ := json.Marshal("foo bar")

	req, err := http.NewRequest("POST", "/v1/collection", bytes.NewBuffer(b))

	if err != nil {
		checkError(err, t)
	}

	r.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("Expected 401, got %d", w.Code)
	}
}

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

func initalizeRepositoryHelper(dbType string) (service.Repository, error) {
	fixture := config.DBConfiguration{
		Type: dbType,
	}

	return initializeRepository(fixture)
}
