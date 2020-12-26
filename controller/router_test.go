package controller

import (
	"bytes"
	"card-keeper-api/config"
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
