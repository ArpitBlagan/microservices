package controllers

import (
	"encoding/json"
	"fmt"
	db "go-backend/db"
	model "go-backend/models"
	"go-backend/utils"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"golang.org/x/crypto/bcrypt"
)

type createUser struct{
	Name string `json:"name"`
	Email string `json:"email"`
	Password string `json:"password"`
	Image string `json:"image"`
}

type changeStatusBody struct{
	Status string `json:"status"`
}

type loginInfo struct{
	Email string `json:"email"`
	Password string `json:"password"`
}

func HandleGetUser(w http.ResponseWriter,r *http.Request){
	params:=mux.Vars(r)
	id:=params["id"]
	var user model.User
	if err:=db.DB.First(&user,id).Error;err!=nil{
		fmt.Println(err)
		http.Error(w,"User not found :(",http.StatusBadRequest)
		return
	}
	res,err:=json.Marshal(user)
	if err!=nil{
		fmt.Println(err)
		http.Error(w,"Not able to marsh the user :(",http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-type","application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(res)
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
	//Bcrypt the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Failed to hash password", http.StatusInternalServerError)
		return
	}
	//Create a new user
	newUser :=createUser{
		Name:user.Name,
		Email:user.Email,
		Password:string(hashedPassword),
		Image:user.Image,
	}
	//Create user entry in DB
	if err := db.DB.Create(&newUser).Error; err != nil {
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

func HandleUserLogin(w http.ResponseWriter,r *http.Request){
	var userInfo loginInfo
	if err:=json.NewDecoder(r.Body).Decode(&userInfo);err!=nil{
		http.Error(w,"Not able to decode the req body.",http.StatusBadRequest)
		return
	}
	var user model.User
	// check in the DB if user is registered or not
	if err:=db.DB.Where("emaiL ?=",userInfo.Email).First(&user).Error;err!=nil{
		http.Error(w, "Invalid email or password", http.StatusUnauthorized)
		return
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(userInfo.Password)); err != nil {
		http.Error(w, "Invalid email or password", http.StatusUnauthorized)
		return
	}
	token,err:=utils.CreateJwt(user.ID)
	if err!=nil{
		http.Error(w,"",http.StatusInternalServerError)
		return 
	}
	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Value:    token,
		Expires:  time.Now().Add(24 * time.Hour),
		HttpOnly: true,
	})

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Login successful"))

}

func HandleGetRidesHistory(w http.ResponseWriter,r *http.Request){
	params:=mux.Vars(r)
	id:=params["id"]
	var user model.User
	//First find the user.
	if err:=db.DB.First(&user,id);err!=nil{
		http.Error(w,"User not found :(",http.StatusInternalServerError)
		return
	}
	var rides []model.Ride
	//Find the rides associated with the found userId.
	if err:=db.DB.Preload("User").Preload("Driver").Where("userId=?",user.ID).Find(&rides).Error; err!=nil{
		http.Error(w,"Error while finding the associated rides for the user :(",http.StatusInternalServerError)
		return
	}
	//If everythink is perfect return the rides associated with that user id.
	res,err:=json.Marshal(rides)
	if err!=nil{
		http.Error(w,"Error while marshaling the rides",http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-type","application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(res)
}

func HandleGetRide(w http.ResponseWriter,r *http.Request){
	params:=mux.Vars(r)
	id:=params["id"]
	var ride model.Ride
	if err:=db.DB.First(&ride,id).Error; err!=nil{
		http.Error(w,"Ride not found :(",http.StatusInternalServerError)
		return
	}
	res,err:=json.Marshal(ride)
	if err!=nil{
		http.Error(w,"Error while marshaling the ride",http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-type","application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(res)
}

func SearchRides(w http.ResponseWriter,r *http.Request){
	//Thinking how we gonna store the user address and geolocation 
	//and using it for getting the radius diameter for searching the rides
	// For particular area radius find the rides whose status are pending.
	var rides []model.Ride
	if err:=db.DB.Preload("User").Preload("Driver").Where("status=?","Pending").Find(&rides).Error;err!=nil{
		http.Error(w,"Rides not found :(",http.StatusInternalServerError)
		return
	}
	res,err:=json.Marshal(rides)
	if err!=nil{
		http.Error(w,"Error while marshaling the ride",http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-type","application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(res)
}

func ChangeRideStatus(w http.ResponseWriter,r *http.Request){
	params:=mux.Vars(r)
	rideId:=params["id"]
	var reqBody changeStatusBody
	if errr:=json.NewDecoder(r.Body).Decode(&reqBody).Error; errr!=nil{
		fmt.Println(errr)
		http.Error(w,"Error while decoding req body",http.StatusInternalServerError)
		return
	}
	var body changeStatusBody
	if err:=db.DB.Find(&body,rideId).Error; err!=nil{
		http.Error(w,"Ride not found :(",http.StatusInternalServerError)
		return
	}
	body.Status=reqBody.Status
	if err:=db.DB.Save(&body).Error;err!=nil{
		http.Error(w,"Failed to update the status of ride",http.StatusInternalServerError)
		return
	}
	res,err:=json.Marshal(body)
	if err!=nil{
		http.Error(w,"Error while marshaling the ride",http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-type","application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(res)
}