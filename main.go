package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
	"os"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var db *gorm.DB

type User struct {
	ID       uint   `gorm:"primaryKey"`
	Username string `gorm:"unique"`
	Name     string
	Surname  string
	Password string
	Email    string
	Age      int
	IsAdmin  bool
	Meetups []*Meetup `gorm:"many2many:user_meetups;"`
}

type Meetup struct {
	ID                  uint `gorm:"primaryKey;autoIncrement:true"`
	MeetupName          string
	MeetupDetails       string
	StartDate           time.Time
	EndDate             time.Time
	Address             string
	Quota               int
	RegisteredUserCount int
	Users []*User `gorm:"many2many:user_meetups;"`
}

// POST Method
// Creates and saves a new user on the database
// @returns Saved user data in JSON if it is saved successfully
// @returns Error message if it is not saved successfully
// TODO Error handling
// TODO Refactoring - http/net & gorm functions
func createUser(w http.ResponseWriter, r *http.Request) {
	var user User
	err := json.NewDecoder(r.Body).Decode(&user)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	// Encrypt password
	hash, hashErr := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)

	if hashErr != nil {
		log.Fatal(hashErr)
	}

	newUser := User{
		Username: user.Username,
		Name:     user.Name,
		Surname:  user.Surname,
		Password: string(hash),
		Email:    user.Email,
		Age:      user.Age,
		IsAdmin:  user.IsAdmin,
	}

	result := db.Create(&newUser)

	w.Header().Set("Content-Type", "application/json")
	if result.Error != nil {
		w.WriteHeader(http.StatusConflict)
		errResp := map[string]string{"message": "Already exist username"}
		jsonBody, _ := json.Marshal(errResp)
		w.Write(jsonBody)
	} else {
		w.WriteHeader(http.StatusCreated)
		respJson, _ := json.Marshal(newUser)
		w.Write(respJson)
	}
}

// GET Method
// Retrieves the user that given username
// Password does not retrieve
// @returns JSON object that includes user information if it is found
// @returns Error message in JSON object if it is not found
func getUser(w http.ResponseWriter, r *http.Request) {

	// Get request parameter from URL
	vars := mux.Vars(r)
	key := vars["username"]
	var user User

	// Fetch user data by comparing given username
	result := db.Where("username = ?", key).Select([]string{"id", "username", "name", "surname", "email", "age", "is_admin"}).First(&user)

	w.Header().Set("Content-Type", "application/json")
	if result.Error != nil {
		w.WriteHeader(http.StatusNotFound)
		res := map[string]string{"message": "User could not found"}
		resBody, _ := json.Marshal(res)
		w.Write(resBody)
	} else {
		w.WriteHeader(http.StatusOK)
		resBody, _ := json.Marshal(user)
		w.Write(resBody)
	}
}

// GET Method
// Retrieves all users
// @returns List of JSON objects which includes users
// @returns Error message if users cannot retrieved
func getUsers(w http.ResponseWriter, r *http.Request) {

	var users []User
	result := db.Select([]string{"id", "username", "name", "surname", "email", "age", "is_admin"}).Find(&users)

	w.Header().Set("Content-Type", "application/json")
	if result.Error != nil {
		w.WriteHeader(http.StatusNotFound)
		res := map[string]string{"message": "Users could not retrieved"}
		resBody, _ := json.Marshal(res)
		w.Write(resBody)
	} else {
		w.WriteHeader(http.StatusOK)
		resBody, _ := json.Marshal(users)
		w.Write(resBody)
	}
}

// POST Method
// Login the system as a user
// Retrieve the user and compare the password
// @returns Error message if username/email is not valid
// or password is not correct
// @returns User data as a JSON object
func login(w http.ResponseWriter, r *http.Request) {
	var user User
	var dbUser User
	err := json.NewDecoder(r.Body).Decode(&user)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	email := user.Email
	username := user.Username
	password := user.Password

	db.Where("username = ?", username).Or("email = ?", email).Find(&dbUser)

	w.Header().Set("Content-Type", "application/json")
	if len(dbUser.Name) == 0 {
		w.WriteHeader(http.StatusNotFound)
		var errMsg string
		if len(username) != 0 {
			errMsg = "Username could not found"
		} else {
			errMsg = "Email could not found"
		}
		res := map[string]string{"message": errMsg}
		resBody, _ := json.Marshal(res)
		w.Write(resBody)
		return
	}

	if hashErr := bcrypt.CompareHashAndPassword([]byte(dbUser.Password), []byte(password)); hashErr != nil {
		w.WriteHeader(http.StatusNotFound)
		res := map[string]string{"message": "Password is not correct"}
		resBody, _ := json.Marshal(res)
		w.Write(resBody)
	} else {
		w.WriteHeader(http.StatusOK)
		resBody, _ := json.Marshal(dbUser)
		w.Write(resBody)
	}
}

