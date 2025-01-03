package route

import (
	"go-backend/controllers"
	"go-backend/middleware"

	"github.com/gorilla/mux"
)

func HandelRoutes(router *mux.Router){
	
	router.HandleFunc("/login",controllers.HandleUserLogin).Methods("POST")
	router.HandleFunc("/register",controllers.HandleRegisterUser).Methods("POST")
	//After this all route will first get validate.
	router.Use(middleware.ValidateMiddleware)
	router.HandleFunc("/getUser",controllers.HandleGetUser).Methods("GET")
	router.HandleFunc("/getUsers/{id}",controllers.HandleGetUsers).Methods("GET")
	router.HandleFunc("/getRides/{id}",controllers.HandleGetRidesHistory).Methods("GET")
	router.HandleFunc("/getRide/{id}",controllers.HandleGetRide).Methods("GET")
	router.HandleFunc("/changeStatus/{id}",controllers.ChangeRideStatus).Methods("POST")
}