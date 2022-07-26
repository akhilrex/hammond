package service

import (
	"bytes"

	"github.com/akhilrex/hammond/db"
)

func WriteToDB(fillups []db.Fillup, expenses []db.Expense) []string {
	var errors []string
	tx := db.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	if err := tx.Error; err != nil {
		errors = append(errors, err.Error())
		return errors
	}
	if fillups != nil {
		if err := tx.Create(&fillups).Error; err != nil {
			tx.Rollback()
			errors = append(errors, err.Error())
			return errors
		}
	}
	if expenses != nil {
		if err := tx.Create(&expenses).Error; err != nil {
			tx.Rollback()
			errors = append(errors, err.Error())
			return errors
		}
	}
	err := tx.Commit().Error
	if err != nil {
		errors = append(errors, err.Error())
	}
	return errors

}

func DrivvoImport(content []byte, userId string, vehicleId string, importLocation bool) []string {
	var errors []string
	user, err := GetUserById(userId)
	if err != nil {
		errors = append(errors, err.Error())
		return errors
	}

	vehicle, err := GetVehicleById(vehicleId)
	if err != nil {
		errors = append(errors, err.Error())
		return errors
	}

	endParseIndex := bytes.Index(content, []byte("#Income"))
	if endParseIndex == -1 {
		endParseIndex = bytes.Index(content, []byte("#Route"))
		if endParseIndex == -1 {
			endParseIndex = len(content)
		}

	}

	serviceEndIndex := bytes.Index(content, []byte("#Expense"))
	if serviceEndIndex == -1 {
		serviceEndIndex = endParseIndex
	}

	refuelEndIndex := bytes.Index(content, []byte("#Service"))
	if refuelEndIndex == -1 {
		refuelEndIndex = serviceEndIndex
	}

	var fillups []db.Fillup
	fillups, errors = DrivvoParseRefuelings(content[:refuelEndIndex], user, vehicle, importLocation)

	var allExpenses []db.Expense
	services, parseErrors := DrivvoParseExpenses(content[refuelEndIndex:serviceEndIndex], user, vehicle)
	if parseErrors != nil {
		errors = append(errors, parseErrors...)
	}
	allExpenses = append(allExpenses, services...)

	expenses, parseErrors := DrivvoParseExpenses(content[serviceEndIndex:endParseIndex], user, vehicle)
	if parseErrors != nil {
		errors = append(errors, parseErrors...)
	}

	allExpenses = append(allExpenses, expenses...)

	if len(errors) != 0 {
		return errors
	}

	return WriteToDB(fillups, allExpenses)
}

func FuellyImport(content []byte, userId string) []string {
	fillups, expenses, errors := FuellyParseAll(content, userId)
	if len(errors) != 0 {
		return errors
	}

	return WriteToDB(fillups, expenses)
}
