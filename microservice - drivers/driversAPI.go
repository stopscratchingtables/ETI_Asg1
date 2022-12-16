package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

var x Driver
var url = "http:/localhost:1000/api/v1/driver"

var driverList map[string]Driver = map[string]Driver{}

// STRUCTS

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
	PassengerName     string         `json:"Passenger"`
	Driver            string         `json:"Driver"`
	TripStatus        string         `json:"Trip Status"`
	StartDateTime     sql.NullString `json:"Pick-Up Timing"`
	EndDateTime       sql.NullString `json:"Drop-Off Timing"`
	PostalCodePickUp  int            `json:"Pick up point"`
	PostalCodeDropOff int            `json:"Drop off point"`
}

// FUNCTIONS

func main() {

	router := mux.NewRouter()
	router.HandleFunc("/api/v1/driver", driverFilter)
	router.HandleFunc("/api/v1/driver/{driverID}", driver).Methods("GET", "DELETE", "PUT", "POST")
	fmt.Println("Listening at port 1000")
	log.Fatal(http.ListenAndServe(":1000", router))

}

func allDriver(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-type", "application/json")
	db, err := sql.Open("mysql", "user:password@tcp(127.0.0.1:3306)/my_db")

	// handle error
	if err != nil {
		panic(err.Error())
	}

	defer db.Close()

	GetDB(db)

	querystringmap := r.URL.Query()
	fmt.Fprintf(w, "List All Drivers\n\n")
	for k, v := range querystringmap {
		fmt.Fprintf(w, "%s = %v\n", k, v[0])
	}

}

func driver(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-type", "application/json")

	params := mux.Vars(r)
	driverVal := params["driverID"]

	db, err := sql.Open("mysql", "user:password@tcp(127.0.0.1:3306)/my_db")
	// handle error
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	GetDB(db)
	e, err := json.Marshal(driverList)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(e))

	if details, exists := driverList[driverVal]; exists {
		if r.Method == "GET" { // GET METHOD

			data, _ := json.Marshal(details)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Fprintf(w, "%s\n", data)

		} else if r.Method == "DELETE" {

			_, err := db.Exec("DELETE FROM PASSENGER WHERE ID=?", driverVal)
			if err != nil {
				panic(err.Error())
			}
			fmt.Fprintf(w, driverVal, "Deleted")

		} else if r.Method == "POST" {

			w.WriteHeader(http.StatusConflict)

		} else if r.Method == "PUT" {
			updatedDriver := Driver{}
			reqBody, _ := ioutil.ReadAll(r.Body)
			json.Unmarshal(reqBody, &updatedDriver)

			_, err := db.Exec("UPDATE Driver SET FirstName=?, LastName=?, MobileNum=?, EmailAddress=?, DriverStatus=?, LicenceNumber=? Where ID=?", updatedDriver.FirstName, updatedDriver.LastName, updatedDriver.MobileNum, updatedDriver.EmailAddr, updatedDriver.DriverStatus, updatedDriver.LicenceNumber, driverVal)
			fmt.Println("Successfull")
			if err != nil {
				panic(err.Error())
			}
			w.WriteHeader(http.StatusAccepted)

		} else {
			fmt.Fprintf(w, "Invalid Passenger ID")
		}
	} else {

		if r.Method == "POST" {

			newDriver := Driver{}
			reqBody, _ := ioutil.ReadAll(r.Body)
			json.Unmarshal(reqBody, &newDriver)
			_, err := db.Exec("INSERT into Driver (ID, FirstName, LastName, MobileNum, EmailAddress, DriverStatus, LicenceNumber) value(?, ?, ?, ?, ?, ?)", driverVal, newDriver.FirstName, newDriver.LastName, newDriver.MobileNum, newDriver.EmailAddr, newDriver.DriverStatus, newDriver.LicenceNumber)
			if err != nil {
				panic(err.Error())
			}
			driverList[driverVal] = newDriver
			w.WriteHeader(http.StatusAccepted)

		} else if r.Method == "PUT" {
			// error message
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprintf(w, "Passenger ID Does Not Exist!")

		}

	}

}

// If there is no query, list all drivers.
// If there is query, list "driverFilter"

