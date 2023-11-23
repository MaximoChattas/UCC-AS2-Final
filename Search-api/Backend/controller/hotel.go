package controller

import (
	"Search/dto"
	"Search/service"
	"encoding/json"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"io"
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
	var err error

	city := c.Query("city")

	if city == "" {
		hotelsDto, err = service.HotelService.GetHotels()
	} else {
		hotelsDto, err = service.HotelService.GetHotelByCity(city)
	}

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, hotelsDto)
}

func InsertData() {

	resp, err := http.Get("http://hotelnginx:8080/hotel")

	if err != nil {
		log.Error("Error in HTTP request: ", err)
		return
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)

	if err != nil {
		log.Error("Error reading response: ", err)
		return
	}

	var hotelsDto dto.HotelsDto

	err = json.Unmarshal(body, &hotelsDto)

	if err != nil {
		log.Error("Error parsing JSON: ", err)
		return
	}

	log.Debug(hotelsDto)

	for _, hotel := range hotelsDto {

		err = service.HotelService.InsertUpdateHotel(hotel)
		if err != nil {
			log.Error(err)
			return
		}
	}

}
