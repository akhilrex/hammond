package controllers

import (
	"net/http"

	"github.com/akhilrex/hammond/common"
	"github.com/akhilrex/hammond/models"
	"github.com/akhilrex/hammond/service"
	"github.com/gin-gonic/gin"
)

func RegisterReportsController(router *gin.RouterGroup) {
	router.GET("/vehicles/:id/mileage", getMileageForVehicle)
}

func getMileageForVehicle(c *gin.Context) {

	var searchByIdQuery models.SearchByIdQuery

	if err := c.ShouldBindUri(&searchByIdQuery); err == nil {
		var model models.MileageQueryModel
		err := c.BindQuery(&model)
		if err != nil {
			c.JSON(http.StatusUnprocessableEntity, common.NewError("getMileageForVehicle", err))
			return
		}

		fillups, err := service.GetMileageByVehicleId(searchByIdQuery.Id, model.Since, model.MileageOption)
		if err != nil {
			c.JSON(http.StatusUnprocessableEntity, common.NewError("getMileageForVehicle", err))
			return
		}
		c.JSON(http.StatusOK, fillups)
	} else {
		c.JSON(http.StatusUnprocessableEntity, common.NewValidatorError(err))
	}
}
