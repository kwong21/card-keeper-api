package config

import "testing"

func TestReturnDefaultConfig(t *testing.T) {
	c := Default()

	testDefault(c, t)
}

func TestReturnDefaultIfCantReadConfig(t *testing.T) {
	c := NewFromFile("../nofile.json")

	testDefault(c, t)
}

func TestReturnValuesInFile(t *testing.T) {
	c := NewFromFile("../config-example.json")

	fail := false

	if c.DBConfigs().Type != "mongodb" {
		t.Errorf("Expected mongodb, but got %s", c.DBConfigs().Type)
		fail = true
	}

	if c.APIAllowedOrigin() != "http://localhost:4200" {
		t.Errorf("Expected http://localhost:4200, but got %s", c.APIAllowedOrigin())
		fail = true
	}

	if c.APIListenPort() != "8080" {
		t.Errorf("Expected 8080, but got %s", c.APIListenPort())
		fail = true
	}

	if fail {
		t.Fail()
	}
}

func testDefault(c Configuration, t *testing.T) {
	fail := false

	if c.DBConfigs().Type != "in-memory" {
		t.Errorf("Expected in-memory, but got %s", c.DBConfigs().Type)
		fail = true
	}

	if c.APIAllowedOrigin() != "http://localhost:4200" {
		t.Errorf("Expected localhost:4200, but got %s", c.APIAllowedOrigin())
		fail = true
	}

	if fail {
		t.Fail()
	}
}
