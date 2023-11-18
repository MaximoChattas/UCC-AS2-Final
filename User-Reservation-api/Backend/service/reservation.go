package service

import (
	"User-Reservation/cache"
	"User-Reservation/client"
	"User-Reservation/dto"
	"User-Reservation/model"
	"encoding/json"
	"errors"
	"fmt"
	log "github.com/sirupsen/logrus"
	"io"
	"math"
	"net/http"
	"strings"
	"sync"
	"time"
)

type reservationService struct {
	HTTPClient httpClientInterface
}

type reservationServiceInterface interface {
	InsertReservation(reservationDto dto.ReservationDto) (dto.ReservationDto, error)
	GetReservationById(id int) (dto.ReservationDto, error)
	GetReservations() (dto.ReservationsDto, error)
	GetReservationsByUser(userId int) (dto.UserReservationsDto, error)
	DeleteReservation(id int) error
	CheckAvailability(hotelId string, startDate time.Time, endDate time.Time) bool
	CheckAllAvailability(city string, startDate string, endDate string) (dto.HotelsDto, error)
	GetAllHotelsByCity(city string) dto.HotelsDto
	GetHotelInfo(hotelId string) (dto.HotelDto, error)
}

type httpClient struct{}

type httpClientInterface interface {
	Get(url string) (*http.Response, error)
}

var ReservationService reservationServiceInterface

func init() {
	ReservationService = &reservationService{
		HTTPClient: &httpClient{},
	}
}

func (h *httpClient) Get(url string) (*http.Response, error) {
	return http.Get(url)
}

func (s *reservationService) InsertReservation(reservationDto dto.ReservationDto) (dto.ReservationDto, error) {

	userDto := client.UserClient.GetUserById(reservationDto.UserId)
	hotelDto, err := s.GetHotelInfo(reservationDto.HotelId)

	if err != nil {
		return dto.ReservationDto{}, errors.New("error retrieving hotel information")
	}

	if userDto.Id == 0 {
		return reservationDto, errors.New("user not found")
	}

	if hotelDto.Id == "000000000000000000000000" {
		return reservationDto, errors.New("hotel not found")
	}

	timeStart, _ := time.Parse("02-01-2006 15:04", reservationDto.StartDate)
	timeEnd, _ := time.Parse("02-01-2006 15:04", reservationDto.EndDate)

	if timeStart.After(timeEnd) {
		return reservationDto, errors.New("a reservation can't end before it starts")
	}

	amadeusDto, err := AmadeusService.GetAmadeusIdByHotelId(hotelDto.Id)
	if err != nil {
		return reservationDto, err
	}

	available, err := AmadeusService.GetAmadeusAvailability(amadeusDto.AmadeusId, timeStart, timeEnd)
	if err != nil {
		return reservationDto, err
	}

	if !available {
		return reservationDto, errors.New("not available through amadeus")
	}

	if s.CheckAvailability(reservationDto.HotelId, timeStart, timeEnd) {
		var reservation model.Reservation

		reservation.StartDate = reservationDto.StartDate
		reservation.EndDate = reservationDto.EndDate
		reservation.HotelId = reservationDto.HotelId
		reservation.UserId = reservationDto.UserId

		hoursAmount := timeEnd.Sub(timeStart).Hours()
		nightsAmount := math.Ceil(hoursAmount / 24)
		rate := hotelDto.Rate

		reservation.Amount = rate * nightsAmount

		reservation = client.ReservationClient.InsertReservation(reservation)

		reservationDto.Id = reservation.Id
		reservationDto.Amount = reservation.Amount

		return reservationDto, nil
	}

	return reservationDto, errors.New("there are no rooms available")
}

func (s *reservationService) GetReservationById(id int) (dto.ReservationDto, error) {
	var reservation model.Reservation
	var reservationDto dto.ReservationDto

	reservation = client.ReservationClient.GetReservationById(id)

	if reservation.Id == 0 {
		return reservationDto, errors.New("reservation not found")
	}

	reservationDto.Id = reservation.Id
	reservationDto.StartDate = reservation.StartDate
	reservationDto.EndDate = reservation.EndDate
	reservationDto.HotelId = reservation.HotelId
	reservationDto.UserId = reservation.UserId
	reservationDto.Amount = reservation.Amount

	return reservationDto, nil
}

func (s *reservationService) GetReservations() (dto.ReservationsDto, error) {

	var reservations model.Reservations = client.ReservationClient.GetReservations()
	var reservationsDto dto.ReservationsDto

	for _, reservation := range reservations {
		var reservationDto dto.ReservationDto

		reservationDto.Id = reservation.Id
		reservationDto.StartDate = reservation.StartDate
		reservationDto.EndDate = reservation.EndDate
		reservationDto.HotelId = reservation.HotelId
		reservationDto.UserId = reservation.UserId
		reservationDto.Amount = reservation.Amount

		reservationsDto = append(reservationsDto, reservationDto)
	}

	return reservationsDto, nil
}

func (s *reservationService) GetReservationsByUser(userId int) (dto.UserReservationsDto, error) {
	var user model.User = client.UserClient.GetUserById(userId)
	var userReservationsDto dto.UserReservationsDto
	var reservationsDto dto.ReservationsDto

	if user.Id == 0 {
		return userReservationsDto, errors.New("user not found")
	}
	var reservations model.Reservations = client.ReservationClient.GetReservationsByUser(userId)

	userReservationsDto.UserId = user.Id
	userReservationsDto.UserName = user.Name
	userReservationsDto.UserLastName = user.LastName
	userReservationsDto.UserDni = user.Dni
	userReservationsDto.UserEmail = user.Email
	userReservationsDto.UserPassword = user.Password

	for _, reservation := range reservations {
		var reservationDto dto.ReservationDto

		reservationDto.Id = reservation.Id
		reservationDto.StartDate = reservation.StartDate
		reservationDto.EndDate = reservation.EndDate
		reservationDto.HotelId = reservation.HotelId
		reservationDto.UserId = reservation.UserId
		reservationDto.Amount = reservation.Amount

		reservationsDto = append(reservationsDto, reservationDto)
	}

	userReservationsDto.Reservations = reservationsDto

	return userReservationsDto, nil
}

