package service

import (
	"sort"
	"time"

	"github.com/akhilrex/hammond/db"
	"github.com/akhilrex/hammond/models"
)

func GetMileageByVehicleId(vehicleId string, since time.Time) (mileage []models.MileageModel, err error) {
	data, err := db.GetFillupsByVehicleIdSince(vehicleId, since)
	if err != nil {
		return nil, err
	}

	fillups := make([]db.Fillup, len(*data))
	copy(fillups, *data)
	sort.Slice(fillups, func(i, j int) bool {
		return fillups[i].OdoReading > fillups[j].OdoReading
	})

	var mileages []models.MileageModel

	for i := 0; i < len(fillups)-1; i++ {
		last := i + 1

		currentFillup := fillups[i]
		lastFillup := fillups[last]

		mileage := models.MileageModel{
			Date:         currentFillup.Date,
			VehicleID:    currentFillup.VehicleID,
			FuelUnit:     currentFillup.FuelUnit,
			FuelQuantity: currentFillup.FuelQuantity,
			PerUnitPrice: currentFillup.PerUnitPrice,
			OdoReading:   currentFillup.OdoReading,
			Currency:     currentFillup.Currency,
			Mileage:      0,
			CostPerMile:  0,
		}

		if currentFillup.IsTankFull != nil && *currentFillup.IsTankFull && (currentFillup.HasMissedFillup == nil || !(*currentFillup.HasMissedFillup)) {
			distance := float32(currentFillup.OdoReading - lastFillup.OdoReading)
			mileage.Mileage = distance / currentFillup.FuelQuantity
			mileage.CostPerMile = distance / currentFillup.TotalAmount

		}

		mileages = append(mileages, mileage)
	}
	if mileages == nil {
		mileages = make([]models.MileageModel, 0)
	}
	return mileages, nil
}
