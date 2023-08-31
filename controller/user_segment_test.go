package controller

import (
	"bytes"
	"encoding/json"
	"fmt"
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

func InitRoutes(router *gin.Engine) {
	router.GET("/", Index)

	router.GET("/users", GetUsers)
	router.POST("/user", CreateUser)
	router.DELETE("/user/:id", DeleteUser)
	router.PUT("/user/:id", UpdateUser)

	router.GET("/segments", GetSegments)
	router.GET("/segment/:id", GetSegment)
	router.POST("/segment", CreateSegment)
	router.DELETE("/segment", DeleteSegment)
	router.PUT("/segment/:id", UpdateSegment)

	router.POST("/bind", Bind)
	router.GET("/binds", GetBinds)
	router.POST("/userbinds", GetUserBinds)
	router.POST("/timedata", GetTimeData)
}

func MakeRequest(method string, url string, body io.Reader, router *gin.Engine, isJson bool) *httptest.ResponseRecorder {
	req, err := http.NewRequest(method, url, body)
	if isJson {
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	}
	fmt.Printf("w.Body: %v\n", req)
	if err != nil {
		log.Fatalf("error creating request %s", err.Error())
	}

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	return w
}

func TestBind(t *testing.T) { // POST
	InitDB()
	router := gin.Default()
	InitRoutes(router)
	var jsonStr = []byte(`{"segments_add":["SEGMENT_1", "SEGMENT_2", "SEGMENT_3"],
	"segments_remove":["SEGMENT_4", "SEGMENT_5", "SEGMENT_6"], "user_id": 10}`)
	data := BindMessage{}
	json.Unmarshal([]byte(jsonStr), &data)
	mockResponse := data.AddUserSegments()
	mockResponseJSON, err := json.Marshal(mockResponse)
	if err != nil {
		log.Fatalf("error when marshaling mock response%s", err.Error())
	}
	w := MakeRequest(http.MethodPost, "/bind", bytes.NewReader(jsonStr), router, true)
	fmt.Printf("w.Body: %v\n", w.Body.String())
	fmt.Printf("mockResponseJSON: %v\n", string(mockResponseJSON))
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, w.Body.String(), string(mockResponseJSON))
}

func TestGetBinds(t *testing.T) { //GET
	InitDB()
	router := gin.Default()
	InitRoutes(router)
	mockResponse := []models.UserSegment{}
	db.DB.Find(&mockResponse)
	mockResponseJSON, err := json.Marshal(mockResponse)
	if err != nil {
		log.Fatalf("error when marshaling %s", err.Error())
	}
	w := MakeRequest(http.MethodGet, "/binds", nil, router, true)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, w.Body.String(), string(mockResponseJSON))
}

func TestGetUserBinds(t *testing.T) { // POST
	InitDB()
	router := gin.Default()
	InitRoutes(router)
	var mockResponse []models.Segment
	if err := db.DB.Table("user_segments").
		Where("user_id = ? AND deleted_at IS NULL", 5).
		Order("segment_id asc").Joins("join segments on segments.id = user_segments.segment_id").
		Select("segments.id", "segments.name", "segments.percentage").
		Find(&mockResponse).Error; err != nil {
		fmt.Println("Wait for bad request")
	}
	mockResponseJSON, err := json.Marshal(mockResponse)
	if err != nil {
		log.Fatalf("error when marshaling mock response%s", err.Error())
	}
	var jsonStr = []byte(`{"user_id":5}`)
	w := MakeRequest(http.MethodPost, "/userbinds", bytes.NewReader(jsonStr), router, true)
	fmt.Printf("w.Body: %v\n", w.Body.String())
	fmt.Printf("mockResponseJSON: %v\n", string(mockResponseJSON))
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, w.Body.String(), string(mockResponseJSON))
}