// PUT Method
// Updates the user profile
// @returns Error message if the user could not updated
// @returns Updated User data as a JSON object
func updateUserProfile(w http.ResponseWriter, r *http.Request)  {
	// Get username from URL
	vars := mux.Vars(r)
	username := vars["username"]

	var newUser User
	err := json.NewDecoder(r.Body).Decode(&newUser)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	// Find and assign data to *user
	var user User
	result := db.Where("username = ?", username).Find(&user)

	// Check username is exist or not
	if user.Username == "" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		res := map[string]string{"message": "Username could not find"}
		resBody, _ := json.Marshal(res)
		w.Write(resBody)
		return
	}

	// Update with new data
	user.Username = newUser.Username
	user.Name = newUser.Name
	user.Surname = newUser.Surname
	user.Email = newUser.Email
	user.Age = newUser.Age

	result = db.Save(&user)

	w.Header().Set("Content-Type", "application/json")
	if result.Error != nil {
		w.WriteHeader(http.StatusBadRequest)
		res := map[string]string{"message": "User profile could not updated"}
		resBody, _ := json.Marshal(res)
		w.Write(resBody)
	} else {
		w.WriteHeader(http.StatusOK)
		resBody, _ := json.Marshal(user)
		w.Write(resBody)
	}
}

// POST Method
// Meetup creation is handled here
// Only admins can create a new meetup
// @returns Error message if user has not admin authorization
// @returns Error message if JSON object could not decode
// @returns Error message if meetup could not save on db
// @returns JSON response after adding new meetup to db
func createMeetup(w http.ResponseWriter, r *http.Request)  {
	vars := mux.Vars(r)
	adminUsername := vars["admin_username"]

	// Check user is admin or not
	var user User
	result := db.Where("username = ? AND is_admin = ?",adminUsername, true).Find(&user)

	if user.Username == "" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		res := map[string]string{"message": "You are not allowed to create a new meetup. Only admins can create a new meetup"}
		errBody, _ :=json.Marshal(res)
		w.Write(errBody)
		return
	}

	var meetup Meetup
	err := json.NewDecoder(r.Body).Decode(&meetup)

	if err != nil {
		http.Error(w, "Error occurred while decoding JSON", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	result = db.Create(&meetup)

	w.Header().Set("Content-Type", "application/json")
	if result.Error != nil {
		w.WriteHeader(http.StatusBadRequest)
		errResp := map[string]string{"message": "Meetup could not created"}
		jsonBody, _ := json.Marshal(errResp)
		w.Write(jsonBody)
	} else {
		w.WriteHeader(http.StatusCreated)
		respBody, _ := json.Marshal(meetup)
		w.Write(respBody)
	}
}

func hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello Hit")
}

func handleRequests() {
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", hello)

	// User
	router.HandleFunc("/user/create", createUser).Methods("POST")
	router.HandleFunc("/user/{username}", getUser).Methods("GET")
	router.HandleFunc("/users", getUsers).Methods("GET")
	router.HandleFunc("/user/login", login).Methods("POST")
	router.HandleFunc("/user/update/{username}", updateUserProfile).Methods("PUT")

	// Meetup
	router.HandleFunc("/meetup/create/{admin_username}", createMeetup).Methods("POST")
	log.Fatal(http.ListenAndServe(":8081", router))
}
func main() {
	fmt.Println("MeetsGo API")

	USERNAME := os.Getenv("PG_USERNAME")
	PASSWORD := os.Getenv("PG_PASSWORD")

	dsn := "host=localhost user=" + USERNAME + " password=" + PASSWORD + " dbname=meetsup port=5432 sslmode=disable TimeZone=Europe/Istanbul"
	DB, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatal("Connection failed to open database")
	}

	log.Println("Connection established to database")

	db = DB

	db.AutoMigrate(&User{}, &Meetup{})

	handleRequests()
}
