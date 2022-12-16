package main

//
// Name/ID : Hugo Von Louwen Dorosan/S10202923
//
// =================[Features]===================
// General:
// 1) Passenger Login x
// 2) Driver Login x
// 3) Exit
//
// ============[Programming Check List]===========
// Pssenger:
// 1) Update Passenger Details x
// 2) View History x
// 3) Book Ride - input pick up point and drop off location + auto assign driver x
//
// Driver:
// 1) Update Driver Details (Except identification number) x
// 2) Option to choose which passenger to pick up (optional)
// 3) Option to end the ride x
//
//
// ============[Debugging Check List]===========
// Pssenger:
// 1) Update Passenger Details
// 2) View History x
// 3) Book Ride - input pick up point and drop off location + auto assign driver x
//
// Driver:
// 1) Update Driver Details (Except identification number)
// 2) Option to choose which passenger to pick up (optional)
// 3) Option to end the ride x
// 4) View History
// Ride:
// 1)
//
//

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"bufio"
)

//==================[Object Structs for Passanger, Driver and Ride]==================

type Passenger struct {
	FirstName string `json:"First Name"`
	LastName  string `json:"Last Name"`
	MobileNum int    `json:"Mobile Number"`
	EmailAddr string `json:"Email Address"`
}

type Driver struct {
	FirstName     string `json:"First Name"`
	LastName      string `json:"Last Name"`
	MobileNum     int    `json:"Mobile Number"`
	EmailAddr     string `json:"Email Address"`
	DriverStatus  string `json:"Driver Status"`
	LicenceNumber string `json:"Licence Number"`
}

type Ride struct {
	PassengerName     string `json:"Passenger"`
	Driver            string `json:"Driver"`
	TripStatus        string `json:"Trip Status"`
	StartDateTime     string `json:"Pick-Up Timing"`
	EndDateTime       string `json:"Drop-Off Timing"`
	PostalCodePickUp  string `json:"Pick up point"`
	PostalCodeDropOff string `json:"Drop off point"`
}

//==================[Variables for the list of availablerides, user's ID (can be either passenger or driver) and passOnRide boolean to check if passenger is currently on a ride]==================

var availableRidesList = map[string]map[string]Ride{}
var userID string
var passOnRide bool

// Menu Prompts
var menuitems = []string{"Passenger", "Driver", "Exit"}
var passMenuItems = []string{"Book Ride", "View Ride History", "Exit"}

