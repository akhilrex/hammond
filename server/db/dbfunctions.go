package db

import (
	"errors"
	"fmt"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func CanInitializeSystem() (bool, error) {
	users, _ := GetAllUsers()
	if len(*users) != 0 {
		// db.MigrateClarkson("root:password@tcp(192.168.0.117:3306)/clarkson?charset=utf8mb4&parseTime=True&loc=Local")
		return false,
			fmt.Errorf("there are already users in the database. Migration can only be done on an empty database")
	}
	return true, nil
}

func CreateUser(user *User) error {
	tx := DB.Create(&user)
	return tx.Error
}
func UpdateUser(user *User) error {
	tx := DB.Omit(clause.Associations).Save(&user)
	return tx.Error
}
func FindOneUser(condition interface{}) (User, error) {

	var model User
	err := DB.Where(condition).First(&model).Error
	return model, err
}
func SetDisabledStatusForUser(userId string, isDisabled bool) error {
	//Cannot do this for admin
	tx := DB.Debug().Model(&User{}).Where("id= ? and role=?", userId, USER).Update("is_disabled", isDisabled)
	return tx.Error
}
func GetAllUsers() (*[]User, error) {

	sorting := "created_at desc"
	var users []User
	result := DB.Order(sorting).Find(&users)
	return &users, result.Error
}

func GetAllVehicles(sorting string) (*[]Vehicle, error) {
	if sorting == "" {
		sorting = "created_at desc"
	}
	var vehicles []Vehicle
	result := DB.Preload("Fillups", func(db *gorm.DB) *gorm.DB {
		return db.Order("fillups.date DESC")
	}).Preload("Expenses", func(db *gorm.DB) *gorm.DB {
		return db.Order("expenses.date DESC")
	}).Order(sorting).Find(&vehicles)
	return &vehicles, result.Error
}

func GetVehicleOwner(vehicleId string) (string, error) {
	var mapping UserVehicle

	tx := DB.Where("vehicle_id = ? AND is_owner = 1", vehicleId).First(&mapping)

	if tx.Error != nil {
		return "", tx.Error
	}
	return mapping.UserID, nil
}

func GetVehicleUsers(vehicleId string) (*[]UserVehicle, error) {
	var mapping []UserVehicle

	tx := DB.Debug().Preload("User").Where("vehicle_id = ?", vehicleId).Find(&mapping)

	if tx.Error != nil {
		return nil, tx.Error
	}
	return &mapping, nil
}

func ShareVehicle(vehicleId, userId string) error {
	var mapping UserVehicle

	tx := DB.Where("vehicle_id = ? AND user_id = ?", vehicleId, userId).First(&mapping)

	if errors.Is(tx.Error, gorm.ErrRecordNotFound) {
		newMapping := UserVehicle{
			UserID:    userId,
			VehicleID: vehicleId,
			IsOwner:   false,
		}
		tx = DB.Create(&newMapping)
		return tx.Error
	}
	return nil
}
func TransferVehicle(vehicleId, ownerId, newUserID string) error {

	tx := DB.Model(&UserVehicle{}).Where("vehicle_id = ? AND user_id = ?", vehicleId, ownerId).Update("is_owner", false)
	if tx.Error != nil {
		return tx.Error
	}
	tx = DB.Model(&UserVehicle{}).Where("vehicle_id = ? AND user_id = ?", vehicleId, newUserID).Update("is_owner", true)

	return tx.Error
}

func UnshareVehicle(vehicleId, userId string) error {
	var mapping UserVehicle

	tx := DB.Where("vehicle_id = ? AND user_id = ?", vehicleId, userId).First(&mapping)

	if errors.Is(tx.Error, gorm.ErrRecordNotFound) {
		return nil
	}
	if mapping.IsOwner {
		return fmt.Errorf("Cannot unshare owner")
	}
	result := DB.Where("id=?", mapping.ID).Delete(&UserVehicle{})
	return result.Error
}

func GetUserVehicles(id string) (*[]Vehicle, error) {
	var toReturn []Vehicle
	user, err := GetUserById(id)
	if err != nil {
		return nil, err
	}
	err = DB.Preload("Fillups", func(db *gorm.DB) *gorm.DB {
		return db.Order("fillups.date DESC")
	}).Preload("Expenses", func(db *gorm.DB) *gorm.DB {
		return db.Order("expenses.date DESC")
	}).Model(user).Select("vehicles.*,user_vehicles.is_owner").Association("Vehicles").Find(&toReturn)
	if err != nil {
		return nil, err
	}
	return &toReturn, nil
}
func GetUserById(id string) (*User, error) {
	var data User
	result := DB.Preload(clause.Associations).First(&data, "id=?", id)
	return &data, result.Error
}
func GetVehicleById(id string) (*Vehicle, error) {
	var vehicle Vehicle
	result := DB.Preload(clause.Associations).First(&vehicle, "id=?", id)
	return &vehicle, result.Error
}
func GetFillupById(id string) (*Fillup, error) {
	var obj Fillup
	result := DB.Preload(clause.Associations).First(&obj, "id=?", id)
	return &obj, result.Error
}

func GetFillupsByVehicleId(id string) (*[]Fillup, error) {
	var obj []Fillup
	result := DB.Preload(clause.Associations).Order("date desc").Find(&obj, &Fillup{VehicleID: id})
	return &obj, result.Error
}
func GetLatestFillupsByVehicleId(id string) (*Fillup, error) {
	var obj Fillup
	result := DB.Preload(clause.Associations).Order("date desc").First(&obj, &Fillup{VehicleID: id})
	return &obj, result.Error
}
func GetFillupsByVehicleIdSince(id string, since time.Time) (*[]Fillup, error) {
	var obj []Fillup
	result := DB.Where("date >= ? AND vehicle_id = ?", since, id).Preload(clause.Associations).Order("date desc").Find(&obj)
	return &obj, result.Error
}
func FindFillups(condition interface{}) (*[]Fillup, error) {

	var model []Fillup
	err := DB.Where(condition).Find(&model).Error
	return &model, err
}

func FindFillupsForDateRange(vehicleIds []string, start, end time.Time) (*[]Fillup, error) {

	var model []Fillup
	err := DB.Where("date <= ? AND date >= ? AND vehicle_id in ?", end, start, vehicleIds).Find(&model).Error
	return &model, err
}
func FindExpensesForDateRange(vehicleIds []string, start, end time.Time) (*[]Expense, error) {

	var model []Expense
	err := DB.Where("date <= ? AND date >= ? AND vehicle_id in ?", end, start, vehicleIds).Find(&model).Error
	return &model, err
}

func GetExpensesByVehicleId(id string) (*[]Expense, error) {
	var obj []Expense
	result := DB.Preload(clause.Associations).Order("date desc").Find(&obj, &Expense{VehicleID: id})
	return &obj, result.Error
}
func GetLatestExpenseByVehicleId(id string) (*Expense, error) {
	var obj Expense
	result := DB.Preload(clause.Associations).Order("date desc").First(&obj, &Expense{VehicleID: id})
	return &obj, result.Error
}
func GetExpenseById(id string) (*Expense, error) {
	var obj Expense
	result := DB.Preload(clause.Associations).First(&obj, "id=?", id)
	return &obj, result.Error
}

func DeleteVehicleById(id string) error {

	result := DB.Where("id=?", id).Delete(&Vehicle{})
	return result.Error
}
func DeleteFillupById(id string) error {

	result := DB.Where("id=?", id).Delete(&Fillup{})
	return result.Error
}
func DeleteExpenseById(id string) error {
	result := DB.Where("id=?", id).Delete(&Expense{})
	return result.Error
}

func DeleteFillupByVehicleId(id string) error {

	result := DB.Where("vehicle_id=?", id).Delete(&Fillup{})
	return result.Error
}
func DeleteExpenseByVehicleId(id string) error {
	result := DB.Where("vehicle_id=?", id).Delete(&Expense{})
	return result.Error
}

func GetAllQuickEntries(sorting string) (*[]QuickEntry, error) {
	if sorting == "" {
		sorting = "created_at desc"
	}
	var quickEntries []QuickEntry
	result := DB.Preload(clause.Associations).Order(sorting).Find(&quickEntries)
	return &quickEntries, result.Error
}
func GetQuickEntriesForUser(userId, sorting string) (*[]QuickEntry, error) {
	if sorting == "" {
		sorting = "created_at desc"
	}
	var quickEntries []QuickEntry
	result := DB.Preload(clause.Associations).Where("user_id = ?", userId).Order(sorting).Find(&quickEntries)
	return &quickEntries, result.Error
}
func GetQuickEntryById(id string) (*QuickEntry, error) {
	var quickEntry QuickEntry
	result := DB.Preload(clause.Associations).First(&quickEntry, "id=?", id)
	return &quickEntry, result.Error
}
func DeleteQuickEntryById(id string) error {
	result := DB.Where("id=?", id).Delete(&QuickEntry{})
	return result.Error
}
func UpdateQuickEntry(entry *QuickEntry) error {
	return DB.Save(entry).Error
}
func SetQuickEntryAsProcessed(id string, processDate time.Time) error {
	result := DB.Model(QuickEntry{}).Where("id=?", id).Update("process_date", processDate)
	return result.Error
}

func GetAttachmentById(id string) (*Attachment, error) {
	var entry Attachment
	result := DB.Preload(clause.Associations).First(&entry, "id=?", id)
	return &entry, result.Error
}
func GetVehicleAttachments(vehicleId string) (*[]Attachment, error) {
	var attachments []Attachment
	vehicle, err := GetVehicleById(vehicleId)
	if err != nil {
		return nil, err
	}
	err = DB.Debug().Model(vehicle).Select("attachments.*,vehicle_attachments.title").Preload("User").Association("Attachments").Find(&attachments)
	if err != nil {
		return nil, err
	}
	return &attachments, nil
}
func GeAlertById(id string) (*VehicleAlert, error) {
	var alert VehicleAlert
	result := DB.Preload(clause.Associations).First(&alert, "id=?", id)
	return &alert, result.Error
}
func GetAlertOccurenceByAlertId(id string) (*[]AlertOccurance, error) {
	var alertOccurance []AlertOccurance
	result := DB.Preload(clause.Associations).Order("created_at desc").Find(&alertOccurance, "vehicle_alert_id=?", id)
	return &alertOccurance, result.Error
}

func GetUnprocessedAlertOccurances() (*[]AlertOccurance, error) {
	var alertOccurance []AlertOccurance
	result := DB.Preload(clause.Associations).Order("created_at desc").Find(&alertOccurance, "process_date is NULL")
	return &alertOccurance, result.Error
}
func MarkAlertOccuranceAsProcessed(id string, alertProcessType AlertType, date time.Time) error {
	tx := DB.Debug().Model(&AlertOccurance{}).Where("id= ?", id).
		Update("alert_process_type", alertProcessType).
		Update("process_date", date)
	return tx.Error

}

func UpdateSettings(setting *Setting) error {
	tx := DB.Save(&setting)
	return tx.Error
}
func GetOrCreateSetting() *Setting {
	var setting Setting
	result := DB.First(&setting)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		DB.Save(&Setting{})
		DB.First(&setting)
	}
	return &setting
}

