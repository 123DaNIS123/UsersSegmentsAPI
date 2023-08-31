package controller

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/123DaNIS123/UsersSegments/db"
	"github.com/123DaNIS123/UsersSegments/models"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/assert/v2"
)

func TestGetSegments(t *testing.T) {
	t.Skip("skipping")
	router := gin.Default()
	router.GET("/segments", GetSegments)
	req, err := http.NewRequest("GET", "/segments", nil)
	if err != nil {
		t.Fatalf("error creating request: %s", err.Error())
	}
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)

	// expected := []models.Segment{}
	// config.DB.Find(&expected)
	// assert.Equal(t, expected, w.Body.Bytes())
}

func TestGetSegment(t *testing.T) {
	t.Skip("Skipping testing in CI environment")
	router := gin.Default()
	router.GET("/segment/:id", GetSegments)
	w := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "/segment/1", nil)
	if err != nil {
		t.Fatalf("error creating request: %s", err.Error())
	}
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)

	expected := []models.Segment{}
	db.DB.First(&expected, 1)
	assert.Equal(t, expected, w.Body.Bytes())
}
