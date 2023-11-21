package service

import (
	"User-Reservation/client"
	"User-Reservation/dto"
	"User-Reservation/model"
	"errors"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"strings"
	"testing"
	"time"
)

type testReservation struct{}
type mockHTTPClient struct{}
type mockCache struct{}

func init() {
	client.ReservationClient = &testReservation{}
	ReservationService = &reservationService{HTTPClient: &mockHTTPClient{}, Cache: &mockCache{}}
	AmadeusService = &amadeusService{HTTPClient: &mockHTTPClient{}}
}

func (m *mockHTTPClient) Get(url string) (*http.Response, error) {

	if url == "http://hotel:8080/hotel" {
		return &http.Response{
				StatusCode: http.StatusNotFound,
			},
			nil
	}

	if strings.Contains(url, "http://hotel:8080/hotel/1") {
		return &http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(strings.NewReader(`{"id":"1"}`)),
		}, nil
	}

	if strings.Contains(url, "http://hotel:8080/hotel/5") {
		return &http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(strings.NewReader(`{"id":"5", "room_amount":1}`)),
		}, nil
	}

	if strings.Contains(url, "http://search:8085/hotel?city=") {

		body := `[{"id":"1"}, {"id":"2"}, {"id":"3"}, {"id":"4"}, {"id":"5"}]`

		return &http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(strings.NewReader(body)),
		}, nil
	}

	return &http.Response{
		StatusCode: http.StatusOK,
		Body:       io.NopCloser(strings.NewReader(`{}`)),
	}, nil
}

func (m *mockHTTPClient) Do(_ *http.Request) (*http.Response, error) {

	body := `{"data":[{"available":true}]}`

	return &http.Response{
		StatusCode: http.StatusOK,
		Body:       io.NopCloser(strings.NewReader(body)),
	}, nil
}

func (c *mockCache) Set(key string, value []byte) {}

func (c *mockCache) Get(key string) ([]byte, error) {
	return nil, errors.New("cache miss")
}

func (t *testReservation) InsertReservation(reservation model.Reservation) model.Reservation {

	if reservation.UserId == 0 {
		return model.Reservation{}
	}

	reservation.Id = 1
	return reservation
}

func (t *testReservation) GetReservationById(id int) model.Reservation {

	if id == 0 {
		return model.Reservation{}
	}

	if id == 2 {
		startTime := time.Now().Add(-72 * time.Hour).Format("02-01-2006 15:04")
		return model.Reservation{
			Id:        id,
			StartDate: startTime,
		}
	}

	if id == 3 {
		startTime := time.Now().Add(72 * time.Hour).Format("02-01-2006 15:04")
		return model.Reservation{
			Id:        id,
			StartDate: startTime,
		}
	}

	return model.Reservation{Id: id}
}

func (t *testReservation) GetReservations() model.Reservations {

	return model.Reservations{}
}

func (t *testReservation) GetReservationsByUser(_ int) model.Reservations {

	return model.Reservations{}
}

func (t *testReservation) GetReservationsByHotel(id string) model.Reservations {

	if id == "5" {
		return model.Reservations{
			model.Reservation{
				Id:        1,
				StartDate: "20-11-2023 15:00",
				EndDate:   "21-11-2023 11:00",
				UserId:    1,
				HotelId:   id,
				Amount:    10900,
			},
		}
	}

	return model.Reservations{}
}

func (t *testReservation) DeleteReservation(reservation model.Reservation) error {

	if reservation.Id == 0 {
		return errors.New("reservation not found")
	}

	return nil
}

func TestInsertReservation_UserError(t *testing.T) {

	a := assert.New(t)

	body := dto.ReservationDto{}

	_, err := ReservationService.InsertReservation(body)

	a.NotNil(err)

	expectedResponse := "user not found"
	a.Equal(expectedResponse, err.Error())
}

func TestInsertReservation_HotelError(t *testing.T) {

	a := assert.New(t)

	body := dto.ReservationDto{
		Id:        1,
		StartDate: "20-11-2023 15:00",
		EndDate:   "21-11-2023 11:00",
		UserId:    10,
		HotelId:   "",
		Amount:    2090807.4,
	}

	_, err := ReservationService.InsertReservation(body)
	a.NotNil(err)

	expectedResponse := "hotel not found"
	a.Equal(expectedResponse, err.Error())
}

func TestInsertReservation_DateError(t *testing.T) {

	a := assert.New(t)

	body := dto.ReservationDto{
		Id:        1,
		StartDate: "22-11-2023 15:00",
		EndDate:   "20-11-2023 11:00",
		UserId:    10,
		HotelId:   "1",
		Amount:    2090807.4,
	}

	_, err := ReservationService.InsertReservation(body)
	a.NotNil(err)

	expectedResponse := "a reservation can't end before it starts"
	a.Equal(expectedResponse, err.Error())
}

