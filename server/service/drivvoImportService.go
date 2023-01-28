package service

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/akhilrex/hammond/db"
)

func DrivvoParseExpenses(content []byte, user *db.User, vehicle *db.Vehicle) ([]db.Expense, []string) {
	expenseReader := csv.NewReader(bytes.NewReader(content))
	expenseReader.Comment = '#'
	// Read headers (there is a trailing comma at the end, that's why we have to read the first line)
	expenseReader.Read()
	expenseReader.FieldsPerRecord = 6
	expenseRecords, err := expenseReader.ReadAll()

	var errors []string
	if err != nil {
		errors = append(errors, err.Error())
		println(err.Error())
		return nil, errors
	}

	var expenses []db.Expense
	for index, record := range expenseRecords {
		date, err := time.Parse("2006-01-02 15:04:05", record[1])
		if err != nil {
			errors = append(errors, "Found an invalid date/time at service/expense row "+strconv.Itoa(index+1))
		}

		totalCost, err := strconv.ParseFloat(record[2], 32)
		if err != nil {
			errors = append(errors, "Found and invalid total cost at service/expense row "+strconv.Itoa(index+1))
		}

		odometer, err := strconv.Atoi(record[0])
		if err != nil {
			errors = append(errors, "Found an invalid odometer reading at service/expense row "+strconv.Itoa(index+1))
		}

		notes := fmt.Sprintf("Location: %s\nNotes: %s\n", record[4], record[5])

		expenses = append(expenses, db.Expense{
			UserID:       user.ID,
			VehicleID:    vehicle.ID,
			Date:         date,
			OdoReading:   odometer,
			Amount:       float32(totalCost),
			ExpenseType:  record[3],
			Currency:     user.Currency,
			DistanceUnit: user.DistanceUnit,
			Comments:     notes,
			Source:       "Drivvo",
		})
	}

	return expenses, errors
}

func DrivvoParseRefuelings(content []byte, user *db.User, vehicle *db.Vehicle, importLocation bool) ([]db.Fillup, []string) {
	refuelingReader := csv.NewReader(bytes.NewReader(content))
	refuelingReader.Comment = '#'
	refuelingRecords, err := refuelingReader.ReadAll()

	var errors []string
	if err != nil {
		errors = append(errors, err.Error())
		println(err.Error())
		return nil, errors
	}

	var fillups []db.Fillup
	for index, record := range refuelingRecords {
		// Skip column titles
		if index == 0 {
			continue
		}

		date, err := time.Parse("2006-01-02 15:04:05", record[1])
		if err != nil {
			errors = append(errors, "Found an invalid date/time at refuel row "+strconv.Itoa(index+1))
		}

		totalCost, err := strconv.ParseFloat(record[4], 32)
		if err != nil {
			errors = append(errors, "Found and invalid total cost at refuel row "+strconv.Itoa(index+1))
		}

		odometer, err := strconv.Atoi(record[0])
		if err != nil {
			errors = append(errors, "Found an invalid odometer reading at refuel row "+strconv.Itoa(index+1))
		}

		location := ""
		if importLocation {
			location = record[17]
		}

		pricePerUnit, err := strconv.ParseFloat(record[3], 32)
		if err != nil {
			unit := strings.ToLower(db.FuelUnitDetails[vehicle.FuelUnit].Key)
			errors = append(errors, fmt.Sprintf("Found an invalid cost per %s at refuel row %d", unit, index+1))
		}

		quantity, err := strconv.ParseFloat(record[5], 32)
		if err != nil {
			errors = append(errors, "Found an invalid quantity at refuel row "+strconv.Itoa(index+1))
		}

		isTankFull := record[6] == "Yes"

		// Unfortunatly, drivvo doesn't expose this info in their export
		fal := false

		notes := fmt.Sprintf("Reason: %s\nNotes: %s\nFuel: %s\n", record[18], record[19], record[2])

		fillups = append(fillups, db.Fillup{
			VehicleID:       vehicle.ID,
			UserID:          user.ID,
			Date:            date,
			HasMissedFillup: &fal,
			IsTankFull:      &isTankFull,
			FuelQuantity:    float32(quantity),
			PerUnitPrice:    float32(pricePerUnit),
			FillingStation:  location,
			OdoReading:      odometer,
			TotalAmount:     float32(totalCost),
			FuelUnit:        vehicle.FuelUnit,
			Currency:        user.Currency,
			DistanceUnit:    user.DistanceUnit,
			Comments:        notes,
			Source:          "Drivvo",
		})

	}
	return fillups, errors
}
