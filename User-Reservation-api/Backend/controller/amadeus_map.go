package controller

import (
	"User-Reservation/dto"
	"User-Reservation/service"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"net/http"
)

func InsertAmadeusMap(c *gin.Context) {
	var amadeusMapDto dto.AmadeusMapDto
	err := c.BindJSON(&amadeusMapDto)

	if err != nil {
		log.Error(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	amadeusMapDto, er := service.AmadeusService.InsertAmadeusMap(amadeusMapDto)

	if er != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": er.Error()})
		return
	}

	c.JSON(http.StatusCreated, amadeusMapDto)
}

func GetAmadeusIdByHotelId(c *gin.Context) {

	hotelId := c.Param("hotel_id")
	var amadeusMapDto dto.AmadeusMapDto

	amadeusMapDto, err := service.AmadeusService.GetAmadeusIdByHotelId(hotelId)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, amadeusMapDto)
}