func (s *reservationService) DeleteReservation(id int) error {

	reservation := client.ReservationClient.GetReservationById(id)

	if reservation.Id == 0 {
		return errors.New("reservation not found")
	}

	reservationStart, _ := time.Parse("02-01-2006 15:04", reservation.StartDate)

	if reservationStart.Before(time.Now().Add(-48 * time.Hour)) {
		return errors.New("can't delete a reservation 48hs before it starts")
	}

	err := client.ReservationClient.DeleteReservation(reservation)

	return err

}

func (s *reservationService) GetHotelInfo(hotelId string) (dto.HotelDto, error) {
	resp, err := http.Get("http://hotel:8080/hotel/" + hotelId)

	if err != nil {
		log.Error("Error in HTTP request ", err)
		return dto.HotelDto{}, err
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)

	if err != nil {
		log.Error("Error reading response ", err)
		return dto.HotelDto{}, err
	}

	var hotelDto dto.HotelDto

	err = json.Unmarshal(body, &hotelDto)

	if err != nil {
		log.Error("Error parsing JSON ", err)
		return dto.HotelDto{}, err
	}

	return hotelDto, nil
}

func (s *reservationService) GetAllHotelsByCity(city string) dto.HotelsDto {

	cityFormatted := strings.ReplaceAll(city, " ", "+")
	resp, err := s.HTTPClient.Get("http://search:8085/hotel?city=" + cityFormatted)

	if err != nil {
		log.Error("Error in HTTP request ", err)
		return dto.HotelsDto{}
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)

	if err != nil {
		log.Error("Error reading response ", err)
		return dto.HotelsDto{}
	}

	var hotelsDto dto.HotelsDto

	err = json.Unmarshal(body, &hotelsDto)

	if err != nil {
		log.Error("Error parsing JSON ", err)
		return dto.HotelsDto{}
	}

	return hotelsDto
}

func (s *reservationService) CheckAvailability(hotelId string, startDate time.Time, endDate time.Time) bool {

	hotel, _ := s.GetHotelInfo(hotelId)
	reservations := client.ReservationClient.GetReservationsByHotel(hotelId)

	roomsAvailable := hotel.RoomAmount

	for _, reservation := range reservations {

		reservationStart, _ := time.Parse("02-01-2006 15:04", reservation.StartDate)
		reservationEnd, _ := time.Parse("02-01-2006 15:04", reservation.EndDate)

		if reservationStart.After(startDate) && reservationEnd.Before(endDate) ||
			reservationStart.Before(startDate) && reservationEnd.After(startDate) ||
			reservationStart.Before(endDate) && reservationEnd.After(endDate) ||
			reservationStart.Before(startDate) && reservationEnd.After(endDate) ||
			reservationStart.Equal(startDate) || reservationEnd.Equal(endDate) {
			roomsAvailable--
		}
		if roomsAvailable == 0 {
			return false
		}
	}

	return true
}

func (s *reservationService) CheckAllAvailability(city string, startDate string, endDate string) (dto.HotelsDto, error) {

	var wg sync.WaitGroup
	var hotelsAvailable dto.HotelsDto

	reservationStart, _ := time.Parse("02-01-2006 15:04", startDate)
	reservationEnd, _ := time.Parse("02-01-2006 15:04", endDate)

	if reservationStart.After(reservationEnd) {
		return hotelsAvailable, errors.New("a reservation can't end before it starts")
	}

	cityFormatted := city[0:3]
	cacheKey := fmt.Sprintf("%s/%s/%s", cityFormatted, reservationStart.Format("02-01-06"), reservationEnd.Format("02-01-06"))

	result, err := cache.Get(cacheKey)

	if err == nil {

		err = json.Unmarshal(result, &hotelsAvailable)

		if err != nil {
			return dto.HotelsDto{}, errors.New("error unmarshaling json")
		}

		return hotelsAvailable, nil

	}

	hotels := s.GetAllHotelsByCity(city)

	resultsCh := make(chan dto.HotelDto)

	for _, hotel := range hotels {

		wg.Add(1)
		go func(hotel dto.HotelDto) {

			defer wg.Done()

			if s.CheckAvailability(hotel.Id, reservationStart, reservationEnd) {
				var hotelDto dto.HotelDto
				hotelDto.Id = hotel.Id
				hotelDto.Name = hotel.Name
				hotelDto.City = hotel.City
				hotelDto.StreetName = hotel.StreetName
				hotelDto.StreetNumber = hotel.StreetNumber
				hotelDto.RoomAmount = hotel.RoomAmount
				hotelDto.Rate = hotel.Rate
				hotelDto.Description = hotel.Description

				if len(hotel.Images) > 0 {
					hotelDto.Images = append(hotelDto.Images, hotel.Images[0])
				}
				resultsCh <- hotelDto
			}
		}(hotel)
	}

	go func() {
		wg.Wait()
		close(resultsCh)
	}()

	for hotelDto := range resultsCh {
		hotelsAvailable = append(hotelsAvailable, hotelDto)
	}

	jsonResult, _ := json.Marshal(hotelsAvailable)
	cache.Set(cacheKey, jsonResult)

	return hotelsAvailable, nil
}
