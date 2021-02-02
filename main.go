package main

import (
	"encoding/json"
	//"bytes"
	//"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var db *gorm.DB

type User struct {
	ID			uint	`gorm:"primaryKey"`
	Username	string	`gorm:"unique"`
	Name 		string
	Surname		string
	Password 	string
	Email		string
	Age 		int
	IsAdmin		bool
}

type Meetup struct {
	ID 			uint	`gorm:"primaryKey;autoIncrement:true"`
	MeetupName	string
	MeetupDetails	string
	StartDate	time.Time
	EndDate		time.Time
	Address 	string
	Quota		int
	RegisteredUserCount	int
}

// POST Method
// Creates and saves a new user on the database
// TODO Error handling
// TODO Refactoring - http/net & gorm functions
// TODO Encrypt password before storing data
func createUser(w http.ResponseWriter, r *http.Request) {
	var user User
	err := json.NewDecoder(r.Body).Decode(&user)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	newUser := User{
		Username: user.Username,
		Name: user.Name,
		Surname: user.Surname,
		Password: user.Password,
		Email: user.Email,
		Age: user.Age,
		IsAdmin: user.IsAdmin,
	}

	result := db.Create(&newUser)

	if result.Error != nil {
		fmt.Println(result.Error)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusConflict)
		errResp := map[string]string{"message": "Already exist username"}
		jsonBody, _ := json.Marshal(errResp)
		w.Write(jsonBody)
	} else {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		respJson, _ := json.Marshal(newUser)
		w.Write(respJson)
	}
}

func hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello Hit")
}

func handleRequests() {
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", hello)
	router.HandleFunc("/user/create", createUser).Methods("POST")
	log.Fatal(http.ListenAndServe(":8081", router))
}
func main() {
	fmt.Println("MeetsGo API")

	USERNAME := os.Getenv("PG_USERNAME")
	PASSWORD := os.Getenv("PG_PASSWORD")

	dsn := "host=localhost user=" + USERNAME + " password=" + PASSWORD +" dbname=meetsup port=5432 sslmode=disable TimeZone=Europe/Istanbul"
	DB, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatal("Connection failed to open database")
	}

	log.Println("Connection established to database")

	db = DB

	db.AutoMigrate(&User{}, &Meetup{})

	handleRequests()
}