package service

import (
	"github.com/akhilrex/hammond/db"
	"github.com/akhilrex/hammond/models"
	"gorm.io/gorm/clause"
)

func CreateVehicle(model models.CreateVehicleRequest, userId string) (*db.Vehicle, error) {
	vehicle := db.Vehicle{
		Nickname:          model.Nickname,
		Registration:      model.Registration,
		Model:             model.Model,
		Make:              model.Make,
		YearOfManufacture: model.YearOfManufacture,
		EngineSize:        model.EngineSize,
		FuelUnit:          *model.FuelUnit,
		FuelType:          *model.FuelType,
	}

	tx := db.DB.Create(&vehicle)
	if tx.Error != nil {
		return nil, tx.Error
	}
	association := db.UserVehicle{
		UserID:    userId,
		VehicleID: vehicle.ID,
		IsOwner:   true,
	}
	tx = db.DB.Create(&association)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return &vehicle, nil

}

func GetVehicleOwner(vehicleId string) (string, error) {
	return db.GetVehicleOwner(vehicleId)
}

func GetVehicleUsers(vehicleId string) (*[]db.UserVehicle, error) {
	return db.GetVehicleUsers(vehicleId)
}
func ShareVehicle(vehicleId, userId string) error {
	return db.ShareVehicle(vehicleId, userId)
}
func UnshareVehicle(vehicleId, userId string) error {
	return db.UnshareVehicle(vehicleId, userId)
}

func GetVehicleById(vehicleID string) (*db.Vehicle, error) {
	return db.GetVehicleById(vehicleID)
}
func GetFillupsByVehicleId(vehicleId string) (*[]db.Fillup, error) {
	return db.GetFillupsByVehicleId(vehicleId)
}
func GetExpensesByVehicleId(vehicleId string) (*[]db.Expense, error) {
	return db.GetExpensesByVehicleId(vehicleId)
}
func GetFillupById(fillupId string) (*db.Fillup, error) {
	return db.GetFillupById(fillupId)
}
func GetExpenseById(expenseId string) (*db.Expense, error) {
	return db.GetExpenseById(expenseId)
}
func UpdateVehicle(vehicleID string, model models.UpdateVehicleRequest) error {
	toUpdate, err := GetVehicleById(vehicleID)
	if err != nil {
		return err
	}
	//return db.DB.Model(&toUpdate).Updates(db.Vehicle{
	toUpdate.Nickname = model.Nickname
	toUpdate.Registration = model.Registration
	toUpdate.Model = model.Model
	toUpdate.Make = model.Make
	toUpdate.YearOfManufacture = model.YearOfManufacture
	toUpdate.EngineSize = model.EngineSize
	toUpdate.FuelUnit = *model.FuelUnit
	toUpdate.FuelType = *model.FuelType
	//}).Error

	return db.DB.Omit(clause.Associations).Save(toUpdate).Error
}
func GetAllVehicles() (*[]db.Vehicle, error) {
	return db.GetAllVehicles("")
}

func GetUserVehicles(id string) (*[]db.Vehicle, error) {
	return db.GetUserVehicles(id)
}

func CreateFillup(model models.CreateFillupRequest) (*db.Fillup, error) {

	user, err := db.GetUserById(model.UserID)
	if err != nil {
		return nil, err
	}

	fillup := db.Fillup{
		VehicleID:       model.VehicleID,
		FuelUnit:        *model.FuelUnit,
		FuelQuantity:    model.FuelQuantity,
		PerUnitPrice:    model.PerUnitPrice,
		TotalAmount:     model.TotalAmount,
		OdoReading:      model.OdoReading,
		IsTankFull:      model.IsTankFull,
		HasMissedFillup: model.HasMissedFillup,
		Comments:        model.Comments,
		FillingStation:  model.FillingStation,
		UserID:          model.UserID,
		Date:            model.Date,
		Currency:        user.Currency,
		DistanceUnit:    user.DistanceUnit,
	}

	tx := db.DB.Create(&fillup)
	if tx.Error != nil {
		return nil, tx.Error
	}

	return &fillup, nil

}

func CreateExpense(model models.CreateExpenseRequest) (*db.Expense, error) {
	user, err := db.GetUserById(model.UserID)
	if err != nil {
		return nil, err
	}

	expense := db.Expense{
		VehicleID:    model.VehicleID,
		Amount:       model.Amount,
		OdoReading:   model.OdoReading,
		ExpenseType:  model.ExpenseType,
		Comments:     model.Comments,
		UserID:       model.UserID,
		Date:         model.Date,
		Currency:     user.Currency,
		DistanceUnit: user.DistanceUnit,
	}

	tx := db.DB.Create(&expense)
	if tx.Error != nil {
		return nil, tx.Error
	}

	return &expense, nil

}

func UpdateFillup(fillupId string, model models.UpdateFillupRequest) error {
	toUpdate, err := GetFillupById(fillupId)
	if err != nil {
		return err
	}
	return db.DB.Model(&toUpdate).Updates(db.Fillup{
		VehicleID:       model.VehicleID,
		FuelUnit:        *model.FuelUnit,
		FuelQuantity:    model.FuelQuantity,
		PerUnitPrice:    model.PerUnitPrice,
		TotalAmount:     model.TotalAmount,
		OdoReading:      model.OdoReading,
		IsTankFull:      model.IsTankFull,
		HasMissedFillup: model.HasMissedFillup,
		Comments:        model.Comments,
		FillingStation:  model.FillingStation,
		UserID:          model.UserID,
		Date:            model.Date,
	}).Error
}

func UpdateExpense(fillupId string, model models.UpdateExpenseRequest) error {
	toUpdate, err := GetExpenseById(fillupId)
	if err != nil {
		return err
	}
	return db.DB.Model(&toUpdate).Updates(db.Expense{
		VehicleID:   model.VehicleID,
		Amount:      model.Amount,
		OdoReading:  model.OdoReading,
		ExpenseType: model.ExpenseType,
		Comments:    model.Comments,
		UserID:      model.UserID,
		Date:        model.Date,
	}).Error
}

func DeleteFillupById(fillupId string) error {
	return db.DeleteFillupById(fillupId)
}
func DeleteExpenseById(expenseId string) error {
	return db.DeleteExpenseById(expenseId)
}

func CreateVehicleAttachment(vehicleId, attachmentId, title string) error {
	model := &db.VehicleAttachment{
		AttachmentID: attachmentId,
		VehicleID:    vehicleId,
		Title:        title,
	}
	return db.DB.Create(model).Error
}
func GetVehicleAttachments(vehicleId string) (*[]db.Attachment, error) {

	return db.GetVehicleAttachments(vehicleId)
}

func GetUserStats(userId string, model models.UserStatsQueryModel) ([]models.VehicleStatsModel, error) {

	vehicles, err := GetUserVehicles(userId)
	if err != nil {
		return nil, err
	}

	var vehicleIds []string
	for _, v := range *vehicles {
		vehicleIds = append(vehicleIds, v.ID)
	}

	expenses, err := db.FindExpensesForDateRange(vehicleIds, model.Start, model.End)
	if err != nil {
		return nil, err
	}
	fillups, err := db.FindFillupsForDateRange(vehicleIds, model.Start, model.End)
	if err != nil {
		return nil, err
	}
	toReturn := models.VehicleStatsModel{}
	stats := toReturn.SetStats(fillups, expenses)

	return stats, nil
}
