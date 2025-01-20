package controllers

import (
	"encoding/json"
	"fmt"
	db "go-backend/db"
	model "go-backend/models"
	"go-backend/redis"
	"go-backend/utils"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"golang.org/x/crypto/bcrypt"
)

type contextKey string
const userKey contextKey = "user"

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
		// key string,value string, expirationTime time.Duration
		userID, ok := r.Context().Value(userKey).(string)
		if !ok{
			fmt.Println("Not able to fetch userID :(")
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
		redis.SetCache(userID,string(data),time.Minute*10)
		w.Write(data)
}

func HandleRegisterUser(w http.ResponseWriter,r *http.Request){
	//Extract the body from the request first
	var user model.User
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
	newUser :=model.User{
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
	data,err:=json.Marshal(newUser)
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
	if err:=db.DB.Where("email = ?",userInfo.Email).First(&user).Error;err!=nil{
		fmt.Println("Error while find user with give email in DB",err)
		http.Error(w, "Invalid email or password", http.StatusUnauthorized)
		return
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(userInfo.Password)); err != nil {
		http.Error(w, "Invalid email or password", http.StatusUnauthorized)
		return
	}
	fmt.Println(user.ID)
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
	fmt.Println(id)
	var user model.User
	//First find the user.
	if err:=db.DB.First(&user,id).Error;err!=nil{
		fmt.Println(err)
		http.Error(w,"User not found :(",http.StatusInternalServerError)
		return
	}
	var rides []model.Ride
	//Find the rides associated with the found userId.
	if err:=db.DB.Preload("User").Preload("Driver").Where("user_id =? ",id).Find(&rides).Error; err!=nil{
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

func HandleCreateRide(w http.ResponseWriter,r *http.Request){
	var ride model.Ride
	if err := json.NewDecoder(r.Body).Decode(&ride); err!=nil{
		fmt.Println(err)
		http.Error(w,"Error while decoding the body :(",http.StatusInternalServerError)
		return;
	}
	
	if err := db.DB.Create(&ride).Error; err!=nil{
		http.Error(w,"Not able to create a new ride :(",http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-type","application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Ride created sucessfully :)"))
}

func HandleCreateDriver(w http.ResponseWriter,r *http.Request){
	var driver model.Driver
	if err:=json.NewDecoder(r.Body).Decode(&driver);err!=nil{
		fmt.Println(err)
		http.Error(w,"Error while decoding the body :(",http.StatusInternalServerError)
		return;
	}
	fmt.Println(driver.CarId)

	var car model.Car
	if err := db.DB.First(&car, 1).Error; err != nil {
		// If the car is not found, return an error
		fmt.Println("Car not found with ID:", 1)
		http.Error(w, "Car not found for the given CarId", http.StatusBadRequest)
		return
	}
	driver.CarId=1;
	if err:=db.DB.Debug().Create(&driver).Error;err!=nil{
		fmt.Println(err)
		http.Error(w,"Not able to create a new Driver :(",http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-type","application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Driver created sucessfully :)"))
}

func HandleCreateCar(w http.ResponseWriter,r *http.Request){
	var car model.Car
	if err := json.NewDecoder(r.Body).Decode(&car); err!=nil{
		fmt.Println(err)
		http.Error(w,"Error while decoding the body :(",http.StatusInternalServerError)
		return;
	}
	
	if err := db.DB.Create(&car).Error; err!=nil{
		http.Error(w,"Not able to create a new car :(",http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-type","application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Car created sucessfully :)"))
}

func HandleGetCars(w http.ResponseWriter,r *http.Request){
	var cars []model.Car
	if err := db.DB.Find(&cars).Error; err!=nil{
		http.Error(w,"Not able to create a new car :(",http.StatusInternalServerError)
		return
	}
	res,errr := json.Marshal(cars);
	if errr!=nil{	
		http.Error(w,"Erorr while marshaling the result",http.StatusInternalServerError)
		return
	}	
	w.Header().Set("Content-type","application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(res)
}

func HandleGetDrivers(w http.ResponseWriter,r *http.Request){
	var drivers []model.Driver
	if err := db.DB.Find(&drivers).Error; err!=nil{
		http.Error(w,"Not able to create a new car :(",http.StatusInternalServerError)
		return
	}
	res,errr := json.Marshal(drivers);
	if errr!=nil{	
		http.Error(w,"Erorr while marshaling the result",http.StatusInternalServerError)
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

type RiDe struct{
	ride model.Ride
	distance float64
}

func SearchRides(w http.ResponseWriter,r *http.Request){
	//Thinking how we gonna store the user address and geolocation 
	//and using it for getting the radius diameter for searching the rides
	// For particular area radius find the rides whose status are pending.
	query:=r.URL.Query()

	latitude, err := strconv.ParseFloat(query.Get("latitude"), 64)
	longitude, err := strconv.ParseFloat(query.Get("longitude"), 64)
	radius, err := strconv.ParseFloat(query.Get("radius"), 64) // radius in kilometers

	if err != nil {
		http.Error(w, "Invalid query parameters", http.StatusBadRequest)
		return
	}

	var allRides []model.Ride
	db.DB.Find(&allRides)

	var nearbyRides []RiDe

	// Filter rides within the radius
	for _, ride := range allRides {
		distance := utils.Haversine(latitude, longitude, ride.PickupLatitude, ride.PickupLongitude)
		if distance <= radius {
			rideInfo:=RiDe{
				ride,
				distance,
			}
			nearbyRides = append(nearbyRides,rideInfo)
		}
	}

	response, err := json.Marshal(nearbyRides)
	if err != nil {
		http.Error(w, "Failed to marshal data", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(response)
}

func ChangeRideStatus(w http.ResponseWriter,r *http.Request){
	params:=mux.Vars(r)
	rideId:=params["id"]
	var reqBody changeStatusBody
	if errr:=json.NewDecoder(r.Body).Decode(&reqBody); errr!=nil{
		fmt.Println(errr)
		http.Error(w,"Error while decoding req body",http.StatusInternalServerError)
		return
	}
	var body model.Ride
	if err:=db.DB.First(&body,rideId).Error; err!=nil{
		http.Error(w,"Ride not found :(",http.StatusInternalServerError)
		return
	}
	fmt.Println(body.Status)
	body.Status=model.Status(reqBody.Status)
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