func main() {

	for {
		// Printing out the Main Meny
		fmt.Println("============")
		for i := 0; i < len(menuitems); i++ {
			fmt.Println(i+1, ". ", menuitems[i])
		}

		// Prompting the user which type of account to log in as
		fmt.Println("Which user are you intending to log in as?: ")
		reader := bufio.NewReader(os.Stdin)
		userInput, _ := reader.ReadString('\n')
		userInput = strings.TrimSpace(userInput)
		fmt.Println("===================================================")

		// Logging in as a Passenger User
		if userInput == "1" {

			// Ask to input email
			fmt.Println("Enter email: ")
			reader := bufio.NewReader(os.Stdin)
			emailInput, _ := reader.ReadString('\n')
			emailInput = strings.TrimSpace(emailInput)

			// Passenger list to store passengers
			// validated variable to see whether email input is in the database
			var passengerList = map[string]map[string]Passenger{}
			var validated bool
			var passenger Passenger

			// Calling Passenger Microservice to get all passangers and store them in a list
			resp, err := http.Get("http://localhost:6969/api/v1/passenger")
			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				log.Fatalln(err)
			}
			json.Unmarshal([]byte(body), &passengerList)

			// Checks whether input email exists
			for key, element := range passengerList["Passenger"] {
				// fmt.Print(element.FirstName, " ( ", key, " ) ", "\nFirst Name ", element.FirstName, "\nLast Name ", element.LastName, "\nMobile Number ", element.MobileNum, "\nEmail Address ", element.EmailAddr, "\n\n")
				if emailInput == element.EmailAddr {
					// returns true if it is in
					validated = true
					userID = key
					passenger = element
				}
			}
			if validated != false {
				// Move on to the Prompt Passenger Menu if passenger email exists
				PromptPassenger(passenger)
			} else {
				// Validation error message if input email does not exist
				fmt.Println("Passenger " + emailInput + " does not exist")
			}

			// Logging in as a Driver User
		} else if userInput == "2" {

			// Asking user to input email
			fmt.Println("Enter email: ")
			reader := bufio.NewReader(os.Stdin)
			emailInput, _ := reader.ReadString('\n')
			emailInput = strings.TrimSpace(emailInput)

			// Driver list to store drivers
			// validated variable to see whether email input is in the database
			var driverList = map[string]map[string]Driver{}
			var validated bool
			var driverStruct Driver

			resp, err := http.Get("http://localhost:1000/api/v1/driver")
			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				log.Fatalln(err)
			}
			json.Unmarshal([]byte(body), &driverList)

			// Check if email is in the database
			for key, element := range driverList["Driver"] {
				if emailInput == element.EmailAddr {
					// returns true if it is in
					validated = true
					// get the ID of the rider which will be useful in updating driver status later on
					driverStruct = element
					// save driverID
					userID = key
				}
			}

			if validated != false {
				PromptDriver(driverStruct)
			} else {
				fmt.Println(emailInput + " does not exist!")
			}

		} else if userInput == "3" {
			fmt.Println("Closing")
			break
		}
	}

}

func PromptDriver(d Driver) {

	LoadAvailableRides()
	inprogmenuitems := []string{"End Ride", "Exit"}
	availablemenuitems := []string{"Update Account Details", "Exit"}
	fmt.Println("\n===================================================")
	fmt.Println("Hello Driver " + d.FirstName + " " + d.LastName + "!\nWelcome to Ride Sharing Platform!\n")
	var finished bool

	for {
		if d.DriverStatus == "Busy" {
			// Grabbing the current Ride
			var ridesList = map[string]map[string]Ride{}
			var currentRide Ride
			var currentRideKey string

			resp, err := http.Get("http://localhost:5000/api/v1/ride?tripstatus=inprogress")
			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				log.Fatalln(err)
			}

			json.Unmarshal([]byte(body), &ridesList)

			for a, b := range ridesList["Ride"] {
				fmt.Println(a + " " + b.Driver)
				if b.Driver == userID {
					currentRide = b
					currentRideKey = a
				}
			}

			fmt.Println("You are currently driving for Passenger " + currentRide.PassengerName + "\n")
			fmt.Println("===================[Ride Details]==================")
			fmt.Println("Pick Up Point: " + currentRide.PostalCodePickUp)
			fmt.Println("Drop Off Point: " + currentRide.PostalCodeDropOff)
			fmt.Println("===================================================")

			if finished != true {
				fmt.Println("===================================================")
				for i := 0; i < len(inprogmenuitems); i++ {
					fmt.Println(i+1, ". ", inprogmenuitems[i])
				}

				fmt.Println("Select Option: ")
				reader := bufio.NewReader(os.Stdin)
				driverInput, _ := reader.ReadString('\n')
				driverInput = strings.TrimSpace(driverInput)
				fmt.Println("===================================================")

				if driverInput == "1" {

					fmt.Println("Are you sure you want to end this ride? (Type YES to confirm)")
					reader := bufio.NewReader(os.Stdin)
					cfmEndInput, _ := reader.ReadString('\n')
					cfmEndInput = strings.TrimSpace(cfmEndInput)
					fmt.Println("===================================================")
					if cfmEndInput == "YES" {
						UpdateDriverStatus("Available", d, userID)
						UpdateRide("Finished", currentRide, currentRideKey)
					} else {
						fmt.Println("Unrecognizble input, please input YES to end the ride")
					}
					finished = true

				} else if driverInput == "2" {
					fmt.Println("Closing Driver Menu...")
					finished = true
				}
				break
			}

		} else if d.DriverStatus == "Available" {

			fmt.Println("===================================================")
			for i := 0; i < len(availablemenuitems); i++ {
				fmt.Println(i+1, ". ", availablemenuitems[i])
			}

			fmt.Println("Select Option: ")
			reader := bufio.NewReader(os.Stdin)
			driverInput, _ := reader.ReadString('\n')
			driverInput = strings.TrimSpace(driverInput)
			fmt.Println("===================================================")

			if driverInput == "1" {
				UpdateDriver(d)
			} else if driverInput == "2" {
				break
			}

		}
	}

}

