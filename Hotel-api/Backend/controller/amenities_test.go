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

type testAmenities struct{}

func init() {
	service.AmenityService = testAmenities{}
}

func (t testAmenities) InsertAmenity(amenityDto dto.AmenityDto) (dto.AmenityDto, error) {
	return amenityDto, nil
}

func (t testAmenities) GetAmenities() (dto.AmenitiesDto, error) {
	return dto.AmenitiesDto{}, nil
}

func (t testAmenities) DeleteAmenityById(id string) error {

	if id == "000000000000000000000000" {
		return errors.New("amenity not found")
	}
	return nil
}

func TestInsertAmenity(t *testing.T) {
	a := assert.New(t)

	r := gin.Default()
	r.POST("/amenity", InsertAmenity)

	body := `{
		"id": "654cf68d807298d99186019f",
        "name": "Amenity Test"
    }`

	req, err := http.NewRequest(http.MethodPost, "/amenity", strings.NewReader(body))
	if err != nil {
		t.Fatalf("New request failed: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	c, _ := gin.CreateTestContext(w)
	c.Request = req

	r.ServeHTTP(w, req)

	a.Equal(http.StatusCreated, w.Code)

	var response dto.AmenityDto
	err = json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	a.Equal("654cf68d807298d99186019f", response.Id)
}

func TestGetAmenities(t *testing.T) {
	a := assert.New(t)

	r := gin.Default()
	r.GET("/amenity", GetAmenities)

	req, err := http.NewRequest(http.MethodGet, "/amenity", nil)

	if err != nil {
		t.Fatalf("New request failed: %v", err)
	}

	w := httptest.NewRecorder()

	c, _ := gin.CreateTestContext(w)
	c.Request = req

	r.ServeHTTP(w, req)

	a.Equal(http.StatusOK, w.Code)

	var response dto.AmenitiesDto
	err = json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	expectedResponse := dto.AmenitiesDto{}

	a.Equal(expectedResponse, response)
}

func TestDeleteAmenityById_NotFound(t *testing.T) {
	a := assert.New(t)

	r := gin.Default()
	r.DELETE("/amenity/:id", DeleteAmenityById)

	req, err := http.NewRequest(http.MethodDelete, "/amenity/000000000000000000000000", nil)

	if err != nil {
		t.Fatalf("New request failed: %v", err)
	}

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	r.ServeHTTP(w, req)

	a.Equal(http.StatusBadRequest, w.Code)

	expectedResponse := `{"error":"amenity not found"}`

	a.Equal(expectedResponse, w.Body.String())
}

func TestDeleteAmenityById_Found(t *testing.T) {

	a := assert.New(t)

	r := gin.Default()
	r.DELETE("/amenity/:id", DeleteAmenityById)

	req, err := http.NewRequest(http.MethodDelete, "/amenity/654cf68d807298d99186019f", nil)

	if err != nil {
		t.Fatalf("New request failed: %v", err)
	}

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	r.ServeHTTP(w, req)

	a.Equal(http.StatusOK, w.Code)

	expectedResponse := `{"message":"amenity deleted successfully"}`

	a.Equal(expectedResponse, w.Body.String())

}
