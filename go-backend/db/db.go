package db

import (
	"errors"
	"fmt"
	"log"

	model "go-backend/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB() error{
	var err error
	dsn := "host=localhost user=postgres password=yourpassword dbname=testdb port=5432 sslmode=disable"
	DB,err =gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if(err!=nil){
		fmt.Println("Error while connecting to DB",err)
		log.Fatal("Not able to connect to DB")
		return errors.New("Getting error while connecting to postgres")
	}
	DB.AutoMigrate(&model.User{})
	fmt.Println("Connect to DB sucessfully :)")
	return nil
}