func PromptPassenger(p Passenger) {

	menuitems := []string{"Book Ride", "View Ride History", "Update Account Details", "Exit"}

	checkIfPassengerOnRide()

	for {
		fmt.Println("===================================================")
		fmt.Println("Hello " + p.FirstName + " " + p.LastName + "!\nWelcome to Ride Sharing Platform!\n")

		for i := 0; i < len(menuitems); i++ {
			fmt.Println(i+1, ". ", menuitems[i])
		}

		reader := bufio.NewReader(os.Stdin)
		userInput, _ := reader.ReadString('\n')
		userInput = strings.TrimSpace(userInput)
		fmt.Println("===================================================")

		if userInput == "1" {

			if passOnRide != true {

				var puInput string
				var doInput string

				fmt.Println("Key in Postal Code of Pick Up Location (Postal Code): ")
				reader := bufio.NewReader(os.Stdin)
				puInput, _ = reader.ReadString('\n')
				puInput = strings.TrimSpace(puInput)

				if _, err := strconv.Atoi(puInput); err == nil {
					fmt.Printf("\nPostal Code " + puInput + " has been selected\n")

					fmt.Println("\nKey in Postal Code of Drop Off Location (Postal Code): ")
					reader = bufio.NewReader(os.Stdin)
					doInput, _ = reader.ReadString('\n')
					doInput = strings.TrimSpace(doInput)

					if _, err := strconv.Atoi(doInput); err == nil {
						fmt.Printf("\nPostal Code " + doInput + " has been selected\n")
					} else {
						fmt.Println("\nImproper Drop Off Location Input")
					}

				} else {
					fmt.Println("\nImproper Pick Up Location Input")
				}

				CreateRide(puInput, doInput)

			} else if passOnRide == true {
				fmt.Println("You are currently on a ride! Wait for your ride to end or inform your driver to cancel this ride")
			} else {
				fmt.Println("Improper Postal Code Input")
			}

		} else if userInput == "2" {
			ViewRideHistory("Passenger")
		} else if userInput == "3" {
			UpdatePassenger(p)
		} else if userInput == "4" {
			break
		}

	}

}

func ViewRideHistory(userType string) {

	if userType == "Passenger" {

		var ridesList = map[string]map[string]Ride{}
		resp, err := http.Get("http://localhost:6969/api/v1/passenger?history=" + userID)
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Fatalln(err)
		}

		json.Unmarshal([]byte(body), &ridesList)

		for key, element := range ridesList["Ride"] {
			if element.TripStatus == "Finished" {
				fmt.Println("=====================[", key, "]===================\n",
					"\nDriver ID: ", element.Driver, "\nTrip Status: ", element.TripStatus, "\nTrip Start Time: ", element.StartDateTime, "\nTrip End Time: ", element.EndDateTime, "\nPick Up Location: ", element.PostalCodePickUp, "\nDrop Off Location ", element.PostalCodeDropOff)
				fmt.Println("Start Date: " + element.StartDateTime)
			} else if element.TripStatus == "In Progress" {
				fmt.Println("=====================[", key, "]===================\n",
					"\nDriver ID: ", element.Driver, "\nTrip Status: ", element.TripStatus, "\nTrip Start Time: ", element.StartDateTime, "\nPick Up Location: ", element.PostalCodePickUp, "\nDrop Off Location ", element.PostalCodeDropOff)
			} else if len(ridesList) == 0 {
				fmt.Println("Empty Riding History")
			}
		}
		fmt.Println("===================================================")
	} else if userType == "Driver" {
		var allrideList = map[string]map[string]Ride{}
		var userideList = map[string]Ride{}

		resp, err := http.Get("http://localhost:1000/api/v1/Driver/History")
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Fatalln(err)
		}
		json.Unmarshal([]byte(body), &allrideList)

		for key, element := range allrideList["Passenger"] {
			fmt.Println(" ( ", key, " ) ", "\nFirst Name ", element.PassengerName, "\nLast Name ", element.Driver, "\nMobile Number ", element.TripStatus, "\nEmail Address ", element.StartDateTime, "\nEmail Address ", element.EndDateTime, "\nEmail Address ", element.PostalCodePickUp, "\nEmail Address ", element.PostalCodeDropOff)
		}

		for a, b := range allrideList["Ride"] {
			if b.Driver != userID {
				userideList[a] = b
			}
		}
	}
}

