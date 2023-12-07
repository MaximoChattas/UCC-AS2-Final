package controller

import (
	"Hotel/dto"
	"Hotel/service"
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

type testHotel struct{}

func init() {
	service.HotelService = testHotel{}
}

func (t testHotel) GetHotelById(id string) (dto.HotelDto, error) {

	if id == "000000000000000000000000" {
		return dto.HotelDto{}, errors.New("hotel not found")
	}

	return dto.HotelDto{}, nil
}

func (t testHotel) GetHotels() (dto.HotelsDto, error) {
	return dto.HotelsDto{}, nil
}

func (t testHotel) InsertHotel(hotelDto dto.HotelDto) (dto.HotelDto, error) {
	return hotelDto, nil
}

func (t testHotel) DeleteHotel(id string) (dto.HotelDto, error) {
	return dto.HotelDto{Id: id}, nil
}

func (t testHotel) UpdateHotel(hotelDto dto.HotelDto) (dto.HotelDto, error) {
	return hotelDto, nil
}

func TestInsertHotel(t *testing.T) {
	a := assert.New(t)

	r := gin.Default()
	r.POST("/hotel", InsertHotel)

	body := `{
		"id": "654cf68d807298d99186019f",
        "name": "Hotel Test",
        "room_amount": 10,
        "description": "Test hotel description",
		"city": "Test City",
        "street_name": "Test Street",
        "street_number": 123,
        "rate": 4.5
    }`

	req, err := http.NewRequest(http.MethodPost, "/hotel", strings.NewReader(body))
	if err != nil {
		t.Fatalf("New request failed: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	c, _ := gin.CreateTestContext(w)
	c.Request = req

	r.ServeHTTP(w, req)

	a.Equal(http.StatusCreated, w.Code)

	var response dto.HotelDto
	err = json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	a.Equal("654cf68d807298d99186019f", response.Id)
}

func TestGetHotelById_NotFound(t *testing.T) {

	a := assert.New(t)

	r := gin.Default()
	r.GET("/hotel/:id", GetHotelById)

	req, err := http.NewRequest(http.MethodGet, "/hotel/000000000000000000000000", nil)

	if err != nil {
		t.Fatalf("New request failed: %v", err)
	}

	w := httptest.NewRecorder()

	c, _ := gin.CreateTestContext(w)
	c.Request = req

	r.ServeHTTP(w, req)

	a.Equal(http.StatusNotFound, w.Code)

	expectedResponse := `{"error":"hotel not found"}`

	a.Equal(expectedResponse, w.Body.String())
}

func TestGetHotelById_Found(t *testing.T) {
	a := assert.New(t)

	r := gin.Default()
	r.GET("/hotel/:id", GetHotelById)

	req, err := http.NewRequest(http.MethodGet, "/hotel/654cf68d807298d99186019f", nil)

	if err != nil {
		t.Fatalf("New request failed: %v", err)
	}

	w := httptest.NewRecorder()

	c, _ := gin.CreateTestContext(w)
	c.Request = req

	r.ServeHTTP(w, req)

	a.Equal(http.StatusOK, w.Code)

	var response dto.HotelDto
	err = json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	expectedResponse := dto.HotelDto{}

	a.Equal(expectedResponse, response)
}

func TestGetHotels(t *testing.T) {
	a := assert.New(t)

	r := gin.Default()
	r.GET("/hotel", GetHotels)

	req, err := http.NewRequest(http.MethodGet, "/hotel", nil)

	if err != nil {
		t.Fatalf("New request failed: %v", err)
	}

	w := httptest.NewRecorder()

	c, _ := gin.CreateTestContext(w)
	c.Request = req

	r.ServeHTTP(w, req)

	a.Equal(http.StatusOK, w.Code)

	var response dto.HotelsDto
	err = json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	expectedResponse := dto.HotelsDto{}

	a.Equal(expectedResponse, response)
}

func TestDeleteHotel(t *testing.T) {
	a := assert.New(t)

	r := gin.Default()
	r.DELETE("/hotel/:id", DeleteHotel)

	req, err := http.NewRequest(http.MethodDelete, "/hotel/654cf68d807298d99186019f", nil)

	if err != nil {
		t.Fatalf("New request failed: %v", err)
	}

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	r.ServeHTTP(w, req)

	a.Equal(http.StatusOK, w.Code)

	expectedResponse := `{"message":"Hotel 654cf68d807298d99186019f deleted"}`

	a.Equal(expectedResponse, w.Body.String())
}

func TestUpdateHotel(t *testing.T) {
	a := assert.New(t)

	r := gin.Default()
	r.PUT("/hotel/:id", UpdateHotel)

	body := `{
		"id": "654cf68d807298d99186019f",
        "name": "Hotel Test",
        "room_amount": 10,
        "description": "Test hotel description",
		"city": "Test City",
        "street_name": "Test Street",
        "street_number": 123,
        "rate": 4.5
    }`

	req, err := http.NewRequest(http.MethodPut, "/hotel/654cf68d807298d99186019f", strings.NewReader(body))
	if err != nil {
		t.Fatalf("New request failed: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	c, _ := gin.CreateTestContext(w)
	c.Request = req

	r.ServeHTTP(w, req)

	a.Equal(http.StatusOK, w.Code)

	var response dto.HotelDto
	err = json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	expectedResponse := dto.HotelDto{
		Id:           "654cf68d807298d99186019f",
		Name:         "Hotel Test",
		RoomAmount:   10,
		Description:  "Test hotel description",
		City:         "Test City",
		StreetName:   "Test Street",
		StreetNumber: 123,
		Rate:         4.5,
	}

	a.Equal(expectedResponse, response)
}
