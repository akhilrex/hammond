package models

import "github.com/akhilrex/hammond/db"

type UpdateSettingModel struct {
	Currency     string           `json:"currency" form:"currency" query:"currency"`
	DistanceUnit *db.DistanceUnit `json:"distanceUnit" form:"distanceUnit" query:"distanceUnit" `
}

type ClarksonMigrationModel struct {
	Url string `json:"url" form:"url" query:"url"`
}
