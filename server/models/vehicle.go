package models

import (
	"time"

	"github.com/akhilrex/hammond/db"
	_ "github.com/go-playground/validator/v10"
)

type SearchByIdQuery struct {
	Id string `binding:"required" uri:"id" json:"id" form:"id"`
}
type SubItemQuery struct {
	Id    string `binding:"required" uri:"id" json:"id" form:"id"`
	SubId string `binding:"required" uri:"subId" json:"subId" form:"subId"`
}
type CreateVehicleRequest struct {
	Nickname          string       `form:"nickname" json:"nickname" binding:"required"`
	Registration      string       `form:"registration" json:"registration" binding:"required"`
	Make              string       `form:"make" json:"make" binding:"required"`
	Model             string       `form:"model" json:"model" binding:"required"`
	YearOfManufacture int          `form:"yearOfManufacture" json:"yearOfManufacture"`
	EngineSize        float32      `form:"engineSize" json:"engineSize"`
	FuelUnit          *db.FuelUnit `form:"fuelUnit" json:"fuelUnit" binding:"required"`

	FuelType *db.FuelType `form:"fuelType" json:"fuelType" binding:"required"`
}

type UpdateVehicleRequest struct {
	CreateVehicleRequest
}
type UserVehicleSimpleModel struct {
	ID        string `json:"id"`
	UserID    string `json:"userId"`
	VehicleID string `json:"vehicleId"`
	IsOwner   bool   `json:"isOwner"`
	Name      string `json:"name"`
}

type CreateFillupRequest struct {
	VehicleID       string       `form:"vehicleId" json:"vehicleId" binding:"required"`
	FuelUnit        *db.FuelUnit `form:"fuelUnit" json:"fuelUnit" binding:"required"`
	FuelQuantity    float32      `form:"fuelQuantity" json:"fuelQuantity" binding:"required"`
	PerUnitPrice    float32      `form:"perUnitPrice" json:"perUnitPrice" binding:"required"`
	TotalAmount     float32      `form:"totalAmount" json:"totalAmount" binding:"required"`
	OdoReading      int          `form:"odoReading" json:"odoReading" binding:"required"`
	IsTankFull      *bool        `form:"isTankFull" json:"isTankFull" binding:"required"`
	HasMissedFillup *bool        `form:"hasMissedFillup" json:"HasMissedFillup"`
	Comments        string       `form:"comments" json:"comments" `
	FillingStation  string       `form:"fillingStation" json:"fillingStation"`
	UserID          string       `form:"userId" json:"userId" binding:"required"`
	Date            time.Time    `form:"date" json:"date" binding:"required" time_format:"2006-01-02"`
}

type UpdateFillupRequest struct {
	CreateFillupRequest
}

type UpdateExpenseRequest struct {
	CreateExpenseRequest
}

type CreateExpenseRequest struct {
	VehicleID string `form:"vehicleId" json:"vehicleId" binding:"required"`

	Amount     float32 `form:"amount" json:"amount" binding:"required"`
	OdoReading int     `form:"odoReading" json:"odoReading"`

	Comments    string    `form:"comments" json:"comments" `
	ExpenseType string    `form:"expenseType" json:"expenseType"`
	UserID      string    `form:"userId" json:"userId" binding:"required"`
	Date        time.Time `form:"date" json:"date" binding:"required" time_format:"2006-01-02"`
}

type CreateVehicleAttachmentModel struct {
	Title string `form:"title" json:"title" binding:"required"`
}

type VehicleStatsModel struct {
	CountFillups        int     `json:"countFillups"`
	CountExpenses       int     `json:"countExpenses"`
	ExpenditureFillups  float32 `json:"expenditureFillups"`
	ExpenditureExpenses float32 `json:"expenditureExpenses"`
	ExpenditureTotal    float32 `json:"expenditureTotal"`
	AvgFillupCost       float32 `json:"avgFillupCost"`
	AvgExpenseCost      float32 `json:"avgExpenseCost"`
	AvgFuelQty          float32 `json:"avgFuelQty"`
	AvgFuelPrice        float32 `json:"avgFuelPrice"`
	Currency            string  `json:"currency"`
}

func (m *VehicleStatsModel) SetStats(fillups *[]db.Fillup, expenses *[]db.Expense) []VehicleStatsModel {

	currencyMap := make(map[string]int)
	for _, v := range *fillups {
		currencyMap[v.Currency] = 1
	}
	for _, v := range *expenses {
		currencyMap[v.Currency] = 1
	}
	var toReturn []VehicleStatsModel
	for currency, _ := range currencyMap {
		model := VehicleStatsModel{}
		var totalExpenditure, fillupTotal, expenseTotal, totalFuel, averageFuelCost, averageFuelQty, averageFillup, averageExpense float32
		var countFillup, countExpense int
		for _, v := range *fillups {
			if v.Currency == currency {
				fillupTotal = fillupTotal + v.TotalAmount
				totalFuel = totalFuel + v.FuelQuantity
				countFillup++
			}
		}
		for _, v := range *expenses {
			if v.Currency == currency {
				expenseTotal = expenseTotal + v.Amount
				countExpense++
			}
		}

		totalExpenditure = expenseTotal + fillupTotal

		if countFillup > 0 {
			averageFillup = fillupTotal / float32(countFillup)
			averageFuelCost = fillupTotal / totalFuel
			averageFuelQty = totalFuel / float32(countFillup)
		}
		if countExpense > 0 {
			averageExpense = expenseTotal / float32(countExpense)
		}

		model.CountFillups = countFillup
		model.CountExpenses = countExpense
		model.ExpenditureFillups = fillupTotal
		model.ExpenditureExpenses = expenseTotal
		model.ExpenditureTotal = totalExpenditure
		model.AvgFillupCost = averageFillup
		model.AvgExpenseCost = averageExpense
		model.AvgFuelPrice = averageFuelCost
		model.AvgFuelQty = averageFuelQty
		model.Currency = currency

		toReturn = append(toReturn, model)
	}
	return toReturn
}

type UserStatsQueryModel struct {
	Start time.Time `json:"start" query:"start" form:"start"`
	End   time.Time `json:"end" query:"end" form:"end"`
}
