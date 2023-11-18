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

type reservationTest struct{}

func init() {
	service.ReservationService = reservationTest{}
}

func (t reservationTest) InsertReservation(reservationDto dto.ReservationDto) (dto.ReservationDto, error) {

	if reservationDto.Id == 0 {
		return dto.ReservationDto{}, errors.New("error inserting reservation")
	}

	return reservationDto, nil
}

func (t reservationTest) GetReservationById(id int) (dto.ReservationDto, error) {

	if id == 0 {
		return dto.ReservationDto{}, errors.New("reservation not found")
	}

	return dto.ReservationDto{Id: id}, nil
}

func (t reservationTest) GetReservations() (dto.ReservationsDto, error) {

	return dto.ReservationsDto{}, nil
}

func (t reservationTest) DeleteReservation(id int) error {

	if id == 0 {
		return errors.New("reservation not found")
	}

	return nil
}

func (t reservationTest) GetReservationsByUser(userId int) (dto.UserReservationsDto, error) {

	if userId == 0 {
		return dto.UserReservationsDto{}, errors.New("user not found")
	}

	return dto.UserReservationsDto{}, nil
}

func (t reservationTest) CheckAvailability(_ string, _ time.Time, _ time.Time) bool {

	return true
}

func (t reservationTest) CheckAllAvailability(city string, _ string, _ string) (dto.HotelsDto, error) {

	if city == "" {
		return dto.HotelsDto{}, errors.New("no city selected")
	}

	return dto.HotelsDto{}, nil
}

func (t reservationTest) GetAllHotelsByCity(city string) dto.HotelsDto {

	return dto.HotelsDto{}
}

func (t reservationTest) GetHotelInfo(hotelId string) (dto.HotelDto, error) {

	return dto.HotelDto{}, nil
}

