package service

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"strconv"
	"time"

	"github.com/akhilrex/hammond/db"
	"github.com/leekchan/accounting"
)

func FuellyParseAll(content []byte, userId string) ([]db.Fillup, []db.Expense, []string) {
	stream := bytes.NewReader(content)
	reader := csv.NewReader(stream)
	records, err := reader.ReadAll()

	var errors []string
	user, err := GetUserById(userId)
	if err != nil {
		errors = append(errors, err.Error())
		return nil, nil, errors
	}

	vehicles, err := GetUserVehicles(userId)
	if err != nil {
		errors = append(errors, err.Error())
		return nil, nil, errors
	}

	if err != nil {
		errors = append(errors, err.Error())
		return nil, nil, errors
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
	return fillups, expenses, errors
}
