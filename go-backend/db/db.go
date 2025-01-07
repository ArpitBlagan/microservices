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

//If you are using docker simply run this command
//docker run -d -p 5432:5432 -e POSTGRES_PASSWORD=yourpassword -e POSTGRES_DB=textdb -e POSTGRES_USER=postgres -e postgres

func InitDB() error{
	var err error
	dsn := "host=localhost user=postgres password=yourpassword dbname=testdb port=5432 sslmode=disable"
	DB,err =gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if(err!=nil){
		fmt.Println("Error while connecting to DB",err)
		log.Fatal("Not able to connect to DB")
		panic("Something went wrong :(")
		return errors.New("Getting error while connecting to postgres")
	}
	DB.AutoMigrate(&model.User{},&model.Car{},&model.Driver{},&model.Ride{})
	fmt.Println("Connect to DB sucessfully :)")
	return nil
}