package db

import (
	"fmt"
	"log"
	"os"
	"path"

	"gorm.io/driver/sqlite"

	"gorm.io/gorm"
)

//DB is
var DB *gorm.DB

//Init is used to Initialize Database
func Init() (*gorm.DB, error) {
	// github.com/mattn/go-sqlite3
	configPath := os.Getenv("CONFIG")
	dbPath := path.Join(configPath, "hammond.db")
	log.Println(dbPath)
	db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true,
	})
	if err != nil {
		fmt.Println("db err: ", err)
		return nil, err
	}

	localDB, _ := db.DB()
	localDB.SetMaxIdleConns(10)
	//db.LogMode(true)
	DB = db
	return DB, nil
}

//Migrate Database
func Migrate() {
	err := DB.AutoMigrate(&Attachment{}, &QuickEntry{}, &User{}, &Vehicle{}, &UserVehicle{}, &VehicleAttachment{}, &Fillup{}, &Expense{}, &Setting{}, &JobLock{}, &Migration{})
	if err != nil {
		fmt.Println("1 " + err.Error())
	}
	err = DB.SetupJoinTable(&User{}, "Vehicles", &UserVehicle{})
	if err != nil {
		fmt.Println(err.Error())
	}
	err = DB.SetupJoinTable(&Vehicle{}, "Attachments", &VehicleAttachment{})
	if err != nil {
		fmt.Println(err.Error())
	}
	RunMigrations()
}

// Using this function to get a connection, you can create your connection pool here.
func GetDB() *gorm.DB {
	return DB
}
