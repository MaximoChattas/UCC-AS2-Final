package controller

import (
	"Docker-Containers/client"
	"github.com/gin-gonic/gin"
	"net/http"
)

func GetStats(c *gin.Context) {

	stats, err := client.GetStats()

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
	}

	c.JSON(http.StatusOK, stats)
}
