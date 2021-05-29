package db

import (
	"time"

	"gorm.io/driver/mysql"

	"gorm.io/gorm"
)

func CanMigrate(connectionString string) (bool, interface{}, error) {

	canInitialize, err := CanInitializeSystem()
	if !canInitialize {
		return canInitialize, nil, err
	}

	cdb, err := gorm.Open(mysql.Open(connectionString), &gorm.Config{})
	if err != nil {
		return false, nil, err
	}

	var usersCount, vehiclesCount, fuelCount int64
	tx := cdb.Table("Users").Count(&usersCount)
	if tx.Error != nil {
		return false, nil, tx.Error
	}
	tx = cdb.Table("Vehicles").Count(&vehiclesCount)
	if tx.Error != nil {
		return false, nil, tx.Error
	}
	tx = cdb.Table("Fuel").Count(&fuelCount)
	if tx.Error != nil {
		return false, nil, tx.Error
	}
	data := struct {
		Users    int64 `json:"users"`
		Vehicles int64 `json:"vehicles"`
		Fillups  int64 `json:"fillups"`
	}{
		Vehicles: vehiclesCount,
		Users:    usersCount,
		Fillups:  fuelCount,
	}

	return true, data, nil
}

func MigrateClarkson(connectionString string) (bool, error) {
	canInitialize, err := CanInitializeSystem()
	if !canInitialize {
		return canInitialize, err
	}

	//dsn := "user:pass@tcp(127.0.0.1:3306)/dbname?charset=utf8mb4&parseTime=True&loc=Local"
	cdb, err := gorm.Open(mysql.Open(connectionString), &gorm.Config{})
	if err != nil {
		return false, nil
	}

	/////Models
	type CUser struct {
		ID                  string `gorm:"column:id"`
		Email               string `gorm:"column:email"`
		Username            string `gorm:"column:username"`
		Password            string `gorm:"column:password"`
		Admin               bool   `gorm:"column:admin"`
		FuelUnit            int    `gorm:"column:fuelUnit"`
		DistanceUnit        int    `gorm:"column:distanceUnit"`
		FuelConsumptionUnit int    `gorm:"column:fuelConsumptionUnit"`
		CurrencyUnit        int    `gorm:"column:currencyUnit"`
	}

	type CVehicle struct {
		ID                string `gorm:"column:id"`
		User              string
		Name              string
		Registration      string
		Make              string
		Model             string
		YearOfManufacture int `gorm:"column:yearOfManufacture"`
		Vin               string
		EngineSizeCC      int `gorm:"column:engineSizeCC"`
		FuelType          int `gorm:"column:fuelType"`
	}

	type CFuel struct {
		ID              string    `gorm:"column:id"`
		Vehicle         string    `gorm:"column:vehicle"`
		Date            time.Time `gorm:"column:date"`
		FuelAmount      float32   `gorm:"column:fuelAmount"`
		TotalCost       float32   `gorm:"column:totalCost"`
		FuelUnitCost    float32   `gorm:"column:fuelUnitCost"`
		OdometerReading int       `gorm:"column:odometerReading"`
		Notes           string    `gorm:"column:notes"`
		FullTank        bool      `gorm:"column:fullTank"`
		MissedFillup    bool      `gorm:"column:missedFillUp"`
	}

	distanceUnitMap := map[int]DistanceUnit{
		1: MILES,
		2: KILOMETERS,
	}

	fuelTypeMap := map[int]FuelType{
		1: PETROL,
		2: DIESEL,
		3: ETHANOL,
		4: LPG,
	}

	fuelUnitsMap := map[int]FuelUnit{
		1: LITRE,
		2: GALLON,
		3: US_GALLON,
	}
	currencyMap := map[int]string{
		1: "GBP",
		2: "USD",
		3: "EUR",
		4: "AUD",
		5: "CAD",
	}

	newUserIdsMap := make(map[string]User)
	oldUserIdsMap := make(map[string]CUser)

	var allUsers []CUser
	cdb.Table("Users").Find(&allUsers)
	for _, v := range allUsers {
		role := USER
		if v.Admin {
			role = ADMIN
		}
		user := User{
			Email:        v.Email,
			Currency:     currencyMap[v.CurrencyUnit],
			DistanceUnit: distanceUnitMap[v.DistanceUnit],
			Role:         role,
			Name:         v.Username,
		}
		user.SetPassword("hammond")
		err = CreateUser(&user)
		if err != nil {
			return false, err
		}

		newUserIdsMap[v.ID] = user
		oldUserIdsMap[v.ID] = v

		if v.Admin {
			setting := GetOrCreateSetting()
			setting.Currency = user.Currency
			setting.DistanceUnit = user.DistanceUnit
			UpdateSettings(setting)
		}
	}

	newVehicleIdsMap := make(map[string]Vehicle)
	oldVehicleIdsMap := make(map[string]CVehicle)
	vehicleUserMap := make(map[string]User)
	var allVehicles []CVehicle
	cdb.Table("Vehicles").Find(&allVehicles)
	for _, model := range allVehicles {
		vehicle := Vehicle{
			Nickname:          model.Name,
			Registration:      model.Registration,
			Model:             model.Model,
			Make:              model.Make,
			YearOfManufacture: model.YearOfManufacture,
			EngineSize:        float32(model.EngineSizeCC),
			FuelUnit:          fuelUnitsMap[oldUserIdsMap[model.User].FuelUnit],
			FuelType:          fuelTypeMap[model.FuelType],
		}

		tx := DB.Create(&vehicle)
		if tx.Error != nil {
			return false, tx.Error
		}
		association := UserVehicle{
			UserID:    newUserIdsMap[model.User].ID,
			VehicleID: vehicle.ID,
			IsOwner:   true,
		}
		vehicleUserMap[vehicle.ID] = newUserIdsMap[model.User]
		tx = DB.Create(&association)

		if tx.Error != nil {
			return false, err
		}

		newVehicleIdsMap[model.ID] = vehicle
		oldVehicleIdsMap[model.ID] = model
	}

	var allFillups []CFuel
	cdb.Table("Fuel").Find(&allFillups)
	for _, model := range allFillups {
		fillup := Fillup{
			VehicleID:       newVehicleIdsMap[model.Vehicle].ID,
			FuelUnit:        newVehicleIdsMap[model.Vehicle].FuelUnit,
			FuelQuantity:    model.FuelAmount,
			PerUnitPrice:    model.FuelUnitCost,
			TotalAmount:     model.TotalCost,
			OdoReading:      model.OdometerReading,
			IsTankFull:      &model.FullTank,
			HasMissedFillup: &model.MissedFillup,
			Comments:        model.Notes,
			UserID:          vehicleUserMap[newVehicleIdsMap[model.Vehicle].ID].ID,
			Date:            model.Date,
			Currency:        vehicleUserMap[newVehicleIdsMap[model.Vehicle].ID].Currency,
			DistanceUnit:    vehicleUserMap[newVehicleIdsMap[model.Vehicle].ID].DistanceUnit,
		}

		tx := DB.Create(&fillup)
		if tx.Error != nil {
			return false, tx.Error
		}

	}

	return true, nil
}
