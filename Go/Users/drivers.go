package Users

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"strconv"
	"main.go/Trip"
	"net/http"
	"bytes"
	"io/ioutil"
	"encoding/json"
)

type Driver struct {
	DriverId      int    `json:"driverId"`
	FirstName     string `json:"firstName"`
	LastName      string `json:"lastName"`
	MobileNo  	  int 	 `json:"mobileNo"`
	EmailAddr     string `json:"emailAddr"`
	LicenseNo     string `json:"licenseNo"`
	IdNum 		  string `json:"idNum"`
	IsAvail       bool   `json:"isAvail"`
	Password 	  string `json:"password"`
}

//Creating the session state variable
var driverInfo Driver

//main page for drivers to view
func DriverPage() {
	for {
		fmt.Println()
		fmt.Println("Ride sharing Platform(Driver)\n",
				"1. Creating an account\n",
				"2. Login into account\n",
				"0. Back")
		fmt.Print("Enter an option: ")
		var choice string
		reader := bufio.NewReader(os.Stdin)
		input, _ := reader.ReadString('\n')
		choice = strings.TrimSpace(input)
		switch choice {
		case "1":
			CreateDriverAccount()
		case "2": 
			LoginToDriverAccount()
		case "0":
			os.Exit(0);
		default:
			fmt.Println("### Invalid Input ###")
			fmt.Println()
		}
	}
}

func LoggedInDriverPage() {
	fmt.Println()
		for {
			fmt.Println()
			fmt.Println("Ride sharing Platform(Logged In)\n",
					"1. Update account\n",
					"2. Trip Status\n",
					"0. Log out")
			fmt.Print("Enter an option: ")
			var choice string
			reader := bufio.NewReader(os.Stdin)
			input, _ := reader.ReadString('\n')
			choice = strings.TrimSpace(input)
			switch choice {
			case "1":
				UpdateDriverAccount()
			case "2": 
				Trip.TripStatus(driverInfo.DriverId, driverInfo.IsAvail)
			case "0":
				driverInfo = Driver{}
				DriverPage()
			default:
				fmt.Println("### Invalid Input ###")
			}
		}
	}


func LoginToDriverAccount(){
		var emailAddr string
		var password string
		for {
			for {
				fmt.Println()
				fmt.Println("======= Logging into account =======")
				fmt.Print("Enter your email address: ")
				reader := bufio.NewReader(os.Stdin)
				input, _ := reader.ReadString('\n')
				emailAddr = strings.TrimSpace(input)

				fmt.Print("Enter your password: ")
				reader = bufio.NewReader(os.Stdin)
				input, _ = reader.ReadString('\n')
				password = strings.TrimSpace(input)

				//Checks for empty field
				if emailAddr == "" || password == "" {
					fmt.Println("Please do not leave any blanks")
					continue
				} else {
					break
				}
			}

		//Creating a client
		client := &http.Client{}
		var driverData Driver
		//Calling login API
		url := "http://localhost:5000/driver/login?emailAddr=" + emailAddr + "&&password=" + password
		if req, err := http.NewRequest("GET", url, nil); err == nil {
			if res, err := client.Do(req); err == nil {
				if body, err2 := ioutil.ReadAll(res.Body); err2 == nil {
					//upon success
					if res.StatusCode == 202 {
						json.Unmarshal(body, &driverData)
						driverInfo = driverData
						fmt.Println("Welcome back "+ driverData.FirstName)
						LoggedInDriverPage() // Send to logged in page
					} else if res.StatusCode == 400 { //upon failure
						fmt.Println("There is no avaliable drivers")
					} else {
						fmt.Println("There is no avaliable drivers")
					}
				}
			}
		}
	}
}

