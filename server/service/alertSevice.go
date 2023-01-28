package service

import (
	"errors"
	"time"

	"github.com/akhilrex/hammond/db"
	"github.com/akhilrex/hammond/models"
)

func CreateAlert(model models.CreateAlertModel, vehicleId, userId string) (*db.VehicleAlert, error) {
	alert := db.VehicleAlert{
		VehicleID:       vehicleId,
		UserID:          userId,
		Title:           model.Title,
		Comments:        model.Comments,
		StartDate:       model.StartDate,
		StartOdoReading: model.StartOdoReading,
		DistanceUnit:    *model.DistanceUnit,
		AlertFrequency:  *model.AlertFrequency,
		OdoFrequency:    model.OdoFrequency,
		DayFrequency:    model.DayFrequency,
		AlertAllUsers:   model.AlertAllUsers,
		IsActive:        model.IsActive,
		AlertType:       *model.AlertType,
	}
	tx := db.DB.Create(&alert)
	if tx.Error != nil {
		return nil, tx.Error
	}
	go CreateAlertInstance(alert.ID)
	return &alert, nil
}

func CreateAlertInstance(alertId string) error {
	alert, err := db.GeAlertById(alertId)
	if err != nil {
		return err
	}
	existingOccurence, err := db.GetAlertOccurenceByAlertId(alertId)
	if err != nil {
		return err
	}
	var lastOccurance db.AlertOccurance
	useOccurance := false

	if len(*existingOccurence) > 0 {
		lastOccurance = (*existingOccurence)[0]
		useOccurance = true
		if alert.AlertFrequency == db.ONETIME {
			return errors.New("Only single occurance is possible for this kind of alert")
		}
	}
	users := []string{alert.UserID}
	if alert.AlertAllUsers {
		allUsers, err := db.GetVehicleUsers(alert.VehicleID)
		if err != nil {
			return err
		}
		users = make([]string, len(*allUsers))
		for i, user := range *allUsers {
			users[i] = user.UserID
		}
	}

	for _, userId := range users {
		model := db.AlertOccurance{
			VehicleID:      alert.VehicleID,
			UserID:         userId,
			VehicleAlertID: alertId,
		}

		if alert.AlertType == db.DISTANCE || alert.AlertType == db.BOTH {
			model.OdoReading = alert.StartOdoReading + alert.OdoFrequency
			if useOccurance {
				model.OdoReading = lastOccurance.OdoReading + alert.OdoFrequency
			}
		}
		if alert.AlertType == db.TIME || alert.AlertType == db.BOTH {
			date := alert.StartDate.Add(time.Duration(alert.DayFrequency) * 24 * time.Hour)
			if useOccurance {
				date = lastOccurance.Date.Add(time.Duration(alert.DayFrequency) * 24 * time.Hour)
			}
			model.Date = &date
		}
		tx := db.DB.Create(&model)
		if tx.Error != nil {
			return tx.Error
		}
	}
	return nil

}

func ProcessAlertOccurance(occurance db.AlertOccurance, today time.Time) error {
	if occurance.ProcessDate != nil {
		return errors.New("Alert occurence already processed")
	}
	alert := occurance.VehicleAlert
	if !alert.IsActive {
		return errors.New("Alert is not active")
	}
	notification := db.Notification{
		Title:      alert.Title,
		Content:    alert.Comments,
		UserID:     occurance.UserID,
		VehicleID:  occurance.VehicleID,
		Date:       today,
		ParentID:   occurance.ID,
		ParentType: "AlertOccurance",
	}
	var alertProcessType db.AlertType
	if alert.AlertType == db.DISTANCE || alert.AlertType == db.BOTH {
		odoReading, err := GetLatestOdoReadingForVehicle(occurance.VehicleID)
		if err != nil {
			return err
		}
		if odoReading >= occurance.OdoReading {
			alertProcessType = db.DISTANCE
		}
	}
	if alert.AlertType == db.TIME || alert.AlertType == db.BOTH {
		if occurance.Date.Before(today) {
			alertProcessType = db.TIME
		}
	}

	db.DB.Create(&notification)
	return db.MarkAlertOccuranceAsProcessed(occurance.ID, alertProcessType, today)

}

func FindAlertOccurancesToProcess(today time.Time) ([]db.AlertOccurance, error) {
	occurances, err := db.GetUnprocessedAlertOccurances()
	if err != nil {
		return nil, err
	}
	if len(*occurances) == 0 {
		return make([]db.AlertOccurance, 0), nil
	}

	var toReturn []db.AlertOccurance

	for _, occurance := range *occurances {
		alert := occurance.VehicleAlert
		if !alert.IsActive {
			continue
		}
		if alert.AlertType == db.DISTANCE || alert.AlertType == db.BOTH {
			odoReading, err := GetLatestOdoReadingForVehicle(occurance.VehicleID)
			if err != nil {
				return nil, err
			}
			if odoReading >= occurance.OdoReading {
				toReturn = append(toReturn, occurance)
				continue
			}
		}
		if alert.AlertType == db.TIME || alert.AlertType == db.BOTH {
			if occurance.Date.Before(today) {
				toReturn = append(toReturn, occurance)
				continue
			}
		}

	}
	return toReturn, nil
}

func MarkAlertOccuranceAsCompleted() {

}
