package Users

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"net/http"
	_ "database/sql"
	"strconv"
	"main.go/Trip"
	"encoding/json"
	"bytes"
	"io/ioutil"
)

type Passenger struct {
	ID  		 int    `json:"id"`
	FirstName    string `json:"firstName"`
	LastName     string `json:"lastName"`
	MobileNo	 int    `json:"mobileNo"`
	EmailAddr    string `json:"emailAddr"`
	Password     string `json:"password"`
}

//Decalring session state variable 
var passengerInfo Passenger

//Main page for the passengers to view
func PassengerPage() {
		for {
			fmt.Println()
			fmt.Println("Ride sharing Platform(passenger)\n",
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
				CreatePassengerAccount()
			case "2":
				LoginToPassgengerAccount()
			case "0":
				os.Exit(0);
			default:
				fmt.Println("### Invalid Input ###")
			}
		}
	}

func LoggedInPassengerPage() {
	for {
		fmt.Println()
		fmt.Println("Ride sharing Platform(Logged In)\n",
				"1. Update account\n",
				"2. Book a trip\n",
				"3. Trip History\n",
				"0. Log out",)
		fmt.Print("Enter an option: ")
		var choice string
		reader := bufio.NewReader(os.Stdin)
		input, _ := reader.ReadString('\n')
		choice = strings.TrimSpace(input)
		switch choice {
		case "1":
			//caling updating of passenger account function
			UpdatePassengerAccount()
		case "2": 
			//Calling booking trip function 
			Trip.BookTrip(passengerInfo.ID)
		case "3":
			//Getting users trip history 
			Trip.TripHistory(passengerInfo.ID)
		case "0":
			//Clearing session state variable
			passengerInfo = Passenger{}
			PassengerPage()
		default:
			fmt.Println("### Invalid Input ###")
			fmt.Println()
		}
	}
}

func CreatePassengerAccount() {
	var newPassenger Passenger

	//Creation of form to collect user information
	fmt.Print("Enter your first name: ")
	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	newPassenger.FirstName = strings.TrimSpace(input)

	fmt.Print("Enter your last name: ")
	reader = bufio.NewReader(os.Stdin)
	input, _ = reader.ReadString('\n')
	newPassenger.LastName = strings.TrimSpace(input)

	//Ensuring that Mobile number is 8 digit value
	for {
		fmt.Print("Enter your Mobile number: ")
		reader = bufio.NewReader(os.Stdin)
		input, _ = reader.ReadString('\n')
		//Check if the input has 8 digits
		if len(strings.TrimSpace(input)) == 8 {
			//Checks if the input is only digits
			store, err := strconv.Atoi(strings.TrimSpace(input)) 
			if err != nil{
				fmt.Println("The input contains non-integer characters!")
				continue;
			} else { //Break the value and store the value
				newPassenger.MobileNo = store
				break;
			}
		} else {
			fmt.Println("Please have 8 digits")
			continue;
		}
	}
	
	fmt.Print("Enter your email address : ")
	reader = bufio.NewReader(os.Stdin)
	input, _ = reader.ReadString('\n')
	newPassenger.EmailAddr = strings.TrimSpace(input)


	//Ensuring that the passwords match
	for {
		fmt.Print("Enter a Password: ")
		reader = bufio.NewReader(os.Stdin)
		input, _ = reader.ReadString('\n')
		newPassenger.Password = strings.TrimSpace(input)
		
		fmt.Print("Re-enter your password: ")
		reader = bufio.NewReader(os.Stdin)
		input, _ = reader.ReadString('\n')
		rePassword := strings.TrimSpace(input)

		//Break if passwords match 
		if rePassword == newPassenger.Password {
			break;
		} else {
			//print error message
			fmt.Println("Passwords do not match please try again!")
			continue;
		}
	}

	client := &http.Client{}
	url := "http://localhost:5000/passenger"
	postBody, _ := json.Marshal(newPassenger)
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

//Updating / changing user account details 
func UpdatePassengerAccount() {
	var updateUser Passenger

	//Creation of form to collect user 
	fmt.Println("First Name: ", passengerInfo.FirstName)
	fmt.Print("Enter your first name: ")
	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	updateUser.FirstName = strings.TrimSpace(input)

	fmt.Println("Last name: ", passengerInfo.LastName)
	fmt.Print("Enter your last name: ")
	reader = bufio.NewReader(os.Stdin)
	input, _ = reader.ReadString('\n')
	updateUser.LastName = strings.TrimSpace(input)

	fmt.Println("Mobile Number: ", passengerInfo.MobileNo)
	//Ensuring that Mobile number is 8 digit value
	for {
		fmt.Print("Enter your Mobile number: ")
		reader = bufio.NewReader(os.Stdin)
		input, _ = reader.ReadString('\n')
		//Check if the input has 8 digits
		if len(strings.TrimSpace(input)) == 8 {
			//Checks if the input is only digits
			store, err := strconv.Atoi(strings.TrimSpace(input)) 
			if err != nil{
				fmt.Println("The input contains non-integer characters!")
				continue;
			} else { //Break and store the value
				updateUser.MobileNo = store
				break;
			}
		} else {
			fmt.Println("Please have 8 digits")
			continue;
		}
	}
	fmt.Println("Email address: ", passengerInfo.EmailAddr)
	fmt.Print("Enter your email address: ")
		reader = bufio.NewReader(os.Stdin)
		input, _ = reader.ReadString('\n')
		updateUser.EmailAddr = strings.TrimSpace(input)
	//Ensuring that the passwords match
	for {
		fmt.Print("Enter a Password: ")
		reader = bufio.NewReader(os.Stdin)
		input, _ = reader.ReadString('\n')
		updateUser.Password = strings.TrimSpace(input)
		
		fmt.Print("Re-enter your password: ")
		reader = bufio.NewReader(os.Stdin)
		input, _ = reader.ReadString('\n')
		rePassword := strings.TrimSpace(input)

		//Break if passwords match 
		if rePassword == updateUser.Password {
			break;
		} else {
			//print error message
			fmt.Println("Passwords do not match please try again!")
			continue;
		}
	}

	updateUser.ID = passengerInfo.ID
	client := &http.Client{}
	url := "http://localhost:5000/passenger"
	patchBody, _ := json.Marshal(updateUser)
	resBody := bytes.NewBuffer(patchBody)
	if req, err := http.NewRequest("PATCH", url, resBody); err == nil {
		if res, err2 := client.Do(req); err2 == nil {
			if res.StatusCode == 202 {
				//updating global variable 
				passengerInfo = updateUser
				fmt.Println("Successfully updated account")
			} else if res.StatusCode == 400 {
				fmt.Println("Error - Bad Request")
			}
		}
	}
	
}

func LoginToPassgengerAccount(){
	var emailAddr string
	var password string
	for {
		//Ensuring that fields are not null
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
		var passengerData Passenger
		//Calling login API
		url := "http://localhost:5000/passenger/login?emailAddr=" + emailAddr + "&&password=" + password
		if req, err := http.NewRequest("GET", url, nil); err == nil {
			if res, err := client.Do(req); err == nil {
				if body, err2 := ioutil.ReadAll(res.Body); err2 == nil {
					//upon success
					if res.StatusCode == 202 {
						json.Unmarshal(body, &passengerData)
						passengerInfo = passengerData
						fmt.Println("Welcome back "+ passengerData.FirstName)
						LoggedInPassengerPage() // Send to logged in page
					} else if res.StatusCode == 400 { //upon failure
						fmt.Println("Email or password is wrong")
					} else {
						fmt.Println("Email or password is wrong")
					}
				}
			}
		}
	}
}