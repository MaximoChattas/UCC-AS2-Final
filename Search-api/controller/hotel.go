package controller

import (
	"Search/dto"
	"Search/service"
	"github.com/gin-gonic/gin"
	"net/http"
)

func GetHotelById(c *gin.Context) {

	id := c.Param("id")
	var hotelDto dto.HotelDto

	hotelDto, err := service.HotelService.GetHotelById(id)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, hotelDto)
}

func GetHotels(c *gin.Context) {

	var hotelsDto dto.HotelsDto

	hotelsDto, err := service.HotelService.GetHotels()

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, hotelsDto)
}
