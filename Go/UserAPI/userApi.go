package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"io/ioutil"
	// "strings"
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

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/passenger", passengerEndpoint).Methods("POST", "PATCH")
	router.HandleFunc("/passenger/login", authPassengerEndpoint).Methods("GET")
	router.HandleFunc("/driver", driverEndpoint).Methods("POST", "PATCH")
	router.HandleFunc("/driver/login", authDriverEndpoint).Methods("GET")
	fmt.Println("Listening at port 5000")
	log.Fatal(http.ListenAndServe(":5000", router))
}
func passengerEndpoint(w http.ResponseWriter, r *http.Request) {
	// params := mux.Vars(r)
	if r.Method =="POST"{
		if body, err := ioutil.ReadAll(r.Body); err == nil {
			var newPassenger Passenger
			if err := json.Unmarshal(body, &newPassenger); err == nil{
				//Opening database connection
				db, err := sql.Open("mysql","root:root@tcp(127.0.0.1:3306)/my_db")
				// handle error upon failure
				if err != nil {
					http.Error(w, err.Error(), http.StatusBadRequest)
				}
				defer db.Close()
				
				//inserting values into passenger table
				_, err = db.Exec("insert into passenger (firstName, lastName, mobileNo,emailAddr, password) values(?, ?, ?, ?, ?)",
				newPassenger.FirstName, newPassenger.LastName, newPassenger.MobileNo, newPassenger.EmailAddr, newPassenger.Password)
				//Handling error of SQL statement
				if err != nil {
					http.Error(w, err.Error(), http.StatusBadRequest)
				}
				w.WriteHeader(http.StatusAccepted)
			}
		}
	} else if r.Method == "PATCH" {
		//Digest Passenger object from Body
		var updateUser Passenger
		err := json.NewDecoder(r.Body).Decode(&updateUser)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
		if updateUser.ID == 0 { //ID Does not have a value means missing id
			http.Error(w, "PassengerId missing", http.StatusBadRequest)
			
		}
		//Opening database connection
		db, err := sql.Open("mysql","root:root@tcp(127.0.0.1:3306)/my_db")
		// handle error upon failure
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
		defer db.Close()
		//Update Passenger in DB
		stmt, err := db.Prepare("UPDATE passenger SET firstName = ?, lastName = ?, mobileNo = ? ,emailAddr = ?, password = ? WHERE ID = ?")
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
		defer stmt.Close()
		_, err = stmt.Exec(updateUser.FirstName, updateUser.LastName, updateUser.MobileNo, updateUser.EmailAddr, updateUser.Password, updateUser.ID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
		} else {
			w.WriteHeader(http.StatusAccepted)
		}
	} else {
		http.Error(w, "Bad Request", http.StatusBadRequest)
	}
}
func authPassengerEndpoint(w http.ResponseWriter, r *http.Request){
	if r.Method == "GET"{
		querystringmap := r.URL.Query()
		emailAddr := querystringmap.Get("emailAddr")
		password := querystringmap.Get("password")
		//Opening database connection
		db, err := sql.Open("mysql","root:root@tcp(127.0.0.1:3306)/my_db")
		// handle error upon failure
		if err != nil {
			http.Error(w, "Unable to connect", http.StatusBadRequest)
		}
		defer db.Close()

		//inserting values into passenger table
		result, err := db.Query("select * from passenger where EmailAddr = '" + emailAddr + "'"+ "and Password = '" + password +"'")
		//Handling error of SQL statement
		if err != nil {
			http.Error(w, "Missing data", http.StatusBadRequest)
			panic(err.Error())
		}
		var newPassenger Passenger
		for result.Next() {
		err = result.Scan( &newPassenger.ID, &newPassenger.FirstName, &newPassenger.LastName, &newPassenger.MobileNo, &newPassenger.EmailAddr, &newPassenger.Password )
			if err != nil {
				http.Error(w, "Missing data", http.StatusBadRequest)
			} else {
				//Jsonfiying the data to send back to UI 
				output, _ := json.Marshal(newPassenger)
				w.WriteHeader(http.StatusAccepted)
				fmt.Fprintf(w, string(output))
			}
		}
	} else {
		http.Error(w, "Bad Request", http.StatusBadRequest)
	}
}

