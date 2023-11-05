package controller

import (
	"User-Reservation/dto"
	"User-Reservation/service"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"net/http"
	"strconv"
)

func InsertReservation(c *gin.Context) {
	var reservationDto dto.ReservationDto
	err := c.BindJSON(&reservationDto)

	if err != nil {
		log.Error(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	reservationDto, er := service.ReservationService.InsertReservation(reservationDto)

	if er != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": er.Error()})
		return
	}

	c.JSON(http.StatusCreated, reservationDto)
}

func GetReservationById(c *gin.Context) {

	id, _ := strconv.Atoi(c.Param("id"))
	var reservationDto dto.ReservationDto

	reservationDto, err := service.ReservationService.GetReservationById(id)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, reservationDto)
}

func GetReservations(c *gin.Context) {

	var reservationsDto dto.ReservationsDto

	reservationsDto, err := service.ReservationService.GetReservations()

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, reservationsDto)
}

func DeleteReservation(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	err := service.ReservationService.DeleteReservation(id)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Reservation deleted"})
}

func CheckAvailability(c *gin.Context) {

	city := c.Query("city")
	startDate := c.Query("start_date")
	endDate := c.Query("end_date")

	hotelsAvailable, err := service.ReservationService.CheckAllAvailability(city, startDate, endDate)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	c.JSON(http.StatusOK, hotelsAvailable)
}
