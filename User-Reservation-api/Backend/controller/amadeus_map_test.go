package controller

import (
	"User-Reservation/dto"
	"User-Reservation/service"
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

type testAmadeus struct{}

func init() {

	service.AmadeusService = &testAmadeus{}
}

func (t testAmadeus) InsertAmadeusMap(amadeusMapDto dto.AmadeusMapDto) (dto.AmadeusMapDto, error) {

	if amadeusMapDto.HotelId == "" || amadeusMapDto.AmadeusId == "" {
		return dto.AmadeusMapDto{}, errors.New("failed to insert amadeus map")
	}

	return amadeusMapDto, nil
}

func (t testAmadeus) GetAmadeusIdByHotelId(hotelId string) (dto.AmadeusMapDto, error) {

	if hotelId == "0" {
		return dto.AmadeusMapDto{}, errors.New("hotel not found")
	}

	return dto.AmadeusMapDto{hotelId, "SBMIASOF"}, nil
}

func (t testAmadeus) GetAmadeusAvailability(amadeusId string, _ time.Time, _ time.Time) (bool, error) {

	if amadeusId == "0" {
		return false, errors.New("amadeus id not found")
	}

	if amadeusId == "1" {
		return false, nil
	}

	return true, nil
}

func (t testAmadeus) DeleteMapping(hotelId string) error {

	if hotelId == "0" {
		return errors.New("mapping not found")
	}

	return nil
}

func TestInsertAmadeusMap(t *testing.T) {

	a := assert.New(t)

	r := gin.Default()
	r.POST("/amadeus", InsertAmadeusMap)

	hotelId := "654cf68d807298d99186019f"
	amadeusId := "SBMIASOF"

	body := `{
    	"hotel_id": "` + hotelId + `",
    	"amadeus_id": "` + amadeusId + `"
	}`

	req, err := http.NewRequest(http.MethodPost, "/amadeus", strings.NewReader(body))
	if err != nil {
		t.Fatalf("New request failed: %v", err)
	}

	w := httptest.NewRecorder()

	c, _ := gin.CreateTestContext(w)
	c.Request = req

	r.ServeHTTP(w, req)
	a.Equal(http.StatusCreated, w.Code)

	var response dto.AmadeusMapDto
	expectedResponse := dto.AmadeusMapDto{hotelId, amadeusId}

	err = json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	a.Equal(expectedResponse, response)
}

func TestInsertAmadeusMap_Error(t *testing.T) {

	a := assert.New(t)

	r := gin.Default()
	r.POST("/amadeus", InsertAmadeusMap)

	hotelId := ""
	amadeusId := ""

	body := `{
    	"hotel_id": "` + hotelId + `",
    	"amadeus_id": "` + amadeusId + `"
	}`

	req, err := http.NewRequest(http.MethodPost, "/amadeus", strings.NewReader(body))
	if err != nil {
		t.Fatalf("New request failed: %v", err)
	}

	w := httptest.NewRecorder()

	c, _ := gin.CreateTestContext(w)
	c.Request = req

	r.ServeHTTP(w, req)
	a.Equal(http.StatusBadRequest, w.Code)

	expectedResponse := `{"error":"failed to insert amadeus map"}`

	a.Equal(expectedResponse, w.Body.String())
}

func TestGetAmadeusIdByHotelId_NotFound(t *testing.T) {

	a := assert.New(t)

	r := gin.Default()
	r.GET("/amadeus/:hotel_id", GetAmadeusIdByHotelId)

	req, err := http.NewRequest(http.MethodGet, "/amadeus/0", nil)
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

func TestGetAmadeusIdByHotelId_Found(t *testing.T) {

	a := assert.New(t)

	r := gin.Default()
	r.GET("/amadeus/:hotel_id", GetAmadeusIdByHotelId)

	req, err := http.NewRequest(http.MethodGet, "/amadeus/1", nil)
	if err != nil {
		t.Fatalf("New request failed: %v", err)
	}

	w := httptest.NewRecorder()

	c, _ := gin.CreateTestContext(w)
	c.Request = req

	r.ServeHTTP(w, req)

	a.Equal(http.StatusOK, w.Code)

	var response dto.AmadeusMapDto
	expectedResponse := dto.AmadeusMapDto{"1", "SBMIASOF"}

	err = json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}
	a.Equal(expectedResponse, response)
}

func TestDeleteMapping_NotFound(t *testing.T) {

	a := assert.New(t)

	r := gin.Default()
	r.DELETE("/amadeus/:hotel_id", DeleteMapping)

	req, err := http.NewRequest(http.MethodDelete, "/amadeus/0", nil)
	if err != nil {
		t.Fatalf("New request failed: %v", err)
	}

	w := httptest.NewRecorder()

	c, _ := gin.CreateTestContext(w)
	c.Request = req

	r.ServeHTTP(w, req)

	a.Equal(http.StatusBadRequest, w.Code)

	expectedResponse := `{"error":"mapping not found"}`
	a.Equal(expectedResponse, w.Body.String())
}

func TestDeleteMapping_Found(t *testing.T) {

	a := assert.New(t)

	r := gin.Default()
	r.DELETE("/amadeus/:hotel_id", DeleteMapping)

	req, err := http.NewRequest(http.MethodDelete, "/amadeus/654cf68d807298d99186019f", nil)
	if err != nil {
		t.Fatalf("New request failed: %v", err)
	}

	w := httptest.NewRecorder()

	c, _ := gin.CreateTestContext(w)
	c.Request = req

	r.ServeHTTP(w, req)

	a.Equal(http.StatusOK, w.Code)

	expectedResponse := `{"message":"Mapping deleted"}`
	a.Equal(expectedResponse, w.Body.String())
}
