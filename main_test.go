package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestMain(t *testing.T) {
	go main()
	time.Sleep(5 * time.Second)
	resp, err := http.Get("http://127.0.0.1:8080/")
	if err != nil {
		t.Fatalf("Cannot make get: %v\n", err)
	}
	bodySb, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("Error reading body: %v\n", err)
	}
	body := string(bodySb)
	fmt.Printf("Body: %v\n", body)
	var decodedResponse interface{}
	err = json.Unmarshal(bodySb, &decodedResponse)
	if err != nil {
		t.Fatalf("Cannot decode response <%p> from server. Err: %v", bodySb, err)
	}
	assert.Equal(t, map[string]interface{}{"message": "go to /docs/index.html for API docs"}, decodedResponse)
}
