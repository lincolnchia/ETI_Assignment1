package Trip

import (
		"bufio"
	"fmt"
	"os"
	"strings"
	"strconv"
	"net/http"
	"bytes"
	"io/ioutil"
	"encoding/json"
	"time"
)

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

func BookTrip(passengerId int) {
	//storing driver name
	var driverName string
	var driverData Driver
	// Create A trip for a passenger
	var newTrip Trip
	//Pick up point
	for {
		fmt.Println()
		fmt.Println("====== Booking Trip ======")
		fmt.Print("Enter postal code of pick up point: ")
		reader := bufio.NewReader(os.Stdin)
		input, _ := reader.ReadString('\n')
		//Check if the input has 5 or 6 digits
		if len(strings.TrimSpace(input)) == 5 ||  len(strings.TrimSpace(input)) == 6{
			//Checks if the input is only digits
			store, err := strconv.Atoi(strings.TrimSpace(input)) 
			if err != nil{
				fmt.Println("The input contains non-integer characters!")
				continue;
			} else { //Break the value and store the value
				newTrip.StartPostal = store
				break;
			}
		} else {
			fmt.Println("Please have 5-6 digits")
			continue;
		}
	}
	//Drop off point
	for {
		fmt.Print("Enter postal code of drop off point: ")
		reader := bufio.NewReader(os.Stdin)
		input, _ := reader.ReadString('\n')
		//Check if the input has 5 or 6 digits
		if len(strings.TrimSpace(input)) == 5 ||  len(strings.TrimSpace(input)) == 6{
			//Checks if the input is only digits
			store, err := strconv.Atoi(strings.TrimSpace(input)) 
			if err != nil{
				fmt.Println("The input contains non-integer characters!")
				continue;
			} else { 
				//Check if the start postal and end postal are the same
				if store == newTrip.StartPostal {
					fmt.Println("You have input in the same location")
					continue;
				} else {
					newTrip.EndPostal = store
					break;
				}
			}
		} else {
			fmt.Println("Please have 5-6 digits")
			continue;
		}
	}
	client := &http.Client{}
		url := "http://localhost:5001/trip?passengerId="+ strconv.Itoa(passengerId)
		postBody, _ := json.Marshal(newTrip)
		resBody := bytes.NewBuffer(postBody)
		if req, err := http.NewRequest(http.MethodPost, url, resBody); err == nil {
			if res, err2 := client.Do(req); err2 == nil {
				if body, err2 := ioutil.ReadAll(res.Body); err2 == nil {
					if res.StatusCode == 202 {
						json.Unmarshal(body, &driverData)
						driverName = driverData.FirstName +" "+driverData.LastName
						fmt.Println("Your trip has been created your driver is "+ driverName)
					} else if res.StatusCode == 400 {
						fmt.Println("No avaliable drivers")
					}
				}
			}
		}
}

func TripHistory(passengerId int) {
	//Creating new trip array
	tripMap := make(map[int]Trip)
	client := &http.Client{}
		url := "http://localhost:5001/trip?passengerId="+ strconv.Itoa(passengerId)
		postBody, _ := json.Marshal(tripMap)
		resBody := bytes.NewBuffer(postBody)
		if req, err := http.NewRequest(http.MethodGet, url, resBody); err == nil {
			if res, err2 := client.Do(req); err2 == nil {
				if body, err2 := ioutil.ReadAll(res.Body); err2 == nil {
					if res.StatusCode == 202 {
						json.Unmarshal(body, &tripMap)
					} else if res.StatusCode == 400 {
						fmt.Println("Error - Bad Request")
					}
				}
			}
		}
	//Number of trips
	count := 0 
	fmt.Println("======= Trip Details =======")
	for _ , element := range tripMap {
		count += 1
		//creating a new trip
		var tripDetails Trip
		//storing map element into tripDetails
		tripDetails = element
		dateTimeString := tripDetails.CreatedAt.Format("2006-01-02 15:04:05")
		fmt.Println(count,". Start Postal: " + strconv.Itoa(tripDetails.StartPostal) + " End Postal: " + strconv.Itoa(tripDetails.StartPostal)  + " Date:" , dateTimeString)
	}
}
func TripStatus(driverId int, isAvail bool) {
	var tripDetails Trip
	//Checking if the driver is driving 
	if isAvail == true {
		fmt.Println("You currently do not have any rides")
	} else {
		client := &http.Client{}
		url := "http://localhost:5001/trip/tripStatus?driverId="+ strconv.Itoa(driverId)
		postBody, _ := json.Marshal(driverId)
		resBody := bytes.NewBuffer(postBody)
		if req, err := http.NewRequest(http.MethodGet, url, resBody); err == nil {
			if res, err2 := client.Do(req); err2 == nil {
				if body, err2 := ioutil.ReadAll(res.Body); err2 == nil {
					if res.StatusCode == 202 {
						json.Unmarshal(body, &tripDetails)
					} else if res.StatusCode == 400 {
						fmt.Println("No avaliable drivers")
					}
				}
			}
		}
		//Checks if the ride has started
		if tripDetails.IsCompleted == false && tripDetails.IsStarted == false{
			fmt.Print("You currently have a trip!\nStart Ride Y/N:  ")
			var choice string
			reader := bufio.NewReader(os.Stdin)
			input, _ := reader.ReadString('\n')
			choice = strings.ToUpper(strings.TrimSpace(input))
			switch choice {
			case "Y":
				startTrip(tripDetails.TripId)
				break;
			case "N":
				break;
			default:
				fmt.Println("### Invalid Input ###")
			}
		}
		//Checks if the ride is currently on going
		if tripDetails.IsCompleted == false && tripDetails.IsStarted != false{
			fmt.Print("You are currently on a trip!\nEnd Ride Y/N:  ")
			var choice string
			reader := bufio.NewReader(os.Stdin)
			input, _ := reader.ReadString('\n')
			choice = strings.ToUpper(strings.TrimSpace(input))
			switch choice {
			case "Y":
				endTrip(tripDetails.TripId)
				break;
			case "N":
				break;
			default:
				fmt.Println("### Invalid Input ###")
			}
		}
		
	}
}
func startTrip(tripId int){
	client := &http.Client{}
		url := "http://localhost:5001/trip/tripStatus?tripId="+ strconv.Itoa(tripId) +"&&isStarted=0"
		postBody, _ := json.Marshal(tripId)
		resBody := bytes.NewBuffer(postBody)
		if req, err := http.NewRequest(http.MethodPatch, url, resBody); err == nil {
			if res, err2 := client.Do(req); err2 == nil {
				if res.StatusCode == 202 {
					fmt.Println("Trip has started!")
				} else if res.StatusCode == 400 {
					fmt.Println("Problems with SQL")
				}
			}
		}
}

func endTrip(tripId int){
	client := &http.Client{}
		url := "http://localhost:5001/trip/tripStatus?tripId="+ strconv.Itoa(tripId) +"&&isStarted=1"
		postBody, _ := json.Marshal(tripId)
		resBody := bytes.NewBuffer(postBody)
		if req, err := http.NewRequest(http.MethodPatch, url, resBody); err == nil {
			if res, err2 := client.Do(req); err2 == nil {
				if res.StatusCode == 202 {
					fmt.Println("Trip has Ended!")
				} else if res.StatusCode == 400 {
					fmt.Println("Problems with SQL")
				}
			}
		}
}