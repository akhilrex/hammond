package controllers

import (
	"net/http"

	"github.com/akhilrex/hammond/common"
	"github.com/akhilrex/hammond/db"
	"github.com/akhilrex/hammond/models"
	"github.com/akhilrex/hammond/service"
	"github.com/gin-gonic/gin"
)

func RegisterUserController(router *gin.RouterGroup) {
	router.GET("/users", allUsers)
	router.POST("/users/:id/enable", ShouldBeAdmin(), enableUser)
	router.POST("/users/:id/disable", ShouldBeAdmin(), disableUser)
}

func allUsers(c *gin.Context) {
	users, err := db.GetAllUsers()
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}
	c.JSON(http.StatusOK, users)

}
func enableUser(c *gin.Context) {
	var searchByIdQuery models.SearchByIdQuery
	if err := c.ShouldBindUri(&searchByIdQuery); err == nil {
		err := service.SetDisabledStatusForUser(searchByIdQuery.Id, false)
		if err != nil {
			c.JSON(http.StatusBadRequest, err)
			return
		}
		c.JSON(http.StatusOK, gin.H{})
	} else {
		c.JSON(http.StatusUnprocessableEntity, common.NewValidatorError(err))
	}

}
func disableUser(c *gin.Context) {
	var searchByIdQuery models.SearchByIdQuery
	if err := c.ShouldBindUri(&searchByIdQuery); err == nil {
		err := service.SetDisabledStatusForUser(searchByIdQuery.Id, true)
		if err != nil {
			c.JSON(http.StatusBadRequest, err)
			return
		}
		c.JSON(http.StatusOK, gin.H{})
	} else {
		c.JSON(http.StatusUnprocessableEntity, common.NewValidatorError(err))
	}

}
