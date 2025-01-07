package main

import (
	"fmt"
	"net/http"

	db "go-backend/db"
	route "go-backend/routes"

	"github.com/gorilla/mux"
	"github.com/swaggo/swag/example/basic/docs"
)

// @title Simple Go backend API with Swagger
// @version 1.0
// @description A simple CRUD API example using Mux, Gorm and Swagger
// @host localhost:8080
// @BasePath /api/v1
func main(){
	fmt.Println("Hello form the server :)")
	DbErr:=db.InitDB()	
	if DbErr!=nil{
		fmt.Println("Error while connecting to the DB :(")
		return;
	}
	
	router:=mux.NewRouter()
	// Swagger documentation endpoint
	docs.SwaggerInfo.BasePath = "/api/v1"
	// router.Handle("/swagger/", ginSwagger.WrapHandler)
	route.HandelRoutes(router)

	err:=http.ListenAndServe(":8080",router)
	if err!=nil{
		fmt.Println("Error while create a server on port 7000 :(")
		return;
	}
}