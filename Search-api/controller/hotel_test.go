package controller

import (
	"Search/dto"
	"Search/service"
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

type testHotel struct{}

func init() {
	service.HotelService = testHotel{}
}

func (t testHotel) InsertUpdateHotel(hotelDto dto.HotelDto) error {
	return nil
}

func (t testHotel) GetHotels() (dto.HotelsDto, error) {
	return dto.HotelsDto{}, nil
}

func (t testHotel) GetHotelById(id string) (dto.HotelDto, error) {

	if id == "000000000000000000000000" {
		return dto.HotelDto{}, errors.New("hotel not found")
	}

	return dto.HotelDto{}, nil
}

func (t testHotel) GetHotelByCity(id string) (dto.HotelsDto, error) {
	return dto.HotelsDto{}, nil
}

func (t testHotel) DeleteHotelById(id string) error {
	return nil
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

func TestGetHotels_WithQuery(t *testing.T) {
	a := assert.New(t)

	r := gin.Default()
	r.GET("/hotel", GetHotels)

	req, err := http.NewRequest(http.MethodGet, "/hotel?city=testCity", nil)

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