func TestInsertReservation(t *testing.T) {

	a := assert.New(t)

	r := gin.Default()
	r.POST("/reserve", InsertReservation)

	body := `{
		"id": 1
   }`

	req, err := http.NewRequest(http.MethodPost, "/reserve", strings.NewReader(body))

	if err != nil {
		t.Fatalf("New request failed: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	c, _ := gin.CreateTestContext(w)
	c.Request = req

	r.ServeHTTP(w, req)

	a.Equal(http.StatusCreated, w.Code)

	var response dto.ReservationDto

	err = json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	expectedDto := dto.ReservationDto{Id: 1}

	a.Equal(expectedDto, response)
}

func TestInsertReservation_Error(t *testing.T) {

	a := assert.New(t)

	r := gin.Default()
	r.POST("/reserve", InsertReservation)

	body := `{
		"id": 0
   }`

	req, err := http.NewRequest(http.MethodPost, "/reserve", strings.NewReader(body))

	if err != nil {
		t.Fatalf("New request failed: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	c, _ := gin.CreateTestContext(w)
	c.Request = req

	r.ServeHTTP(w, req)

	a.Equal(http.StatusBadRequest, w.Code)

	expectedResponse := `{"error":"error inserting reservation"}`
	a.Equal(expectedResponse, w.Body.String())
}

func TestGetReservationById_NotFound(t *testing.T) {

	a := assert.New(t)

	r := gin.Default()
	r.GET("/reservation/:id", GetReservationById)

	req, err := http.NewRequest(http.MethodGet, "/reservation/0", nil)
	if err != nil {
		t.Fatalf("New request failed: %v", err)
	}

	w := httptest.NewRecorder()

	c, _ := gin.CreateTestContext(w)
	c.Request = req

	r.ServeHTTP(w, req)

	a.Equal(http.StatusNotFound, w.Code)

	expectedResponse := `{"error":"reservation not found"}`
	a.Equal(expectedResponse, w.Body.String())

}

func TestGetReservationById_Found(t *testing.T) {

	a := assert.New(t)

	r := gin.Default()
	r.GET("/reservation/:id", GetReservationById)

	req, err := http.NewRequest(http.MethodGet, "/reservation/1", nil)
	if err != nil {
		t.Fatalf("New request failed: %v", err)
	}

	w := httptest.NewRecorder()

	c, _ := gin.CreateTestContext(w)
	c.Request = req

	r.ServeHTTP(w, req)

	a.Equal(http.StatusOK, w.Code)

	var response dto.ReservationDto
	expectedResponse := dto.ReservationDto{Id: 1}

	err = json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	a.Equal(expectedResponse, response)
}

func TestGetReservations(t *testing.T) {

	a := assert.New(t)

	r := gin.Default()
	r.GET("/reservation", GetReservations)

	req, err := http.NewRequest(http.MethodGet, "/reservation", nil)
	if err != nil {
		t.Fatalf("New request failed: %v", err)
	}

	w := httptest.NewRecorder()

	c, _ := gin.CreateTestContext(w)
	c.Request = req

	r.ServeHTTP(w, req)

	a.Equal(http.StatusOK, w.Code)

	expectedResponse := dto.ReservationsDto{}
	var response dto.ReservationsDto

	err = json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	a.Equal(expectedResponse, response)
}

func TestDeleteReservation_NotFound(t *testing.T) {

	a := assert.New(t)

	r := gin.Default()
	r.DELETE("/reservation/:id", DeleteReservation)

	req, err := http.NewRequest(http.MethodDelete, "/reservation/0", nil)
	if err != nil {
		t.Fatalf("New request failed: %v", err)
	}

	w := httptest.NewRecorder()

	c, _ := gin.CreateTestContext(w)
	c.Request = req

	r.ServeHTTP(w, req)

	a.Equal(http.StatusBadRequest, w.Code)

	expectedResponse := `{"error":"reservation not found"}`
	a.Equal(expectedResponse, w.Body.String())
}

func TestDeleteReservation_Found(t *testing.T) {

	a := assert.New(t)

	r := gin.Default()
	r.DELETE("/reservation/:id", DeleteReservation)

	req, err := http.NewRequest(http.MethodDelete, "/reservation/1", nil)
	if err != nil {
		t.Fatalf("New request failed: %v", err)
	}

	w := httptest.NewRecorder()

	c, _ := gin.CreateTestContext(w)
	c.Request = req

	r.ServeHTTP(w, req)

	a.Equal(http.StatusOK, w.Code)

	expectedResponse := `{"message":"Reservation deleted"}`
	a.Equal(expectedResponse, w.Body.String())
}

func TestCheckAvailability_Error(t *testing.T) {

	a := assert.New(t)

	r := gin.Default()
	r.GET("/available", CheckAvailability)

	req, err := http.NewRequest(http.MethodGet, "/available", nil)
	if err != nil {
		t.Fatalf("New request failed: %v", err)
	}

	w := httptest.NewRecorder()

	c, _ := gin.CreateTestContext(w)
	c.Request = req

	r.ServeHTTP(w, req)

	a.Equal(http.StatusBadRequest, w.Code)

	expectedResponse := `{"error":"no city selected"}`
	a.Equal(expectedResponse, w.Body.String())
}

func TestCheckAvailability(t *testing.T) {

	a := assert.New(t)

	r := gin.Default()
	r.GET("/available", CheckAvailability)

	req, err := http.NewRequest(http.MethodGet, "/available?city=Cordoba", nil)
	if err != nil {
		t.Fatalf("New request failed: %v", err)
	}

	w := httptest.NewRecorder()

	c, _ := gin.CreateTestContext(w)
	c.Request = req

	r.ServeHTTP(w, req)

	a.Equal(http.StatusOK, w.Code)

	expectedResponse := dto.HotelsDto{}
	var response dto.HotelsDto

	err = json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("Failed to unmarhsal response: %v", err)
	}
	a.Equal(expectedResponse, response)
}

func TestGetReservationsByUser_NotFound(t *testing.T) {

	a := assert.New(t)

	r := gin.Default()
	r.GET("/user/reservations/:id", GetReservationsByUser)

	req, err := http.NewRequest(http.MethodGet, "/user/reservations/0", nil)
	if err != nil {
		t.Fatalf("New request failed: %v", err)
	}

	w := httptest.NewRecorder()

	c, _ := gin.CreateTestContext(w)
	c.Request = req

	r.ServeHTTP(w, req)

	a.Equal(http.StatusNotFound, w.Code)

	expectedResponse := `{"error":"user not found"}`
	a.Equal(expectedResponse, w.Body.String())
}

func TestGetReservationsByUser_Found(t *testing.T) {

	a := assert.New(t)

	r := gin.Default()
	r.GET("/user/reservations/:id", GetReservationsByUser)

	req, err := http.NewRequest(http.MethodGet, "/user/reservations/1", nil)
	if err != nil {
		t.Fatalf("New request failed: %v", err)
	}

	w := httptest.NewRecorder()

	c, _ := gin.CreateTestContext(w)
	c.Request = req

	r.ServeHTTP(w, req)

	a.Equal(http.StatusOK, w.Code)

	expectedResponse := dto.UserReservationsDto{}
	var response dto.UserReservationsDto

	err = json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	a.Equal(expectedResponse, response)
}
