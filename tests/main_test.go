package main

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHttpServer(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello World"))
	})

	server := httptest.NewServer(mux)
	defer server.Close()

	resp, err := http.Get(server.URL)
	if err != nil {
		t.Fatalf("Could not send GET request: %v", err)
	}

	if resp.StatusCode != 200 {
		t.Fatalf("Expected status 200, got %v", resp.StatusCode)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("Could not read response: %v", err)
	}

	if string(body) != "Hello World" {
		t.Fatalf("Expected 'Hello World', got '%s'", body)
	}
}