func UpdatePassenger(pass Passenger) {

	var updatedFname string
	var updatedLname string
	var updatedMobileNum int
	var updatedEmailAddr string

	// Prompting user for the new updates
	fmt.Println("==============[Changing Passenger Details Prompt (type x to skip)]==============")
	fmt.Println("Key in new First Name:")
	reader := bufio.NewReader(os.Stdin)
	updatedFname, _ = reader.ReadString('\n')
	updatedFname = strings.TrimSpace(updatedFname)
	if updatedFname == "x" {
		updatedFname = pass.FirstName
	}
	fmt.Println("===================================================")

	fmt.Println("Key in new Last Name:")
	reader = bufio.NewReader(os.Stdin)
	updatedLname, _ = reader.ReadString('\n')
	updatedLname = strings.TrimSpace(updatedLname)
	if updatedLname == "x" {
		updatedLname = pass.LastName
	}
	fmt.Println("===================================================")

	fmt.Println("Key in new Mobile Number:")
	reader = bufio.NewReader(os.Stdin)
	mn, _ := reader.ReadString('\n')
	mn = strings.TrimSpace(mn)
	if mn == "x" {
		updatedMobileNum = pass.MobileNum
	} else {
		updatedMobileNum, _ = strconv.Atoi(mn)
	}
	fmt.Println("===================================================")

	fmt.Println("Key in new Email Address:")
	reader = bufio.NewReader(os.Stdin)
	updatedEmailAddr, _ = reader.ReadString('\n')
	updatedEmailAddr = strings.TrimSpace(updatedEmailAddr)
	if updatedEmailAddr == "x" {
		updatedEmailAddr = pass.EmailAddr
	}
	fmt.Println("===================================================")

	// UPDATE RIDE TO DATABASE
	var updatedPassenger Passenger
	updatedPassenger.FirstName = updatedFname
	updatedPassenger.LastName = updatedLname
	updatedPassenger.MobileNum = updatedMobileNum
	updatedPassenger.EmailAddr = updatedEmailAddr
	jsonBody, _ := json.Marshal(updatedPassenger)

	client := &http.Client{}
	if req, err := http.NewRequest("PUT", "http://localhost:6969/api/v1/passenger/"+userID, bytes.NewBuffer(jsonBody)); err == nil {
		if _, err := client.Do(req); err == nil {
			fmt.Print("You have updated your details successfully!\n")
		}
	}

}

