package controller

import (
	"Hotel/dto"
	"Hotel/service"
	"fmt"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"io"
	"net/http"
	"os"
	"path"
)

func InsertHotel(c *gin.Context) {
	var hotelDto dto.HotelDto
	err := c.BindJSON(&hotelDto)

	if err != nil {
		log.Error(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	hotelDto, er := service.HotelService.InsertHotel(hotelDto)

	if er != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": er.Error()})
		return
	}

	c.JSON(http.StatusCreated, hotelDto)
}

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

func DeleteHotel(c *gin.Context) {
	id := c.Param("id")

	hotel, err := service.HotelService.DeleteHotel(id)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	for _, image := range hotel.Images {
		err = os.Remove("Images/" + image)

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"message": "Hotel " + hotel.Id + " deleted"})
}

func UpdateHotel(c *gin.Context) {
	id := c.Param("id")
	var hotelDto dto.HotelDto
	err := c.BindJSON(&hotelDto)

	if err != nil {
		log.Error(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	hotelDto.Id = id

	hotelDto, err = service.HotelService.UpdateHotel(hotelDto)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, hotelDto)
}

func InsertImages(c *gin.Context) {
	var hotelDto dto.HotelDto

	id := c.Param("id")

	hotelDto, err := service.HotelService.GetHotelById(id)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	form, _ := c.MultipartForm()
	files := form.File["images"]

	imageCount := len(hotelDto.Images)

	hotelDto.Images = []string{}

	for i, file := range files {

		fileExt := path.Ext(file.Filename)

		//Filename as [hotel_id]-[image_number].[file_extension]
		fileName := fmt.Sprintf("%s-%d%s", id, i+1+imageCount, fileExt)

		err = c.SaveUploadedFile(file, "Images/"+fileName)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save file"})
			return
		}

		hotelDto.Images = append(hotelDto.Images, fileName)
	}

	hotelDto, err = service.HotelService.UpdateHotel(hotelDto)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, hotelDto)
}

func GetImageByName(c *gin.Context) {

	name := c.Query("name")

	filePath := "Images/" + name

	file, err := os.Open(filePath)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	defer file.Close()

	c.Header("Content-Type", "image/jpg")

	_, err = io.Copy(c.Writer, file)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
}
