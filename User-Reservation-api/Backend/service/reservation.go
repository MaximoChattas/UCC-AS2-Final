package service

import (
	"User-Reservation/client"
	"User-Reservation/dto"
	"User-Reservation/model"
	"encoding/json"
	"errors"
	log "github.com/sirupsen/logrus"
	"io"
	"math"
	"net/http"
	"sync"
	"time"
)

type reservationService struct{}

type reservationServiceInterface interface {
	InsertReservation(reservationDto dto.ReservationDto) (dto.ReservationDto, error)
	GetReservationById(id int) (dto.ReservationDto, error)
	GetReservations() (dto.ReservationsDto, error)
	DeleteReservation(id int) error
	getHotelInfo(hotelId string) (dto.HotelDto, error)
	CheckAvailability(hotelId string, startDate time.Time, endDate time.Time) bool
	CheckAllAvailability(startDate string, endDate string) (dto.HotelsDto, error)
}

var ReservationService reservationServiceInterface

func init() {
	ReservationService = &reservationService{}
}

func (s *reservationService) InsertReservation(reservationDto dto.ReservationDto) (dto.ReservationDto, error) {

	userDto := client.GetUserById(reservationDto.UserId)
	hotelDto, err := s.getHotelInfo(reservationDto.HotelId)

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

	if timeStart.Before(time.Now()) {
		return reservationDto, errors.New("a reservation can't start before current time")
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

		reservation = client.InsertReservation(reservation)

		reservationDto.Id = reservation.Id
		reservationDto.Amount = reservation.Amount

		return reservationDto, nil
	}

	return reservationDto, errors.New("there are no rooms available")
}

func (s *reservationService) GetReservationById(id int) (dto.ReservationDto, error) {
	var reservation model.Reservation
	var reservationDto dto.ReservationDto

	reservation = client.GetReservationById(id)

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

	var reservations model.Reservations = client.GetReservations()
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

func (s *reservationService) DeleteReservation(id int) error {

	reservation := client.GetReservationById(id)

	if reservation.Id == 0 {
		return errors.New("reservation not found")
	}

	reservationStart, _ := time.Parse("02-01-2006 15:04", reservation.StartDate)

	if reservationStart.Before(time.Now().Add(-48 * time.Hour)) {
		return errors.New("can't delete a reservation 48hs before it starts")
	}

	err := client.DeleteReservation(reservation)

	return err

}

func (s *reservationService) getHotelInfo(hotelId string) (dto.HotelDto, error) {
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

func (s *reservationService) getAllHotels() dto.HotelsDto {
	resp, err := http.Get("http://hotel:8080/hotel/")

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

	hotel, _ := s.getHotelInfo(hotelId)
	reservations := client.GetReservationsByHotel(hotelId)

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

func (s *reservationService) CheckAllAvailability(startDate string, endDate string) (dto.HotelsDto, error) {

	var wg sync.WaitGroup
	var hotelsAvailable dto.HotelsDto

	reservationStart, _ := time.Parse("02-01-2006 15:04", startDate)
	reservationEnd, _ := time.Parse("02-01-2006 15:04", endDate)

	if reservationStart.After(reservationEnd) {
		return hotelsAvailable, errors.New("a reservation can't end before it starts")
	}

	hotels := s.getAllHotels()

	resultsCh := make(chan dto.HotelDto)

	for _, hotel := range hotels {

		wg.Add(1)
		go func(hotel dto.HotelDto) {

			defer wg.Done()

			if s.CheckAvailability(hotel.Id, reservationStart, reservationEnd) {
				var hotelDto dto.HotelDto
				hotelDto.Id = hotel.Id
				hotelDto.Name = hotel.Name
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

	return hotelsAvailable, nil
}