func CreateDriverAccount() {
	var newDriver Driver
	fmt.Print("Enter your first name: ")
	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	newDriver.FirstName = strings.TrimSpace(input)

	fmt.Print("Enter your last name: ")
	reader = bufio.NewReader(os.Stdin)
	input, _ = reader.ReadString('\n')
	newDriver.LastName = strings.TrimSpace(input)

	fmt.Print("Enter your Mobile number: ")
	reader = bufio.NewReader(os.Stdin)
	input, _ = reader.ReadString('\n')
	store := strings.TrimSpace(input) // string
	//Conver the value to int for storing
	newDriver.MobileNo, _ = strconv.Atoi(store)
	
	fmt.Print("Enter your email address: ")
	reader = bufio.NewReader(os.Stdin)
	input, _ = reader.ReadString('\n')
	newDriver.EmailAddr = strings.TrimSpace(input)

	fmt.Print("Enter your license number: ")
	reader = bufio.NewReader(os.Stdin)
	input, _ = reader.ReadString('\n')
	newDriver.LicenseNo = strings.TrimSpace(input)

	fmt.Print("Enter your Identification number: ")
	reader = bufio.NewReader(os.Stdin)
	input, _ = reader.ReadString('\n')
	newDriver.IdNum = strings.TrimSpace(input)

	//Ensuring that the passwords match
	for {
		fmt.Print("Enter a Password: ")
		reader = bufio.NewReader(os.Stdin)
		input, _ = reader.ReadString('\n')
		newDriver.Password = strings.TrimSpace(input)
		
		fmt.Print("Re-enter your password: ")
		reader = bufio.NewReader(os.Stdin)
		input, _ = reader.ReadString('\n')
		rePassword := strings.TrimSpace(input)

		//Break if passwords match 
		if rePassword == newDriver.Password {
			break;
		} else {
			//print error message
			fmt.Println("Passwords do not match please try again!")
			continue;
		}
	}

	client := &http.Client{}
	url := "http://localhost:5000/driver"
	postBody, _ := json.Marshal(newDriver)
	resBody := bytes.NewBuffer(postBody)
	if req, err := http.NewRequest(http.MethodPost, url, resBody); err == nil {
		if res, err2 := client.Do(req); err2 == nil {
			if res.StatusCode == 202 {
				fmt.Println("Successfully created account!")
			} else if res.StatusCode == 400 {
				fmt.Println("Error - Bad Request")
			}
		}
	}
}

func UpdateDriverAccount() {
	var newDriver Driver
	fmt.Println("First Name: ",driverInfo.FirstName)
	fmt.Print("Enter your first name: ")
	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	newDriver.FirstName = strings.TrimSpace(input)

	fmt.Println("Last name: ",driverInfo.LastName)
	fmt.Print("Enter your last name: ")
	reader = bufio.NewReader(os.Stdin)
	input, _ = reader.ReadString('\n')
	newDriver.LastName = strings.TrimSpace(input)

	fmt.Println("Mobile Number: ",driverInfo.MobileNo)
	fmt.Print("Enter your Mobile number: ")
	reader = bufio.NewReader(os.Stdin)
	input, _ = reader.ReadString('\n')
	store := strings.TrimSpace(input) // string
	//Conver the value to int for storing
	newDriver.MobileNo, _ = strconv.Atoi(store)
	
	fmt.Println("Email address: ",driverInfo.EmailAddr)	
	fmt.Print("Enter your email address: ")
	reader = bufio.NewReader(os.Stdin)
	input, _ = reader.ReadString('\n')
	newDriver.EmailAddr = strings.TrimSpace(input)

	fmt.Println("License number: ",driverInfo.LicenseNo)	
	fmt.Print("Enter your license number: ")
	reader = bufio.NewReader(os.Stdin)
	input, _ = reader.ReadString('\n')
	newDriver.LicenseNo = strings.TrimSpace(input)

	//Setting the driver idmumber as cannot change
	newDriver.IdNum = driverInfo.IdNum
	//Setting the drivers id
	newDriver.DriverId = driverInfo.DriverId
	//Ensuring that the passwords match
	for {
		fmt.Print("Enter a Password: ")
		reader = bufio.NewReader(os.Stdin)
		input, _ = reader.ReadString('\n')
		newDriver.Password = strings.TrimSpace(input)
		
		fmt.Print("Re-enter your password: ")
		reader = bufio.NewReader(os.Stdin)
		input, _ = reader.ReadString('\n')
		rePassword := strings.TrimSpace(input)

		//Break if passwords match 
		if rePassword == newDriver.Password {
			break;
		} else {
			//print error message
			fmt.Println("Passwords do not match please try again!")
			continue;
		}
	}

	client := &http.Client{}
	url := "http://localhost:5000/driver"
	postBody, _ := json.Marshal(newDriver)
	resBody := bytes.NewBuffer(postBody)
	if req, err := http.NewRequest(http.MethodPatch, url, resBody); err == nil {
		if res, err2 := client.Do(req); err2 == nil {
			if res.StatusCode == 202 {
				driverInfo = newDriver
				fmt.Println("Successfully Updated account!")
			} else if res.StatusCode == 400 {
				fmt.Println("Error - Bad Request")
			}
		}
	}
}