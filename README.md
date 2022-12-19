# ETI_Assignment1

Name: lincoln<br />
Class: P03<br />
ID: 10203100B<br />

## Contents
1. Requirements and design Considerations
2. Architecture diagrams
3. Instructions for setting up and running microservices

## Requirements 
![image](https://user-images.githubusercontent.com/73088199/208426491-ff3294b5-78bf-4dae-8e76-35a9213e2580.png)
1. **Passenger menu** Allows user to select the services that they would like to use
2. **Login Passenger** Ask the user to input in their email address as well as their passwords in order to login in to their account
3. **Create Passenger** Allow users to create their own account by inputing their first name , last name, mobile number and email address while having a password. Ids are auto assigned through auto increment.
4. **Quit menu** Allow users to leave the system or to go back one menu
5. **Update Passegner** Allows user to change their information including their password information 
6. **Display Trip History** Allows user to view their trip histories in reverse chronological order.
7. **Request trip** Allows user to input in their current location and destination location and checks if there is an avaliable driver to start the trip.
8. **Login Driver** Allow driver to login using their email address and their password
9. **Create Driver** Allow new driver to create an account with information such as first name, last name, email address, mobile number, identification number and car license number. 
10. **Update driver** Allow drivers to change information on their account except for their identification number. 
11. **Trip status** This tells the driver if they currently have a ride 
12. **Check for a ride** The system will check if the ride is currently avaliable and then procced to show them starting or ending a ride.
13. **Start a ride** This allows the driver to start a ride with a yes no trigger and would reflect on the database 
14. **End a ride** This allows the driver to end a ride witha yes or no trigger which would end the ride for the user.

## Design Consideration 
**Type of checking** 
To ensure that a user is able to have both a driver and user account, code allows for user to only user to only select either passenger or driver menu

**Checking of trips for driver**
For riders to track their trips i have a isAvail property in their database. This is used to check if the driver currently has a ride for assigning rides to the trips. Assigning of rides is done automatically.

**Checking of start and stoping rides** 
For starting and stopping rides the database has a column which contians isstarted and isCompleted this is to ensure that the driver has started a ride and is ready to end it. Once a ride is completed the driver will be reset to available.


## Architecture diagrams 
![image](https://user-images.githubusercontent.com/73088199/208431452-169584aa-522b-4200-9e01-23eb20ed3b86.png)
Passenger and Driver connects to user api which connects passenger and driver database. this includes functions such as logining in, creating account. However, they also interact with the trip api in order to start and maintain the rides which is connected to the trip api. 


## Class Diagram
![image](https://user-images.githubusercontent.com/73088199/208432853-b223f1fc-8ca2-4d90-840e-f2f1ce2fab5c.png)
There are 3 entities which are passenger driver and trip. Passenger and driver is able to have 1 trip at a time. However, the passenger and driver can have multiple trips in the past.

## Instructions for setting up and running microservices
## Start trip api
1. cd into tripAPI folder
2. run command "go run tripApi.go"

## Start user api
1. cd into userAPI folder
2. run command "go run userApi.go"

## Start Cli 
1. cd into go file
2. go run main.go

## SQL script for setting up databases
1. Driver Database
CREATE TABLE `driver` (
  `driverId` int NOT NULL AUTO_INCREMENT,
  `firstName` varchar(45) NOT NULL,
  `lastName` varchar(45) NOT NULL,
  `emailAddr` varchar(45) NOT NULL,
  `mobileNo` int NOT NULL,
  `identificationNo` varchar(45) NOT NULL,
  `licenseNo` varchar(45) NOT NULL,
  `isAvail` tinyint NOT NULL DEFAULT '1',
  `Password` varchar(45) DEFAULT NULL,
  PRIMARY KEY (`driverId`)
)

2. Passenger Database
CREATE TABLE `passenger` (
  `ID` int NOT NULL AUTO_INCREMENT,
  `FirstName` varchar(30) NOT NULL,
  `LastName` varchar(30) NOT NULL,
  `MobileNo` int NOT NULL,
  `EmailAddr` varchar(45) NOT NULL,
  `Password` varchar(45) NOT NULL,
  PRIMARY KEY (`ID`)
)

3. Trip Database
CREATE TABLE `trip` (
  `tripId` int NOT NULL AUTO_INCREMENT,
  `startPostal` int NOT NULL,
  `endPostal` int NOT NULL,
  `driverId` int NOT NULL,
  `passengerId` int NOT NULL,
  `isCompleted` tinyint NOT NULL DEFAULT '0',
  `isStarted` tinyint NOT NULL DEFAULT '0',
  `createdAt` datetime NOT NULL,
  PRIMARY KEY (`tripId`)
)

