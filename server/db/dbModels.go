package db

import (
	"encoding/json"
	"errors"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	Base
	Email        string       `gorm:"unique" json:"email"`
	Password     string       `json:"-"`
	Currency     string       `json:"currency"`
	DistanceUnit DistanceUnit `json:"distanceUnit"`
	DateFormat   string       `json:"dateFormat"`
	Role         Role         `json:"role"`
	Name         string       `json:"name"`
	Vehicles     []Vehicle    `gorm:"many2many:user_vehicles;" json:"vehicles"`
	IsDisabled   bool         `json:"isDisabled"`
}

func (b *User) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		User
		RoleDetail         EnumDetail `json:"roleDetail"`
		DistanceUnitDetail EnumDetail `json:"distanceUnitDetail"`
	}{
		User:               *b,
		RoleDetail:         b.RoleDetail(),
		DistanceUnitDetail: b.DistanceUnitDetail(),
	})
}
func (v *User) RoleDetail() EnumDetail {
	return RoleDetails[v.Role]
}
func (v *User) DistanceUnitDetail() EnumDetail {
	return DistanceUnitDetails[v.DistanceUnit]
}

func (u *User) SetPassword(password string) error {
	if len(password) == 0 {
		return errors.New("password should not be empty")
	}
	bytePassword := []byte(password)
	// Make sure the second param `bcrypt generator cost` between [4, 32)
	passwordHash, _ := bcrypt.GenerateFromPassword(bytePassword, bcrypt.DefaultCost)
	u.Password = string(passwordHash)
	return nil
}

func (u *User) CheckPassword(password string) error {
	bytePassword := []byte(password)
	byteHashedPassword := []byte(u.Password)
	return bcrypt.CompareHashAndPassword(byteHashedPassword, bytePassword)
}

type Vehicle struct {
	Base
	Nickname          string       `json:"nickname"`
	Registration      string       `json:"registration"`
	Make              string       `json:"make"`
	Model             string       `json:"model"`
	YearOfManufacture int          `json:"yearOfManufacture"`
	EngineSize        float32      `json:"engineSize"`
	FuelUnit          FuelUnit     `json:"fuelUnit"`
	FuelType          FuelType     `json:"fuelType"`
	Users             []User       `gorm:"many2many:user_vehicles;" json:"users"`
	Fillups           []Fillup     `json:"fillups"`
	Expenses          []Expense    `json:"expenses"`
	Attachments       []Attachment `gorm:"many2many:vehicle_attachments;" json:"attachments"`
	IsOwner           bool         `gorm:"->" json:"isOwner"`
}

func (b *Vehicle) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Vehicle
		FuelTypeDetail EnumDetail `json:"fuelTypeDetail"`
		FuelUnitDetail EnumDetail `json:"fuelUnitDetail"`
	}{
		Vehicle:        *b,
		FuelTypeDetail: b.FuelTypeDetail(),
		FuelUnitDetail: b.FuelUnitDetail(),
	})
}
func (v *Vehicle) FuelTypeDetail() EnumDetail {
	return FuelTypeDetails[v.FuelType]
}

func (v *Vehicle) FuelUnitDetail() EnumDetail {
	return FuelUnitDetails[v.FuelUnit]
}

type UserVehicle struct {
	Base
	UserID    string `gorm:"primaryKey"`
	User      User   `json:"user"`
	VehicleID string `gorm:"primaryKey"`
	IsOwner   bool   `json:"isOwner"`
}

type Fillup struct {
	Base
	VehicleID       string       `json:"vehicleId"`
	Vehicle         Vehicle      `json:"-"`
	FuelUnit        FuelUnit     `json:"fuelUnit"`
	FuelQuantity    float32      `json:"fuelQuantity"`
	PerUnitPrice    float32      `json:"perUnitPrice"`
	TotalAmount     float32      `json:"totalAmount"`
	OdoReading      int          `json:"odoReading"`
	IsTankFull      *bool        `json:"isTankFull"`
	HasMissedFillup *bool        `json:"hasMissedFillup"`
	Comments        string       `json:"comments"`
	FillingStation  string       `json:"fillingStation"`
	UserID          string       `json:"userId"`
	User            User         `json:"user"`
	Date            time.Time    `json:"date"`
	Currency        string       `json:"currency"`
	DistanceUnit    DistanceUnit `json:"distanceUnit"`
	Source          string       `json:"source"`
}

func (v *Fillup) FuelUnitDetail() EnumDetail {
	return FuelUnitDetails[v.FuelUnit]
}
func (b *Fillup) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Fillup
		FuelUnitDetail EnumDetail `json:"fuelUnitDetail"`
	}{
		Fillup:         *b,
		FuelUnitDetail: b.FuelUnitDetail(),
	})
}

type Expense struct {
	Base
	VehicleID    string       `json:"vehicleId"`
	Vehicle      Vehicle      `json:"-"`
	Amount       float32      `json:"amount"`
	OdoReading   int          `json:"odoReading"`
	Comments     string       `json:"comments"`
	ExpenseType  string       `json:"expenseType"`
	UserID       string       `json:"userId"`
	User         User         `json:"user"`
	Date         time.Time    `json:"date"`
	Currency     string       `json:"currency"`
	DistanceUnit DistanceUnit `json:"distanceUnit"`
	Source       string       `json:"source"`
}

type Setting struct {
	Base
	Currency     string       `json:"currency" gorm:"default:INR"`
	DistanceUnit DistanceUnit `json:"distanceUnit" gorm:"default:1"`
}
type Migration struct {
	Base
	Date time.Time
	Name string
}
type JobLock struct {
	Base
	Date     time.Time
	Name     string
	Duration int
}

type Attachment struct {
	Base
	Path         string `json:"path"`
	OriginalName string `json:"originalName"`
	Size         int64  `json:"size"`
	ContentType  string `json:"contentType"`
	Title        string `gorm:"->" json:"title"`
	UserID       string `json:"userId"`
	User         User   `json:"user"`
}

type QuickEntry struct {
	Base
	AttachmentID string     `json:"attachmentId"`
	Attachment   Attachment `json:"attachment"`
	ProcessDate  *time.Time `json:"processDate"`
	UserID       string     `json:"userId"`
	User         User       `json:"user"`
	Comments     string     `json:"comments"`
}

type VehicleAttachment struct {
	Base
	AttachmentID string `gorm:"primaryKey" json:"attachmentId"`
	VehicleID    string `gorm:"primaryKey" json:"vehicleId"`
	Title        string `json:"title"`
}
