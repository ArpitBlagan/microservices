package route

import (
	"go-backend/controllers"

	"github.com/gorilla/mux"
)

func HandelRoutes(router *mux.Router){
	router.HandleFunc("/getUser",controllers.HandleGetUsers).Methods("GET")
}