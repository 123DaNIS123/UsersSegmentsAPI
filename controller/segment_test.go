package controller

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"testing"

	"github.com/123DaNIS123/UsersSegments/db"
	"github.com/123DaNIS123/UsersSegments/models"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/assert/v2"
)

func TestCreateSegment(t *testing.T) { //POST
	InitDB()
	router := gin.Default()
	InitRoutes(router)
	var jsonStr = []byte(`{"name": "SEGMENT_20"}`) // if passed segment not created returns 201
	w := MakeRequest(http.MethodPost, "/segment", bytes.NewReader(jsonStr), router, true)
	fmt.Printf("w.Body: %v\n", w.Body.String())
	assert.Equal(t, http.StatusCreated, w.Code)
}

func TestDeleteSegment(t *testing.T) { //DELETE
	InitDB()
	router := gin.Default()
	InitRoutes(router)
	var jsonStr = []byte(`{"name": "AVITO_SEGMENT_3"}`)
	w := MakeRequest(http.MethodDelete, "/segment", bytes.NewReader(jsonStr), router, true)
	fmt.Printf("w.Body: %v\n", w.Body.String())
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestGetSegments(t *testing.T) { //GET
	InitDB()
	router := gin.Default()
	InitRoutes(router)
	mockResponse := []models.Segment{}
	db.DB.Find(&mockResponse)
	mockResponseJSON, err := json.Marshal(mockResponse)
	if err != nil {
		log.Fatalf("error when marshaling %s", err.Error())
	}
	w := MakeRequest(http.MethodGet, "/segments", nil, router, true)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, w.Body.String(), string(mockResponseJSON))
}

func TestGetSegment(t *testing.T) { //GET
	InitDB()
	router := gin.Default()
	InitRoutes(router)
	var mockResponse models.Segment
	db.DB.Where("id = ?", 5).First(&mockResponse)
	mockResponseJSON, err := json.Marshal(mockResponse)
	if err != nil {
		log.Fatalf("error when marshaling mock response%s", err.Error())
	}
	var paramStr = "id=5"
	w := MakeRequest(http.MethodGet, "/segment/5", strings.NewReader(paramStr), router, false)
	fmt.Printf("w.Body: %v\n", w.Body.String())
	fmt.Printf("mockResponseJSON: %v\n", string(mockResponseJSON))
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, w.Body.String(), string(mockResponseJSON))
}
