package service

import (
	"sort"
	"time"

	"github.com/akhilrex/hammond/common"
	"github.com/akhilrex/hammond/db"
	"github.com/akhilrex/hammond/models"
)

func GetMileageByVehicleId(vehicleId string, since time.Time, mileageOption string) (mileage []models.MileageModel, err error) {
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
			DistanceUnit: currentFillup.DistanceUnit,
			Mileage:      0,
			CostPerMile:  0,
		}

		if currentFillup.IsTankFull != nil && *currentFillup.IsTankFull && (currentFillup.HasMissedFillup == nil || !(*currentFillup.HasMissedFillup)) {
			currentOdoReading := float32(currentFillup.OdoReading);
			lastFillupOdoReading := float32(lastFillup.OdoReading);
			currentFuelQuantity := float32(currentFillup.FuelQuantity);
			// If miles per gallon option and distanceUnit is km, convert from km to miles 
			// 	then check if fuel unit is litres. If it is, convert to gallons
			if (mileageOption == "mpg" && mileage.DistanceUnit == db.KILOMETERS) {
				currentOdoReading = common.KmToMiles(currentOdoReading);
				lastFillupOdoReading = common.KmToMiles(lastFillupOdoReading);
				if (mileage.FuelUnit == db.LITRE) {
					currentFuelQuantity = common.LitreToGallon(currentFuelQuantity);
				}
			}

			// If km_litre option or litre per 100km and distanceUnit is miles, convert from miles to km 
			// 	then check if fuel unit is not litres. If it isn't, convert to litres

			if ((mileageOption == "km_litre" || mileageOption == "litre_100km") && mileage.DistanceUnit == db.MILES) {
				currentOdoReading = common.MilesToKm(currentOdoReading);
				lastFillupOdoReading = common.MilesToKm(lastFillupOdoReading);

				if (mileage.FuelUnit == db.US_GALLON) {
					currentFuelQuantity = common.GallonToLitre(currentFuelQuantity);
				}
			} 

			


			distance := float32(currentOdoReading - lastFillupOdoReading);
			if (mileageOption == "litre_100km") {
				mileage.Mileage = currentFuelQuantity / distance * 100;
			} else {
				mileage.Mileage = distance / currentFuelQuantity;
			}

			mileage.CostPerMile = distance / currentFillup.TotalAmount;

		}

		mileages = append(mileages, mileage)
	}
	if mileages == nil {
		mileages = make([]models.MileageModel, 0)
	}
	return mileages, nil
}