func driverEndpoint(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST"{
		if body, err := ioutil.ReadAll(r.Body); err == nil {
			var newDriver Driver

			if err := json.Unmarshal(body, &newDriver); err == nil{
				//Opening database connection
				db, err := sql.Open("mysql","root:root@tcp(127.0.0.1:3306)/my_db")
				// handle error upon failure
				if err != nil {
					http.Error(w, err.Error(), http.StatusBadRequest)
				}
				defer db.Close()
				
				//Creating SQL statement to push data into driver
				_, err = db.Exec("insert into driver (firstName, lastName, emailAddr, mobileNo, identificationNo, licenseNo, password) values(?, ?, ?, ?, ?, ?, ?)",
				newDriver.FirstName, newDriver.LastName, newDriver.EmailAddr, newDriver.MobileNo, newDriver.LicenseNo, newDriver.IdNum, newDriver.Password)
				if err != nil {
					http.Error(w, err.Error(), http.StatusBadRequest)
				} else {
					w.WriteHeader(http.StatusAccepted)
				}
			}
		}
	} else if r.Method == "PATCH"{
		//Digest Passenger object from Body
		var updateDriver Driver
		err := json.NewDecoder(r.Body).Decode(&updateDriver)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
		if updateDriver.DriverId == 0 { //ID Does not have a value means missing id
			http.Error(w, "PassengerId missing", http.StatusBadRequest)
			
		}
		//Opening database connection
		db, err := sql.Open("mysql","root:root@tcp(127.0.0.1:3306)/my_db")
		// handle error upon failure
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
		defer db.Close()
		//Update Passenger in DB
		stmt, err := db.Prepare("UPDATE driver SET firstName = ?, lastName = ?, emailAddr = ? ,mobileNo = ?,identificationNo = ?,licenseNo = ?, password = ? WHERE driverId = ?")
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
		defer stmt.Close()
		_, err = stmt.Exec(updateDriver.FirstName, updateDriver.LastName, updateDriver.EmailAddr, updateDriver.MobileNo,updateDriver.IdNum, updateDriver.LicenseNo, updateDriver.Password, updateDriver.DriverId)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
		} else {
			w.WriteHeader(http.StatusAccepted)
		}
	} else {
		http.Error(w, "Bad Request", http.StatusBadRequest)
	}
}

func authDriverEndpoint(w http.ResponseWriter, r *http.Request){
	if r.Method == "GET"{
		querystringmap := r.URL.Query()
		emailAddr := querystringmap.Get("emailAddr")
		password := querystringmap.Get("password")
		//Opening database connection
		db, err := sql.Open("mysql","root:root@tcp(127.0.0.1:3306)/my_db")
		// handle error upon failure
		if err != nil {
			http.Error(w, "Unable to connect", http.StatusBadRequest)
		}
		defer db.Close()

		//inserting values into passenger table
		result, err := db.Query("select * from driver where EmailAddr = '" + emailAddr + "'"+ "and Password = '" + password +"'")
		//Handling error of SQL statement
		if err != nil {
			http.Error(w, "Missing data", http.StatusBadRequest)
		}
		var newDriver Driver
		for result.Next() {
			err = result.Scan( &newDriver.DriverId, &newDriver.FirstName, &newDriver.LastName, &newDriver.EmailAddr, &newDriver.MobileNo, &newDriver.LicenseNo, &newDriver.IdNum, &newDriver.IsAvail, &newDriver.Password )
			if err != nil {
				fmt.Print(err.Error())
				http.Error(w, "Missing data", http.StatusBadRequest)
			} else {
				//Jsonfiying the data to send back to UI 
				output, _ := json.Marshal(newDriver)
				w.WriteHeader(http.StatusAccepted)
				fmt.Fprintf(w, string(output))
			}
		}
	} else {
		http.Error(w, "Bad Request", http.StatusBadRequest)
	}
}