func UpdateDriver(d Driver) {
	var updatedFname string
	var updatedLname string
	var updatedMobileNum int
	var updatedEmailAddr string
	var updatedLN string

	fmt.Println("==============[Changing Driver Details Prompt (type x to skip)]==============")
	// Prompting user for the new updates
	fmt.Println("Key in new First Name:")
	reader := bufio.NewReader(os.Stdin)
	updatedFname, _ = reader.ReadString('\n')
	updatedFname = strings.TrimSpace(updatedFname)
	if updatedFname == "x" {
		updatedFname = d.FirstName
	}
	fmt.Println("===================================================")

	fmt.Println("Key in new Last Name:")
	reader = bufio.NewReader(os.Stdin)
	updatedLname, _ = reader.ReadString('\n')
	updatedLname = strings.TrimSpace(updatedLname)
	if updatedLname == "x" {
		updatedLname = d.LastName
	}
	fmt.Println("===================================================")

	fmt.Println("Key in new Mobile Number:")
	reader = bufio.NewReader(os.Stdin)
	mn, _ := reader.ReadString('\n')
	mn = strings.TrimSpace(mn)
	if mn == "x" {
		updatedMobileNum = d.MobileNum
	} else {
		updatedMobileNum, _ = strconv.Atoi(mn)
	}
	fmt.Println("===================================================")

	fmt.Println("Key in new Email Address:")
	reader = bufio.NewReader(os.Stdin)
	updatedEmailAddr, _ = reader.ReadString('\n')
	updatedEmailAddr = strings.TrimSpace(updatedEmailAddr)
	if updatedEmailAddr == "x" {
		updatedEmailAddr = d.EmailAddr
	}
	fmt.Println("===================================================")

	fmt.Println("Key in new Licence Number:")
	reader = bufio.NewReader(os.Stdin)
	updatedLN, _ = reader.ReadString('\n')
	updatedLN = strings.TrimSpace(updatedLN)
	if mn == "x" {
		updatedLN = d.LicenceNumber
	}
	fmt.Println("===================================================")

	// UPDATE RIDE TO DATABASE
	var updatedDriver Driver
	updatedDriver.FirstName = updatedFname
	updatedDriver.LastName = updatedLname
	updatedDriver.MobileNum = updatedMobileNum
	updatedDriver.EmailAddr = updatedEmailAddr
	updatedDriver.LicenceNumber = updatedLN
	updatedDriver.DriverStatus = d.DriverStatus
	jsonBody, _ := json.Marshal(updatedDriver)

	client := &http.Client{}
	if req, err := http.NewRequest("PUT", "http://localhost:1000/api/v1/driver/"+userID, bytes.NewBuffer(jsonBody)); err == nil {
		if _, err := client.Do(req); err == nil {
			fmt.Print("Your account has been updated successfully!\n")
		}
	}

}

func UpdateDriverStatus(status string, driverStruct Driver, driverId string) {

	var updatedDriver Driver

	updatedDriver = driverStruct
	updatedDriver.DriverStatus = status
	jsonBody, _ := json.Marshal(updatedDriver)

	client := &http.Client{}
	if req, err := http.NewRequest("PUT", "http://localhost:1000/api/v1/driver/"+driverId, bytes.NewBuffer(jsonBody)); err == nil {
		if _, err := client.Do(req); err == nil {
			fmt.Println("")
		}
	}

}

func LoadAvailableRides() {

	resp, err := http.Get("http://localhost:5000/api/v1/Ride/Open")
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}

	// Add all
	json.Unmarshal([]byte(body), &availableRidesList)
	// Using dummyID for now -> change to randomized ID later
	for key, element := range availableRidesList["Ride"] {
		fmt.Print(" ( ", key, " ) ", "\nPassenger Name ", element.PassengerName, "\nDriver Name ", element.Driver, "\nPostal Code Pick Up Point ", element.PostalCodePickUp, "\nPostal Code Drop Off Point ", element.PostalCodeDropOff, "\nTrip Status ", element.TripStatus, "\nPick up timing ", element.StartDateTime, "\nDrop off timing ", element.EndDateTime, "\n\n")
	}

}

