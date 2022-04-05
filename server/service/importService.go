package service

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/akhilrex/hammond/db"
	"github.com/leekchan/accounting"
)

// TODO: Move drivvo stuff to separate file

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

func DrivvoParseRefuelings(content []byte, user *db.User, vehicle *db.Vehicle) ([]db.Fillup, []string) {
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

		// TODO: Make optional
		location := record[17]

		pricePerUnit, err := strconv.ParseFloat(record[3], 32)
		if err != nil {
			unit := strings.ToLower(db.FuelUnitDetails[vehicle.FuelUnit].Long)
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

func DrivvoImport(content []byte, userId string, vehicleId string) []string {
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

	serviceSectionIndex := bytes.Index(content, []byte("#Service"))

	endParseIndex := bytes.Index(content, []byte("#Income"))
	if endParseIndex == -1 {
		endParseIndex = bytes.Index(content, []byte("#Route"))
		if endParseIndex == -1 {
			endParseIndex = len(content)
		}

	}

	expenseSectionIndex := bytes.Index(content, []byte("#Expense"))
	if expenseSectionIndex == -1 {
		expenseSectionIndex = endParseIndex
	}

	fillups, errors := DrivvoParseRefuelings(content[:serviceSectionIndex], user, vehicle)
	_ = fillups

	var allExpenses []db.Expense
	if serviceSectionIndex != -1 {
		services, parseErrors := DrivvoParseExpenses(content[serviceSectionIndex:expenseSectionIndex], user, vehicle)
		if parseErrors != nil {
			errors = append(errors, parseErrors...)
		}
		allExpenses = append(allExpenses, services...)
	}

	if expenseSectionIndex != endParseIndex {
		expenses, parseErrors := DrivvoParseExpenses(content[expenseSectionIndex:endParseIndex], user, vehicle)
		if parseErrors != nil {
			errors = append(errors, parseErrors...)
		}
		allExpenses = append(allExpenses, expenses...)
	}

	if len(errors) != 0 {
		return errors
	}

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
	if err := tx.Create(&fillups).Error; err != nil {
		tx.Rollback()
		errors = append(errors, err.Error())
		return errors
	}
	if err := tx.Create(&allExpenses).Error; err != nil {
		tx.Rollback()
		errors = append(errors, err.Error())
		return errors
	}
	if err := tx.Commit().Error; err != nil {
		errors = append(errors, err.Error())
	}
	return errors
}

func FuellyImport(content []byte, userId string) []string {
	stream := bytes.NewReader(content)
	reader := csv.NewReader(stream)
	records, err := reader.ReadAll()

	var errors []string
	if err != nil {
		errors = append(errors, err.Error())
		return errors
	}

	vehicles, err := GetUserVehicles(userId)
	if err != nil {
		errors = append(errors, err.Error())
		return errors
	}
	user, err := GetUserById(userId)

	if err != nil {
		errors = append(errors, err.Error())
		return errors
	}

	var vehicleMap map[string]db.Vehicle = make(map[string]db.Vehicle)
	for _, vehicle := range *vehicles {
		vehicleMap[vehicle.Nickname] = vehicle
	}

	var fillups []db.Fillup
	var expenses []db.Expense
	layout := "2006-01-02 15:04"
	altLayout := "2006-01-02 3:04 PM"

	for index, record := range records {
		if index == 0 {
			continue
		}

		var vehicle db.Vehicle
		var ok bool
		if vehicle, ok = vehicleMap[record[4]]; !ok {
			errors = append(errors, "Found an unmapped vehicle entry at row "+strconv.Itoa(index+1))
		}
		dateStr := record[2] + " " + record[3]
		date, err := time.Parse(layout, dateStr)
		if err != nil {
			date, err = time.Parse(altLayout, dateStr)
		}
		if err != nil {
			errors = append(errors, "Found an invalid date/time at row "+strconv.Itoa(index+1))
		}

		totalCostStr := accounting.UnformatNumber(record[9], 3, user.Currency)
		totalCost64, err := strconv.ParseFloat(totalCostStr, 32)
		if err != nil {
			errors = append(errors, "Found an invalid total cost at row "+strconv.Itoa(index+1))
		}

		totalCost := float32(totalCost64)
		odoStr := accounting.UnformatNumber(record[5], 0, user.Currency)
		odoreading, err := strconv.Atoi(odoStr)
		if err != nil {
			errors = append(errors, "Found an invalid odo reading at row "+strconv.Itoa(index+1))
		}
		location := record[12]

		//Create Fillup
		if record[0] == "Gas" {
			rateStr := accounting.UnformatNumber(record[7], 3, user.Currency)
			ratet64, err := strconv.ParseFloat(rateStr, 32)
			if err != nil {
				errors = append(errors, "Found an invalid cost per gallon at row "+strconv.Itoa(index+1))
			}
			rate := float32(ratet64)

			quantity64, err := strconv.ParseFloat(record[8], 32)
			if err != nil {
				errors = append(errors, "Found an invalid quantity at row "+strconv.Itoa(index+1))
			}
			quantity := float32(quantity64)

			notes := fmt.Sprintf("Octane:%s\nGas Brand:%s\nLocation%s\nTags:%s\nPayment Type:%s\nTire Pressure:%s\nNotes:%s\nMPG:%s",
				record[10], record[11], record[12], record[13], record[14], record[15], record[16], record[1],
			)

			isTankFull := record[6] == "Full"
			fal := false
			fillups = append(fillups, db.Fillup{
				VehicleID:       vehicle.ID,
				FuelUnit:        vehicle.FuelUnit,
				FuelQuantity:    quantity,
				PerUnitPrice:    rate,
				TotalAmount:     totalCost,
				OdoReading:      odoreading,
				IsTankFull:      &isTankFull,
				Comments:        notes,
				FillingStation:  location,
				HasMissedFillup: &fal,
				UserID:          userId,
				Date:            date,
				Currency:        user.Currency,
				DistanceUnit:    user.DistanceUnit,
				Source:          "Fuelly",
			})

		}
		if record[0] == "Service" {
			notes := fmt.Sprintf("Tags:%s\nPayment Type:%s\nNotes:%s",
				record[13], record[14], record[16],
			)
			expenses = append(expenses, db.Expense{
				VehicleID:    vehicle.ID,
				Amount:       totalCost,
				OdoReading:   odoreading,
				Comments:     notes,
				ExpenseType:  record[17],
				UserID:       userId,
				Currency:     user.Currency,
				Date:         date,
				DistanceUnit: user.DistanceUnit,
				Source:       "Fuelly",
			})
		}

	}
	if len(errors) != 0 {
		return errors
	}

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
	if err := tx.Create(&fillups).Error; err != nil {
		tx.Rollback()
		errors = append(errors, err.Error())
		return errors
	}
	if err := tx.Create(&expenses).Error; err != nil {
		tx.Rollback()
		errors = append(errors, err.Error())
		return errors
	}
	err = tx.Commit().Error
	if err != nil {
		errors = append(errors, err.Error())
	}
	return errors
}
