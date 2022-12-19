package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	_ "io/ioutil"
	"time"
	// "strings"
	"strconv"
)

type Passenger struct {
	ID  		 int    `json:"id"`
	FirstName    string `json:"firstName"`
	LastName     string `json:"lastName"`
	MobileNo	 int    `json:"mobileNo"`
	EmailAddr    string `json:"emailAddr"`
	Password     string `json:"password"`
}

type Driver struct {
	DriverId      int    `json:"driverId"`
	FirstName     string `json:"firstName"`
	LastName      string `json:"lastName"`
	MobileNo  	  int 	`json:"mobileNo"`
	EmailAddr     string `json:"emailAddr"`
	LicenseNo     string `json:"licenseNo"`
	IdNum 		  string `json:"idNum"`
	IsAvail       bool   `json:"isAvail"`
	Password 	  string `json:"password"`
}

type Trip struct {
	TripId 		int		 `json:"tripId"`
	StartPostal int		 `json:"startPostal"`
	EndPostal 	int		 `json:"endPostal"`
	DriverId 	int		 `json:"driverId"`
	PassengerId int		 `json:"passengerId"`
	IsCompleted bool 	 `json:"isCompleted"`
	IsStarted 	bool	 `json:"isStarted"`
	CreatedAt 	time.Time `json:"createdAt"`

}

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/trip", tripEndpoint).Methods("GET", "POST")
	router.HandleFunc("/trip/tripStatus", driverTripEndpoint).Methods("GET", "PATCH")
	fmt.Println("Listening at port 5001")
	log.Fatal(http.ListenAndServe(":5001", router))
}

