package controllers

import (
	"fmt"
	"net/http"

	"github.com/akhilrex/hammond/common"
	"github.com/akhilrex/hammond/db"
	"github.com/akhilrex/hammond/models"
	"github.com/akhilrex/hammond/service"
	"github.com/gin-gonic/gin"
)

func RegisterSetupController(router *gin.RouterGroup) {
	router.POST("/clarkson/check", canMigrate)
	router.POST("/clarkson/migrate", migrate)
	router.GET("/system/status", appInitialized)
}

func appInitialized(c *gin.Context) {
	canInitialize, err := service.CanInitializeSystem()
	message := ""
	if err != nil {
		message = err.Error()
	}
	c.JSON(http.StatusOK, gin.H{"initialized": !canInitialize, "message": message})
}

func canMigrate(c *gin.Context) {
	var request models.ClarksonMigrationModel
	if err := c.ShouldBind(&request); err == nil {
		canMigrate, data, errr := db.CanMigrate(request.Url)
		errorMessage := ""
		if errr != nil {
			errorMessage = errr.Error()
		}

		c.JSON(http.StatusOK, gin.H{
			"canMigrate": canMigrate,
			"data":       data,
			"message":    errorMessage,
		})
	} else {
		c.JSON(http.StatusUnprocessableEntity, common.NewValidatorError(err))
	}
}

func migrate(c *gin.Context) {
	var request models.ClarksonMigrationModel
	if err := c.ShouldBind(&request); err == nil {
		canMigrate, _, _ := db.CanMigrate(request.Url)

		if !canMigrate {
			c.JSON(http.StatusBadRequest, fmt.Errorf("cannot migrate database. please check connection string."))
			return
		}

		success, err := db.MigrateClarkson(request.Url)
		if !success {
			c.JSON(http.StatusBadRequest, err)
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"success": success,
		})
	} else {
		c.JSON(http.StatusUnprocessableEntity, common.NewValidatorError(err))
	}
}
