package controllers

import (
	"encoding/json"
	"fmt"
	db "go-backend/db"
	model "go-backend/models"
	"net/http"
)

type createUser struct{
	Name string `json:"name"`
	Email string `json:"email"`
	Password string `json:"password"`
	Image string `json:"image"`
}

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

func HandleRegisterUser(w http.ResponseWriter,r *http.Request){
	//Extract the body from the request first
	var user createUser
	if err:=json.NewDecoder(r.Body).Decode(&user); err!=nil{
		fmt.Println(err)
		http.Error(w,"Error while decoding the req body :(",http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()
	//Create user entry in DB
	if err := db.DB.Create(&user).Error; err != nil {
		// If database insertion fails, return an error
		http.Error(w, "Failed to create user", http.StatusInternalServerError)
		fmt.Println("Error while creating user in the database:", err)
		return
	}
	data,err:=json.Marshal(user)
	if(err!=nil){
		http.Error(w,"Error while marshaling the user",http.StatusInternalServerError)
		fmt.Println("Not able to convert the data into json format")
		return
	}
	w.Header().Set("Content-Type","application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}