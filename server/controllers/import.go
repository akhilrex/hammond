package controllers

import (
	"net/http"

	"github.com/akhilrex/hammond/service"
	"github.com/gin-gonic/gin"
)

func RegisteImportController(router *gin.RouterGroup) {
	router.POST("/import/fuelly", fuellyImport)
	router.POST("/import/drivvo", drivvoImport)
}

func fuellyImport(c *gin.Context) {
	bytes, err := getFileBytes(c, "file")
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, err)
		return
	}
	errors := service.FuellyImport(bytes, c.MustGet("userId").(string))
	if len(errors) > 0 {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"errors": errors})
		return
	}
	c.JSON(http.StatusOK, gin.H{})
}

func drivvoImport(c *gin.Context) {
	bytes, err := getFileBytes(c, "file")
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, err)
		return
	}
	vehicleId := c.PostForm("vehicleID")
	if vehicleId == "" {
		c.JSON(http.StatusUnprocessableEntity, "Missing Vehicle ID")
		return
	}
	errors := service.DrivvoImport(bytes, c.MustGet("userId").(string), vehicleId)
	if len(errors) > 0 {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"errors": errors})
		return
	}
	c.JSON(http.StatusOK, gin.H{})
}