func tripEndpoint(w http.ResponseWriter, r *http.Request) {
	
	if r.Method == "POST"{
		//Getting variables from query string
		querystringmap := r.URL.Query()
		passengerId := querystringmap.Get("passengerId")
		
		//New trip variable
		var newTrip Trip
		//New Driver variable
		var driver Driver

		err := json.NewDecoder(r.Body).Decode(&newTrip)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}

		//Opening database connection
		db, err := sql.Open("mysql","root:root@tcp(127.0.0.1:3306)/my_db")
		// handle error upon failure
		if err != nil {
			http.Error(w, "Unable to connect", http.StatusBadRequest)
		}
		defer db.Close()

		//Checking for driver who is avaliable
		result, err := db.Query("select * from driver where isAvail = 1")
			//Handling error of SQL statement
			if err != nil {
				http.Error(w, "Unable to connect", http.StatusBadRequest)
			}
			// Checking to see if there is a return variable
			if result == nil {
				fmt.Println("There is no avaliable drivers")
			} else {
				for result.Next() {
					err = result.Scan( &driver.DriverId, &driver.FirstName, &driver.LastName, &driver.EmailAddr, &driver.MobileNo, &driver.IdNum, &driver.LicenseNo, &driver.IsAvail, &driver.Password )
					if err != nil {
					http.Error(w, "Unable to connect", http.StatusBadRequest)
					}
					newTrip.DriverId = driver.DriverId
				}
				//Checking if the returned driver contains values 
				if newTrip.DriverId != 0 {
					//Getting the time that trip is created
					timeNow := time.Now()
					dateTimeString := timeNow.Format("2006-01-02 15:04:05")
					//Creating a trip 
					tripResult, err := db.Exec("insert into trip (startPostal, endPostal, driverId, passengerId, createdAt) values(?, ?, ?, ?, ?)",
					newTrip.StartPostal, newTrip.EndPostal, newTrip.DriverId, passengerId, dateTimeString)
					//Handling error of SQL statement
					if err != nil {
						fmt.Println("Error with sending data to database")
						http.Error(w, "Unable to connect", http.StatusBadRequest)
					}
					//Checking to see if the trip is created
					id, err := tripResult.LastInsertId()
					if err != nil || id == 0{
						fmt.Println("Your Trip creation has failed, please try again!")
						http.Error(w, "Unable to connect", http.StatusBadRequest)
					}
					//Update the value of the driver database
					intDriverId := strconv.Itoa(newTrip.DriverId)
					_, err = db.Exec("UPDATE driver SET isAvail = 0 WHERE driverId = " + intDriverId)
					//Handling error of SQL statement
					if err != nil {
						fmt.Println("Error with sending data to database")
						http.Error(w, "Unable to connect", http.StatusBadRequest)
					} else {
						//sending back driver information
						output, _ := json.Marshal(driver)
						w.WriteHeader(http.StatusAccepted)
						fmt.Fprintf(w, string(output))
					}

				} else{ // When there is no avaliable drivers
					http.Error(w, "Missing data", http.StatusBadRequest)
				}
			}
	} else if r.Method == "GET" {
		//Getting variables from query string
		querystringmap := r.URL.Query()
		passengerId := querystringmap.Get("passengerId")
		count := 0

		//Creating new trip details
		var tripDetails Trip
		tripMap := make(map[int]Trip)

		//Opening database connection
		db, err := sql.Open("mysql","root:root@tcp(127.0.0.1:3306)/my_db?parseTime=true")
		// handle error upon failure
		if err != nil {
			http.Error(w, "Unable to connect", http.StatusBadRequest)
		}
		defer db.Close()

		result, err := db.Query("select * from Trip where PassengerId = " + passengerId + " ORDER BY createdAt DESC")
		//Handling error of SQL statement
		if err != nil {
			fmt.Println("Error with getting data from database")
			http.Error(w, "Unable to connect", http.StatusBadRequest)
		}
		//declare new passenger list
		for result.Next() {
			count++
			err = result.Scan( &tripDetails.TripId, &tripDetails.StartPostal, &tripDetails.EndPostal, &tripDetails.DriverId, &tripDetails.PassengerId, &tripDetails.IsCompleted, &tripDetails.IsStarted, &tripDetails.CreatedAt )
			if err != nil {
				http.Error(w, "Unable to connect", http.StatusBadRequest)
			} else {
				//Add to the list
				tripMap[count] = tripDetails
			}
		}
		output, _ := json.Marshal(tripMap)
		w.WriteHeader(http.StatusAccepted)
		fmt.Fprintf(w, string(output))

	} else {
		fmt.Print("hi")
		http.Error(w, "Bad Request", http.StatusBadRequest)
	}

	
}
func driverTripEndpoint(w http.ResponseWriter, r *http.Request){
	if r.Method == "GET"{
		querystringmap := r.URL.Query()
		driverId := querystringmap.Get("driverId")

		//Opening database connection
		db, err := sql.Open("mysql","root:root@tcp(127.0.0.1:3306)/my_db?parseTime=true")
		// handle error upon failure
		if err != nil {
			fmt.Println("Error connecting to database")
			panic(err.Error())
		}
		defer db.Close()
		//Check if the driver is currently avaliable
		result, err := db.Query("select * from trip where driverId = "+ driverId +" and isCompleted = 0")
		//Handling error of SQL statement
		if err != nil {
			fmt.Println("Error with getting data from database")
			panic(err.Error())
		} 
		for result.Next(){
			var tripDetails Trip
			err = result.Scan( &tripDetails.TripId, &tripDetails.StartPostal, &tripDetails.EndPostal, &tripDetails.DriverId, &tripDetails.PassengerId, &tripDetails.IsCompleted, &tripDetails.IsStarted, &tripDetails.CreatedAt )
			if err != nil {
				panic(err.Error())
			} else {
				output, _ := json.Marshal(tripDetails)
				w.WriteHeader(http.StatusAccepted)
				fmt.Fprintf(w, string(output))
			}
		}

	}else if r.Method == "PATCH"{
		//getting tripid 
		querystringmap := r.URL.Query()
		tripId := querystringmap.Get("tripId")
		//Checking if api call is for is started or for is completed
		isStarted := querystringmap.Get("isStarted")
		//Opening database connection
		db, err := sql.Open("mysql","root:root@tcp(127.0.0.1:3306)/my_db")
		// handle error upon failure
		if err != nil {
			fmt.Println("Error connecting to database")
			panic(err.Error())
		}
		defer db.Close()
		if (isStarted == "0"){
			//Update the trip to start the ride
			_, err = db.Query("UPDATE trip SET isStarted = 1 WHERE tripId = " + tripId)
			//handle error upon failure
			if err != nil {
				http.Error(w, "Problem with SQL", http.StatusBadRequest)
			} else {
				w.WriteHeader(http.StatusAccepted)
			}
		} else if (isStarted == "1") {
			_, err = db.Query("UPDATE trip SET isCompleted = 1 WHERE tripId ="+ tripId)
			//handle error upon failure
			if err != nil {
				http.Error(w, "Problem with SQL", http.StatusBadRequest)
			} else {
				w.WriteHeader(http.StatusAccepted)
			}
		}
	} else {
		http.Error(w, "Bad Request", http.StatusBadRequest)
	}
}