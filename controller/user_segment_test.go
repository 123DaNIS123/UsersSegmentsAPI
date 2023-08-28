package controller

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/assert/v2"
)

func TestBind(t *testing.T) {
	router := gin.Default()
	router.POST("/bind", GetSegments)
	w := httptest.NewRecorder()

	jsonData := []byte(`{
		"segments_add1": "["AVITO_SEGMENT_1", "AVITO_SEGMENT_2"]",
        "segments_remove1": "["AVITO_SEGMENT_3", "AVITO_SEGMENT_4"]",
        "user_id1": 1
	}`)

	// Create a new HTTP request with the POST method and "/bind" URL
	req, err := http.NewRequest(http.MethodPost, "/bind", bytes.NewBuffer(jsonData))
	if err != nil {
		t.Fatalf("error creating request: %s", err.Error())
	}
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
}