func driverFilter(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-type", "application/json")

	db, err := sql.Open("mysql", "user:password@tcp(127.0.0.1:3306)/my_db")
	// handle error
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()
	GetDB(db)

	query := r.URL.Query()
	results := map[string]Driver{}

	if value := query.Get("q"); len(value) > 0 {
		for k, v := range driverList { // e.g. k="IT", v="{"Diplpma in IT, ......, ...., ....."}"
			if strings.Contains(strings.ToLower(v.FirstName), strings.ToLower(value)) {
				results[k] = v
			}
		}
		if len(results) == 0 {
			w.WriteHeader(http.StatusNotFound)
		} else {

			// prints results
			data, err := json.Marshal(map[string]map[string]Driver{"Driver": results})
			if err != nil {
				log.Fatal(err)
			}
			if len(results) != 0 {
				fmt.Fprintf(w, "%s\n", data)
			}
		}
	} else if value := r.URL.Query().Get("driverstatus"); len(value) > 0 {

		var availabledriverList map[string]Driver = map[string]Driver{}

		// CONTINUE WITH GETTING AVAIL DRIVERS

		if value == "available" {

			results, err := db.Query("SELECT * FROM Driver WHERE DriverStatus = 'Available'; ")
			if err != nil {
				panic(err.Error())
			}
			for results.Next() {
				var driverID string
				var scannedDriver Driver
				err = results.Scan(&driverID, &scannedDriver.FirstName, &scannedDriver.LastName, &scannedDriver.EmailAddr, &scannedDriver.MobileNum, &scannedDriver.DriverStatus, &scannedDriver.LicenceNumber)
				if err != nil {
					panic(err.Error())
				}
				fmt.Println(scannedDriver.FirstName, scannedDriver.LastName, scannedDriver.MobileNum, scannedDriver.EmailAddr, scannedDriver.DriverStatus, scannedDriver.LicenceNumber)
				availabledriverList[driverID] = scannedDriver
			}

			data, err := json.Marshal(map[string]map[string]Driver{"Driver": availabledriverList}) // prints "Driver" + json of the drivers
			if err != nil {
				log.Fatal(err)
			}
			if len(driverList) != 0 {
				fmt.Fprintf(w, "%s\n", data)
			}
		}

	} else if value := query.Get("History"); len(value) > 0 {

		var rideHistory = map[string]Ride{}

		db, err := sql.Open("mysql", "user:password@tcp(127.0.0.1:3306)/my_db")
		if err != nil {
			panic(err.Error())
		}
		defer db.Close()

		results, err := db.Query(`SELECT * FROM Ride ORDER BY StartDateTime DESC`)
		if err != nil {
			panic(err.Error())
		}
		for results.Next() {
			var rideID string
			var scannedRide Ride
			err = results.Scan(&rideID, &scannedRide.PassengerName, &scannedRide.Driver, &scannedRide.TripStatus, &scannedRide.StartDateTime, &scannedRide.EndDateTime, &scannedRide.PostalCodePickUp, &scannedRide.PostalCodeDropOff)
			if err != nil {
				panic(err.Error())
			}
			fmt.Println(&scannedRide.PassengerName, scannedRide.Driver, scannedRide.TripStatus, scannedRide.StartDateTime, scannedRide.EndDateTime, scannedRide.PostalCodePickUp, scannedRide.PostalCodeDropOff)
			rideHistory[rideID] = scannedRide
		}

		data, _ := json.Marshal(rideHistory)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Fprintf(w, "%s\n", data)

	} else {

		data, err := json.Marshal(map[string]map[string]Driver{"Driver": driverList}) // prints "Driver" + json of the drivers
		if err != nil {
			log.Fatal(err)
		}
		if len(driverList) != 0 {
			fmt.Fprintf(w, "%s\n", data)
		}

	}

}

func GetDB(db *sql.DB) {
	results, err := db.Query("SELECT * FROM Driver")
	if err != nil {
		panic(err.Error())
	}
	for results.Next() {
		var driverID string
		var scannedDriver Driver
		err = results.Scan(&driverID, &scannedDriver.FirstName, &scannedDriver.LastName, &scannedDriver.EmailAddr, &scannedDriver.MobileNum, &scannedDriver.DriverStatus, &scannedDriver.LicenceNumber)
		if err != nil {
			panic(err.Error())
		}
		fmt.Println(driverID, scannedDriver.FirstName, scannedDriver.LastName, scannedDriver.EmailAddr, scannedDriver.DriverStatus, scannedDriver.MobileNum, scannedDriver.LicenceNumber)
		driverList[driverID] = scannedDriver
	}
}
