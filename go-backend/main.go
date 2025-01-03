package main

import (
	"fmt"
	"net/http"

	db "go-backend/db"
	route "go-backend/routes"

	"github.com/gorilla/mux"
)




func main(){
	fmt.Println("Hello form the server :)")
	DbErr:=db.InitDB()
	if DbErr!=nil{
		fmt.Println("Error while connecting to the DB :(")
		return;
	}
	router:=mux.NewRouter()
	route.HandelRoutes(router)

	err:=http.ListenAndServe(":7000",router)
	if err!=nil{
		fmt.Println("Error while create a server on port 7000 :(")
		return;
	}
}