func CreateRide(pcPickUp string, pcDropOff string) {

	currentTime := time.Now()
	rideId := "R" + strconv.Itoa(time.Now().Year()) + strconv.Itoa(int(time.Now().Month())) + strconv.Itoa(time.Now().Day()) + strconv.Itoa(time.Now().Day()) + strconv.Itoa(time.Now().Hour()) + strconv.Itoa(time.Now().Minute()) + strconv.Itoa(time.Now().Second()) + userID

	var newRide Ride

	newRide.PassengerName = userID
	newRide.TripStatus = "In Progress"
	newRide.StartDateTime = currentTime.Format("2006-01-02 15:04:05")
	newRide.EndDateTime = time.Time{}.Format("2006-01-02 15:04:05")
	newRide.PostalCodePickUp = pcPickUp
	newRide.PostalCodeDropOff = pcDropOff

	// Checking if there are available drivers

	var obtaineddriverList = map[string]map[string]Driver{}
	var availabledriverList = map[string]Driver{}

	resp, err := http.Get("http://localhost:1000/api/v1/driver?driverstatus=available")
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}
	json.Unmarshal([]byte(body), &obtaineddriverList)

	for a, b := range obtaineddriverList["Driver"] {
		availabledriverList[a] = b
	}

	if len(availabledriverList) > 0 {

		client := &http.Client{}
		jsonBody, _ := json.Marshal(newRide)
		if req, err := http.NewRequest("POST", "http://localhost:5000/api/v1/ride/"+rideId, bytes.NewBuffer(jsonBody)); err == nil {
			if _, err := client.Do(req); err == nil {
				fmt.Print("")
			}
		}

		allocateRideToDriver(rideId, newRide)
		checkIfPassengerOnRide()

	} else {
		fmt.Println("There are no available drivers right now. Please book again later")

	}

}

func allocateRideToDriver(rideId string, updatedRide Ride) {

	var obtaineddriverList = map[string]map[string]Driver{}
	var availabledriverList = map[string]Driver{}
	var driverIdList = []string{}

	resp, err := http.Get("http://localhost:1000/api/v1/driver?driverstatus=available")
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}
	json.Unmarshal([]byte(body), &obtaineddriverList)

	for a, b := range obtaineddriverList["Driver"] {
		availabledriverList[a] = b
		driverIdList = append(driverIdList, a)
	}

	a := driverIdList[0]

	UpdateDriverStatus("Busy", availabledriverList[a], a)
	updatedRide.Driver = a
	client := &http.Client{}
	jsonBody, _ := json.Marshal(updatedRide)
	if req, err := http.NewRequest("PUT", "http://localhost:5000/api/v1/ride/"+rideId, bytes.NewBuffer(jsonBody)); err == nil {
		if _, err := client.Do(req); err == nil {
			fmt.Print("===============[Ride Details]===============\nDriver ID: ", updatedRide.Driver, "\nRide ID: ", rideId, "!\nPick up point: ", updatedRide.PostalCodePickUp, "\nDrop off point: ", updatedRide.PostalCodeDropOff)
		}
	}

}

func UpdateRide(status string, riderStruct Ride, rideId string) {

	currentTime := time.Now()

	var updatedRide Ride
	updatedRide = riderStruct
	updatedRide.TripStatus = status
	updatedRide.EndDateTime = currentTime.Format("2006-01-02 15:04:05")
	jsonBody, _ := json.Marshal(updatedRide)

	client := &http.Client{}
	if req, err := http.NewRequest("PUT", "http://localhost:5000/api/v1/ride/"+rideId, bytes.NewBuffer(jsonBody)); err == nil {
		if _, err := client.Do(req); err == nil {
			fmt.Print("Driver ", rideId, " Updated Successfully!\n")
		}
	}

}

func checkIfPassengerOnRide() {

	// Grabbing the current Ride
	var ridesList = map[string]map[string]Ride{}

	resp, err := http.Get("http://localhost:5000/api/v1/ride?tripstatus=inprogress")
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}

	json.Unmarshal([]byte(body), &ridesList)

	passOnRide = false

	for _, b := range ridesList["Ride"] {
		if b.PassengerName == userID {
			passOnRide = true
		}
	}

}