func GetLock(name string) *JobLock {
	var jobLock JobLock
	result := DB.Where("name = ?", name).First(&jobLock)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return &JobLock{
			Name: name,
		}
	}
	return &jobLock
}
func Lock(name string, duration int) {
	jobLock := GetLock(name)
	if jobLock == nil {
		jobLock = &JobLock{
			Name: name,
		}
	}
	jobLock.Duration = duration
	jobLock.Date = time.Now()
	if jobLock.ID == "" {
		DB.Create(&jobLock)
	} else {
		DB.Save(&jobLock)
	}
}
func Unlock(name string) {
	jobLock := GetLock(name)
	if jobLock == nil {
		return
	}
	jobLock.Duration = 0
	jobLock.Date = time.Time{}
	DB.Save(&jobLock)
}

func UnlockMissedJobs() {
	var jobLocks []JobLock

	result := DB.Find(&jobLocks)
	if result.Error != nil {
		return
	}
	for _, job := range jobLocks {
		if (job.Date == time.Time{}) {
			continue
		}
		var duration time.Duration
		duration = time.Duration(job.Duration)
		d := job.Date.Add(time.Minute * duration)
		if d.Before(time.Now()) {
			fmt.Println(job.Name + " is unlocked")
			Unlock(job.Name)
		}
	}
}
