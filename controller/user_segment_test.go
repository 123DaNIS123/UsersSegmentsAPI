package controller

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/123DaNIS123/UsersSegments/db"
	"github.com/123DaNIS123/UsersSegments/models"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/assert/v2"
	"github.com/joho/godotenv"
)

func InitDB() {
	if err := godotenv.Load("../.env"); err != nil {
		log.Fatalf("error loading env variables %s", err.Error())
	}
	db.Connect()
}

func InitRouter(method string, url string, body io.Reader) *httptest.ResponseRecorder {
	r := gin.Default()
	r.GET(url, GetBinds)
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		log.Fatalf("error creating request %s", err.Error())
	}
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}

func TestGetBinds(t *testing.T) {
	InitDB()
	mockResponse := []models.UserSegment{}
	db.DB.Find(&mockResponse)
	mockResponseJSON, err := json.Marshal(mockResponse)
	if err != nil {
		log.Fatalf("error when marshaling %s", err.Error())
	}
	w := InitRouter("/binds", http.MethodGet, nil)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, w.Body.String(), string(mockResponseJSON))
}

func TestGetUserBinds(t *testing.T) {
	InitDB()
	mockResponse := []models.UserSegment{}
	data := UserRequest{ID: 1}
	var buf bytes.Buffer
	err := json.NewEncoder(&buf).Encode(data)
	if err != nil {
		log.Fatal(err)
	}
	db.DB.Find(&mockResponse)
	mockResponseJSON, err := json.Marshal(mockResponse)
	if err != nil {
		log.Fatalf("error when marshaling %s", err.Error())
	}
	w := InitRouter("/userbinds", http.MethodPost, &buf)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, w.Body.String(), string(mockResponseJSON))
}
