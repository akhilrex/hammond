package controllers

import (
	"net/http"

	"github.com/akhilrex/hammond/common"
	"github.com/akhilrex/hammond/db"
	"github.com/akhilrex/hammond/models"
	"github.com/akhilrex/hammond/service"
	"github.com/gin-gonic/gin"
)

func RegisterAnonMasterConroller(router *gin.RouterGroup) {
	router.GET("/masters", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"fuelUnits":     db.FuelUnitDetails,
			"fuelTypes":     db.FuelTypeDetails,
			"distanceUnits": db.DistanceUnitDetails,
			"roles":         db.RoleDetails,
			"currencies":    models.GetCurrencyMasterList(),
		})
	})
}
func RegisterMastersController(router *gin.RouterGroup) {

	router.GET("/settings", getSettings)
	router.POST("/settings", udpateSettings)
	router.POST("/me/settings", udpateMySettings)

}

func getSettings(c *gin.Context) {

	c.JSON(http.StatusOK, service.GetSettings())
}
func udpateSettings(c *gin.Context) {
	var model models.UpdateSettingModel
	if err := c.ShouldBind(&model); err == nil {
		err := service.UpdateSettings(model.Currency, *model.DistanceUnit)
		if err != nil {
			c.JSON(http.StatusUnprocessableEntity, common.NewError("udpateSettings", err))
			return
		}
		c.JSON(http.StatusOK, gin.H{})
	} else {
		c.JSON(http.StatusUnprocessableEntity, common.NewValidatorError(err))
	}

}

func udpateMySettings(c *gin.Context) {
	var model models.UpdateSettingModel
	if err := c.ShouldBind(&model); err == nil {
		err := service.UpdateUserSettings(c.MustGet("userId").(string), model.Currency, *model.DistanceUnit)
		if err != nil {
			c.JSON(http.StatusUnprocessableEntity, common.NewError("udpateMySettings", err))
			return
		}
		c.JSON(http.StatusOK, gin.H{})
	} else {
		c.JSON(http.StatusUnprocessableEntity, common.NewValidatorError(err))
	}

}