func TestInsertReservation(t *testing.T) {

	a := assert.New(t)

	body := dto.ReservationDto{
		StartDate: "20-11-2023 15:00",
		EndDate:   "21-11-2023 11:00",
		UserId:    10,
		HotelId:   "1",
		Amount:    2090807.4,
	}

	response, err := ReservationService.InsertReservation(body)
	a.Nil(err)
	a.Equal(1, response.Id)
}

func TestGetReservationById_NotFound(t *testing.T) {

	a := assert.New(t)

	_, err := ReservationService.GetReservationById(0)

	a.NotNil(err)

	expectedResponse := "reservation not found"
	a.Equal(expectedResponse, err.Error())
}

func TestGetReservationById_Found(t *testing.T) {

	a := assert.New(t)

	id := 1

	response, err := ReservationService.GetReservationById(id)

	a.Nil(err)

	expectedResponse := dto.ReservationDto{Id: id}
	a.Equal(expectedResponse, response)
}

func TestGetReservations(t *testing.T) {

	a := assert.New(t)

	response, err := ReservationService.GetReservations()

	a.Nil(err)

	var expectedResponse dto.ReservationsDto
	a.Equal(expectedResponse, response)
}

func TestGetReservationsByUser_NotFound(t *testing.T) {

	a := assert.New(t)

	userId := 0

	_, err := ReservationService.GetReservationsByUser(userId)
	a.NotNil(err)

	expectedResponse := "user not found"
	a.Equal(expectedResponse, err.Error())
}

func TestGetReservationsByUser_Found(t *testing.T) {

	a := assert.New(t)

	userId := 1

	response, err := ReservationService.GetReservationsByUser(userId)
	a.Nil(err)

	expectedResponse := dto.UserReservationsDto{
		UserId: userId,
	}
	a.Equal(expectedResponse, response)
}

func TestDeleteReservation_NotFound(t *testing.T) {

	a := assert.New(t)

	reservationId := 0
	err := ReservationService.DeleteReservation(reservationId)

	a.NotNil(err)

	expectedResponse := "reservation not found"
	a.Equal(expectedResponse, err.Error())
}

func TestDeleteReservation_Error(t *testing.T) {

	a := assert.New(t)

	reservationId := 2
	err := ReservationService.DeleteReservation(reservationId)

	a.NotNil(err)

	expectedResponse := "can't delete a reservation 48hs before it starts"
	a.Equal(expectedResponse, err.Error())
}

func TestDeleteReservation_Success(t *testing.T) {

	a := assert.New(t)

	reservationId := 3
	err := ReservationService.DeleteReservation(reservationId)

	a.Nil(err)
}

func TestCheckAvailability_NotAvailable(t *testing.T) {

	a := assert.New(t)

	hotelId := "5"
	startDate := time.Date(2023, time.November, 20, 15, 00, 00, 00, time.UTC)
	endDate := time.Date(2023, time.November, 21, 11, 00, 00, 00, time.UTC)

	available := ReservationService.CheckAvailability(hotelId, startDate, endDate)

	a.False(available)

}

func TestCheckAvailability_Available(t *testing.T) {

	a := assert.New(t)

	hotelId := "5"
	startDate := time.Date(2023, time.November, 24, 15, 00, 00, 00, time.UTC)
	endDate := time.Date(2023, time.November, 26, 11, 00, 00, 00, time.UTC)

	available := ReservationService.CheckAvailability(hotelId, startDate, endDate)

	a.True(available)

}

func TestCheckAllAvailability_Error(t *testing.T) {

	a := assert.New(t)

	city := "city"
	startDate := "24-11-2023 15:00"
	endDate := "22-11-2023 11:00"

	_, err := ReservationService.CheckAllAvailability(city, startDate, endDate)

	a.NotNil(err)

	expectedResponse := "a reservation can't end before it starts"
	a.Equal(expectedResponse, err.Error())

}

func TestCheckAllAvailability_Success(t *testing.T) {

	a := assert.New(t)

	city := "city"
	startDate1 := "20-11-2023 15:00"
	endDate1 := "22-11-2023 11:00"

	response1, err := ReservationService.CheckAllAvailability(city, startDate1, endDate1)
	a.Nil(err)
	a.Equal(4, len(response1))

	startDate2 := "28-11-2023 15:00"
	endDate2 := "30-11-2023 11:00"

	response2, err := ReservationService.CheckAllAvailability(city, startDate2, endDate2)
	a.Nil(err)
	a.Equal(5, len(response2))
}
