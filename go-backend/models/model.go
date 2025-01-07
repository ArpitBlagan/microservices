package model

import "gorm.io/gorm"

// User Model
type User struct {
	gorm.Model
	Name     string `json:"name"`
	Email    string `json:"email" gorm:"uniqueIndex;`
	Password string `json:"password"`
	Image    string `json:"image"`
}

// Car Model
type Car struct {
	gorm.Model
	Name   string `json:"name"`
	Number string `json:"number"`
}

// Driver Model
type Driver struct {
	gorm.Model
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Image    string `json:"image"`
	CarId    uint   `json:"car_id"`
	Car      Car    `gorm:"foreignKey:CarId" json:"car"`
}

// Status Enum
type Status string

const (
	Pending   Status = "Pending"
	Confirm   Status = "Confirm"
	Cancel    Status = "Cancel"
	Completed Status = "Completed"
)

// Ride Model
type Ride struct {
	gorm.Model
	Pickup      string `json:"pickup"`
	Destination string `json:"destination"`
	PickupLatitude    float64 `json:"pickup_latitude"`
    PickupLongitude   float64 `json:"pickup_longitude"`
    DestinationLatitude float64 `json:"destination_latitude"`
    DestinationLongitude float64 `json:"destination_longitude"`
	DriverId    uint   `json:"driver_id"`
	Driver      Driver `gorm:"foreignKey:DriverId" json:"driver"`
	UserID      uint   `json:"user_id"` // Foreign key for User
	User        User   `gorm:"foreignKey:UserID" json:"user"`
	Status      Status `json:"status"`
}
