package model

import "gorm.io/gorm"

type User struct{
	gorm.Model
	Id uint `gorm:"json":"id"`
	Name string `gorm:"json":"name"`
	Email string `gorm:"json":"email"`
	Password string `gorm:"json":"email"`
	Image string `gorm:"json":"image"`
}