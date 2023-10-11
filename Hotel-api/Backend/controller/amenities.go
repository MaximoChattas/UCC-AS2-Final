package controller

import (
	"Hotel/dto"
	"Hotel/service"
	"net/http"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

func InsertAmenity(c *gin.Context) {
	var amenityDto dto.AmenityDto
	err := c.BindJSON(&amenityDto)

	if err != nil {
		log.Error(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	amenityDto, er := service.AmenityService.InsertAmenity(amenityDto)

	if er != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": er.Error()})
		return
	}

	c.JSON(http.StatusCreated, amenityDto)
}

func GetAmenities(c *gin.Context) {

	var amenitiesDto dto.AmenitiesDto

	amenitiesDto, err := service.AmenityService.GetAmenities()

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, amenitiesDto)
}

func DeleteAmenityById(c *gin.Context) {
	id := c.Param("id")

	err := service.AmenityService.DeleteAmenityById(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "amenity deleted successfully"})
}
