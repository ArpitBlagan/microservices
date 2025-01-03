package controllers

import (
	"encoding/json"
	"fmt"
	db "go-backend/db"
	model "go-backend/models"
	"net/http"
)

func HandleGetUsers(w http.ResponseWriter,r *http.Request){
	var users []model.User
		if(db.DB==nil){
			http.Error(w, "Error in DB connectivity", http.StatusInternalServerError)
			fmt.Println("DB is not there to run any type of operations")
			return
		}
		if err:=db.DB.Find(&users).Error; err!=nil{
			http.Error(w, "Error retrieving users", http.StatusInternalServerError)
		fmt.Println("Error while retrieving users:", err)
		return
		}
		w.Header().Set("Content-type","application/json")
		w.WriteHeader(http.StatusOK)
		data,err:=json.Marshal(users)
		if(err!=nil){
			http.Error(w, "Error retrieving users", http.StatusInternalServerError)
			fmt.Println("Error while marshaling the users into json",err)
			return;
		}
		w.Write(data)
}

// Define other controllers here....