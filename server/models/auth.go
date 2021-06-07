package models

import "github.com/akhilrex/hammond/db"

type LoginResponse struct {
	Name         string `json:"name"`
	Email        string `json:"email"`
	Token        string `json:"token"`
	RefreshToken string `json:"refreshToken"`
	Role         string `json:"role"`
}

type LoginRequest struct {
	Email    string `form:"email" json:"email" binding:"required,email"`
	Password string `form:"password" json:"password" binding:"required,max=255"`
}

type RegisterRequest struct {
	Name         string           `form:"name" json:"name"`
	Email        string           `form:"email" json:"email" binding:"required,email"`
	Password     string           `form:"password" json:"password" binding:"required,min=8,max=255"`
	Currency     string           `json:"currency" form:"currency" query:"currency"`
	DistanceUnit *db.DistanceUnit `json:"distanceUnit" form:"distanceUnit" query:"distanceUnit" `
	Role         *db.Role         `json:"role" form:"role" query:"role" `
}

type ChangePasswordRequest struct {
	OldPassword string `form:"oldPassword" json:"oldPassword" binding:"required,max=255"`
	NewPassword string `form:"newPassword" json:"newPassword" binding:"required,min=8,max=255"`
}
