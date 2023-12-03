package controller

import (
	"Docker-Containers/client"
	"github.com/gin-gonic/gin"
	"net/http"
)

func GetStats(c *gin.Context) {

	stats, err := client.GetStats()

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, stats)
}

func GetStatsByService(c *gin.Context) {

	service := c.Param("service")

	stats, err := client.GetStatsByService(service)

	if err != nil {

		if err.Error() == "service does not exist" {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, stats)

}
