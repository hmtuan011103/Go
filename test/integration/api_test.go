package integration

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gostructure/app/internal/app"
	"github.com/gostructure/app/internal/config"
)

func setupTestApp() *app.App {
	cfg := &config.Config{
		App: config.AppConfig{
			Name:        "Test App",
			Version:     "1.0.0",
			Environment: "test",
			Debug:       true,
		},
	}
	return app.New(cfg)
}

func TestHealthEndpoint(t *testing.T) {
	application := setupTestApp()

	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	rec := httptest.NewRecorder()

	application.Router().ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, rec.Code)
	}

	var response map[string]string
	if err := json.NewDecoder(rec.Body).Decode(&response); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if response["status"] != "healthy" {
		t.Errorf("Expected status 'healthy', got '%s'", response["status"])
	}
}

func TestReadyEndpoint(t *testing.T) {
	application := setupTestApp()

	req := httptest.NewRequest(http.MethodGet, "/ready", nil)
	rec := httptest.NewRecorder()

	application.Router().ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, rec.Code)
	}
}

func TestInfoEndpoint(t *testing.T) {
	application := setupTestApp()

	req := httptest.NewRequest(http.MethodGet, "/api/v1/info", nil)
	rec := httptest.NewRecorder()

	application.Router().ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, rec.Code)
	}

	var response map[string]interface{}
	if err := json.NewDecoder(rec.Body).Decode(&response); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if response["name"] != "Test App" {
		t.Errorf("Expected name 'Test App', got '%s'", response["name"])
	}
}

func TestUserCRUD(t *testing.T) {
	application := setupTestApp()

	// Create user
	createBody := []byte(`{"name": "John Doe", "email": "john@example.com"}`)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/users", bytes.NewBuffer(createBody))
	rec := httptest.NewRecorder()

	application.Router().ServeHTTP(rec, req)

	if rec.Code != http.StatusCreated {
		t.Errorf("Create user: Expected status %d, got %d", http.StatusCreated, rec.Code)
	}

	var createdUser map[string]interface{}
	if err := json.NewDecoder(rec.Body).Decode(&createdUser); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	userID := int(createdUser["id"].(float64))

	// Get user
	req = httptest.NewRequest(http.MethodGet, "/api/v1/users/"+itoa(userID), nil)
	rec = httptest.NewRecorder()

	application.Router().ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("Get user: Expected status %d, got %d", http.StatusOK, rec.Code)
	}

	// Update user
	updateBody := []byte(`{"name": "Jane Doe"}`)
	req = httptest.NewRequest(http.MethodPut, "/api/v1/users/"+itoa(userID), bytes.NewBuffer(updateBody))
	rec = httptest.NewRecorder()

	application.Router().ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("Update user: Expected status %d, got %d", http.StatusOK, rec.Code)
	}

	// Delete user
	req = httptest.NewRequest(http.MethodDelete, "/api/v1/users/"+itoa(userID), nil)
	rec = httptest.NewRecorder()

	application.Router().ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("Delete user: Expected status %d, got %d", http.StatusOK, rec.Code)
	}
}

func itoa(i int) string {
	return string(rune('0' + i))
}
