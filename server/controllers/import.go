package controllers

import (
	"net/http"

	"github.com/akhilrex/hammond/service"
	"github.com/gin-gonic/gin"
)

func RegisteImportController(router *gin.RouterGroup) {
	router.POST("/import/fuelly", fuellyImport)
}

func fuellyImport(c *gin.Context) {
	bytes, err := getFileBytes(c, "file")
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, err)
		return
	}
	err = service.FuellyImport(bytes, c.MustGet("userId").(string))
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"errors": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{})
}
