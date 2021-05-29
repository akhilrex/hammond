package controllers

import (
	"net/http"

	"github.com/akhilrex/hammond/db"
	"github.com/gin-gonic/gin"
)

func RegisterUserController(router *gin.RouterGroup) {
	router.GET("/users", allUsers)
}

func allUsers(c *gin.Context) {
	users, err := db.GetAllUsers()
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}
	c.JSON(http.StatusOK, users)

}
