package models

import (
	"time"

	"github.com/akhilrex/hammond/db"
)

type CreateAlertModel struct {
	Comments        string             `json:"comments"`
	Title           string             `json:"title"`
	StartDate       time.Time          `json:"date"`
	StartOdoReading int                `json:"startOdoReading"`
	DistanceUnit    *db.DistanceUnit   `json:"distanceUnit"`
	AlertFrequency  *db.AlertFrequency `json:"alertFrequency"`
	OdoFrequency    int                `json:"odoFrequency"`
	DayFrequency    int                `json:"dayFrequency"`
	AlertAllUsers   bool               `json:"alertAllUsers"`
	IsActive        bool               `json:"isActive"`
	AlertType       *db.AlertType      `json:"alertType"`
}
