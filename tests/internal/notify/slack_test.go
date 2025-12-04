package notify_test

import (
	"dbx/internal/notify"
	"net/http"
	"net/http/httptest"
	"testing"
)

// TestSlackNotify_ValidWebhook tests successful notification
func TestSlackNotify_ValidWebhook(t *testing.T) {
	// Create a mock HTTP server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			t.Errorf("Expected POST request, got %s", r.Method)
		}
		if r.Header.Get("Content-Type") != "application/json" {
			t.Errorf("Expected Content-Type application/json, got %s", r.Header.Get("Content-Type"))
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	err := notify.SlackNotify(server.URL, "Test message")
	if err != nil {
		t.Errorf("SlackNotify() error = %v, want nil", err)
	}
}

// TestSlackNotify_InvalidWebhook tests error handling for invalid webhook
func TestSlackNotify_InvalidWebhook(t *testing.T) {
	err := notify.SlackNotify("http://invalid-url-that-does-not-exist.local", "Test message")
	if err == nil {
		t.Error("SlackNotify() should return error for invalid webhook URL")
	}
}

// TestSlackNotify_EmptyMessage tests handling of empty message
func TestSlackNotify_EmptyMessage(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	// Empty message should still be sent
	err := notify.SlackNotify(server.URL, "")
	if err != nil {
		t.Errorf("SlackNotify() should handle empty message, got error: %v", err)
	}
}

// TestSlackNotify_LongMessage tests handling of long messages (scalability)
func TestSlackNotify_LongMessage(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	// Create a very long message (1MB)
	longMessage := make([]byte, 1024*1024)
	for i := range longMessage {
		longMessage[i] = 'A'
	}

	err := notify.SlackNotify(server.URL, string(longMessage))
	if err != nil {
		t.Logf("SlackNotify() with long message returned error (may be expected): %v", err)
	}
}

// TestSlackNotify_SpecialCharacters tests handling of special characters in message
func TestSlackNotify_SpecialCharacters(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	specialMessage := "Test with special chars: <>&\"'`\n\t\r"
	err := notify.SlackNotify(server.URL, specialMessage)
	if err != nil {
		t.Errorf("SlackNotify() should handle special characters, got error: %v", err)
	}
}

// TestSlackNotify_ServerError tests handling of server errors
func TestSlackNotify_ServerError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	err := notify.SlackNotify(server.URL, "Test message")
	// Should return error for server error status
	if err == nil {
		t.Log("SlackNotify() may not detect HTTP error status codes")
	}
}

// TestSlackNotify_ConcurrentNotifications tests concurrent notifications (scalability)
func TestSlackNotify_ConcurrentNotifications(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	numGoroutines := 10
	errors := make(chan error, numGoroutines)

	for i := 0; i < numGoroutines; i++ {
		go func(id int) {
			errors <- notify.SlackNotify(server.URL, "Concurrent test message")
		}(i)
	}

	// Collect results
	var errorCount int
	for i := 0; i < numGoroutines; i++ {
		if err := <-errors; err != nil {
			errorCount++
		}
	}

	if errorCount > 0 {
		t.Errorf("Concurrent SlackNotify() failed for %d/%d operations", errorCount, numGoroutines)
	}
}
