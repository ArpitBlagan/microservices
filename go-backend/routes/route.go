package route

import (
	"go-backend/controllers"
	// "go-backend/middleware"
	"go-backend/redis"

	"github.com/gorilla/mux"
)

func HandleRoutes(router *mux.Router) {
	// Routes without middleware
	router.HandleFunc("/login", controllers.HandleUserLogin).Methods("POST")
	router.HandleFunc("/register", controllers.HandleRegisterUser).Methods("POST")

	// Create a subrouter for routes that require middleware
	authRoutes := router.PathPrefix("/").Subrouter()
	// authRoutes.Use(middleware.ValidateMiddleware)
	authRoutes.Use(redis.RateLimitRequest)
	authRoutes.Use(redis.GetCache)

	// Define routes with middleware
	authRoutes.HandleFunc("/getUser", controllers.HandleGetUser).Methods("GET")
	authRoutes.HandleFunc("/getUsers", controllers.HandleGetUsers).Methods("GET")
	authRoutes.HandleFunc("/createDriver", controllers.HandleCreateDriver).Methods("POST")
	authRoutes.HandleFunc("/createCar", controllers.HandleCreateCar).Methods("POST")
	authRoutes.HandleFunc("/getCars", controllers.HandleGetCars).Methods("GET")
	authRoutes.HandleFunc("/get", controllers.HandleGetDrivers).Methods("GET")
	authRoutes.HandleFunc("/createRide", controllers.HandleCreateRide).Methods("POST")
	authRoutes.HandleFunc("/getRides/{id}", controllers.HandleGetRidesHistory).Methods("GET")
	authRoutes.HandleFunc("/searchRides", controllers.SearchRides).Methods("GET")
	authRoutes.HandleFunc("/getRide/{id}", controllers.HandleGetRide).Methods("GET")
	authRoutes.HandleFunc("/changeStatus/{id}", controllers.ChangeRideStatus).Methods("POST")
}
