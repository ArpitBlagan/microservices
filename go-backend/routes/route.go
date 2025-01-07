package route

import (
	"go-backend/controllers"

	"github.com/gorilla/mux"
)

func HandelRoutes(router *mux.Router){
	
	router.HandleFunc("/login",controllers.HandleUserLogin).Methods("POST")
	router.HandleFunc("/register",controllers.HandleRegisterUser).Methods("POST")

	//After this all route will first get validate.
	// router.Use(middleware.ValidateMiddleware)
	
	router.HandleFunc("/getUser",controllers.HandleGetUser).Methods("GET")
	router.HandleFunc("/getUsers",controllers.HandleGetUsers).Methods("GET")
	router.HandleFunc("/createDriver",controllers.HandleCreateDriver).Methods("POST")
	router.HandleFunc("/createCar",controllers.HandleCreateCar).Methods("POST")
	router.HandleFunc("/getCars",controllers.HandleGetCars).Methods("GET")
	router.HandleFunc("/get",controllers.HandleGetDrivers).Methods("GET")
	router.HandleFunc("/createRide",controllers.HandleCreateRide).Methods("POST")
	router.HandleFunc("/getRides/{id}",controllers.HandleGetRidesHistory).Methods("GET")
	router.HandleFunc("/searchRides",controllers.SearchRides).Methods("GET")
	router.HandleFunc("/getRide/{id}",controllers.HandleGetRide).Methods("GET")
	router.HandleFunc("/changeStatus/{id}",controllers.ChangeRideStatus).Methods("POST")
}