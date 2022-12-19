package main

import (
	_ "encoding/json"
	"fmt"
	_ "io/ioutil"
	_ "log"
	_ "net/http"
	_ "strconv"
	"strings"
	_ "database/sql"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/gorilla/mux"
	"main.go/Users"
	"bufio"
	"os"
)


func main(){
outer:
	for {
		fmt.Println()
		fmt.Println("Ride sharing Platform\n",
				"1. Using app as a user\n",
				"2. Using app as a driver\n",
				"0. Leave application")
		fmt.Print("Enter an option: ")
		var choice string
		reader := bufio.NewReader(os.Stdin)
		input, _ := reader.ReadString('\n')
		choice = strings.TrimSpace(input)
		switch choice {
		case "1":
			Users.PassengerPage()
		case "2":
			Users.DriverPage()
		case "0":
			break outer
		default:
			fmt.Println("### Invalid Input ###")
		}
	}
}

