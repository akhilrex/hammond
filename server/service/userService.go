package service

import (
	"strings"

	"github.com/akhilrex/hammond/db"
	"github.com/akhilrex/hammond/models"
)

func CreateUser(userModel *models.RegisterRequest, role db.Role) error {
	setting := db.GetOrCreateSetting()
	toCreate := db.User{
		Email:        strings.ToLower(userModel.Email),
		Name:         userModel.Name,
		Role:         role,
		Currency:     setting.Currency,
		DistanceUnit: setting.DistanceUnit,
		DateFormat:   "MM/dd/yyyy",
	}

	toCreate.SetPassword(userModel.Password)

	return db.CreateUser(&toCreate)

}

func GetUserById(id string) (*db.User, error) {
	var myUserModel db.User
	tx := db.DB.Debug().Preload("Vehicles").First(&myUserModel, map[string]string{
		"ID": id,
	})
	return &myUserModel, tx.Error
}

func GetAllUsers() (*[]db.User, error) {
	return db.GetAllUsers()
}

func UpdatePassword(id, password string) (bool, error) {
	user, err := GetUserById(id)
	if err != nil {
		return false, err
	}
	user.SetPassword(password)
	err = db.UpdateUser(user)
	if err != nil {
		return false, err
	}
	return true, nil
}
func SetDisabledStatusForUser(userId string, isDisabled bool) error {
	return db.SetDisabledStatusForUser(userId, isDisabled)
}
