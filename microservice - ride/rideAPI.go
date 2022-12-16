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

var x Ride
var url = "http:/localhost:5000/api/v1/ride"

var rideList map[string]Ride = map[string]Ride{}

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
	//ID              string `json:Ride ID`
	PassengerName     string `json:"Passenger"`
	Driver            string `json:"Driver"`
	TripStatus        string `json:"Trip Status"`
	StartDateTime     string `json:"Pick-Up Timing"`
	EndDateTime       string `json:"Drop-Off Timing"`
	PostalCodePickUp  string `json:"Pick up point"`
	PostalCodeDropOff string `json:"Drop off point"`
}

// FUNCTIONS

func main() {

	router := mux.NewRouter()
	router.HandleFunc("/api/v1/ride", rideFilter)
	//router.HandleFunc("/api/v1/ride", allRides)
	router.HandleFunc("/api/v1/ride/{rideID}", ride).Methods("GET", "DELETE", "PUT", "POST")
	fmt.Println("Listening at port 5000")
	log.Fatal(http.ListenAndServe(":5000", router))

}

func allRides(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-type", "application/json")
	db, err := sql.Open("mysql", "user:password@tcp(127.0.0.1:3306)/my_db")

	// handle error
	if err != nil {
		panic(err.Error())
	}

	defer db.Close()

	GetDB(db)

	querystringmap := r.URL.Query()
	fmt.Fprintf(w, "List All Rides\n\n")
	for k, v := range querystringmap {
		fmt.Fprintf(w, "%s = %v\n", k, v[0])
	}

}

func ride(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-type", "application/json")

	params := mux.Vars(r)
	rideVal := params["rideID"]

	db, err := sql.Open("mysql", "user:password@tcp(127.0.0.1:3306)/my_db")
	// handle error
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	GetDB(db)
	e, err := json.Marshal(rideList)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(e))

	if rideDetails, exists := rideList[rideVal]; exists {
		if r.Method == "GET" { // GET METHOD

			data, _ := json.Marshal(rideDetails)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Fprintf(w, "%s\n", data)

		} else if r.Method == "DELETE" {

			_, err := db.Exec("DELETE FROM RIDE WHERE ID=?", rideVal)
			if err != nil {
				panic(err.Error())
			}
			fmt.Fprintf(w, rideVal, "Deleted")

		} else if r.Method == "POST" {

			w.WriteHeader(http.StatusConflict)

		} else if r.Method == "PUT" {
			updatedRide := Ride{}
			reqBody, _ := ioutil.ReadAll(r.Body)
			json.Unmarshal(reqBody, &updatedRide)

			fmt.Println(rideVal)

			_, err := db.Exec("UPDATE Ride SET PassengerID=?, DriverID=?, TripStatus=?, StartDateTime=?, EndDateTime=?, PostalCodePickUp=?, PostalCodeDropOff=? WHERE ID=?", updatedRide.PassengerName, updatedRide.Driver, updatedRide.TripStatus, updatedRide.StartDateTime, updatedRide.EndDateTime, updatedRide.PostalCodePickUp, updatedRide.PostalCodeDropOff, rideVal)
			if err != nil {
				panic(err.Error())
			}
			w.WriteHeader(http.StatusAccepted)

		} else {
			fmt.Fprintf(w, "Invalid Ride ID")
		}
	} else {

		if r.Method == "POST" {

			newRide := Ride{}
			reqBody, _ := ioutil.ReadAll(r.Body)
			json.Unmarshal(reqBody, &newRide)
			_, err := db.Exec("INSERT INTO Ride (ID, PassengerID, DriverID, TripStatus, StartDateTime, EndDateTime, PostalCodePickUp, PostalCodeDropOff) VALUES (?, ?, ?, ?, ?, ?, ?, ?);", rideVal, newRide.PassengerName, newRide.Driver, newRide.TripStatus, newRide.StartDateTime, newRide.EndDateTime, newRide.PostalCodePickUp, newRide.PostalCodeDropOff)
			if err != nil {
				panic(err.Error())
			}
			rideList[rideVal] = newRide
			w.WriteHeader(http.StatusAccepted)

		} else if r.Method == "PUT" {
			// error message
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprintf(w, "Rider ID Does Not Exist!")

		}

	}
}

// If there is no query, list all Ruders.
// If there is query, list "riderFilter"

func rideFilter(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-type", "application/json")

	db, err := sql.Open("mysql", "user:password@tcp(127.0.0.1:3306)/my_db")
	// handle error
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()
	GetDB(db)

	query := r.URL.Query()
	results := map[string]Ride{}

	if value := query.Get("q"); len(value) > 0 {
		for k, v := range rideList { // e.g. k="IT", v="{"Diplpma in IT, ......, ...., ....."}"
			if strings.Contains(strings.ToLower(v.PassengerName), strings.ToLower(value)) {
				results[k] = v
			}
		}
		if len(results) == 0 {
			// fmt.Fprintf(w, "No Passengers Elligible")
			w.WriteHeader(http.StatusNotFound)
		} else {

			// prints results

			data, err := json.Marshal(map[string]map[string]Ride{"Passenger": results})
			if err != nil {
				log.Fatal(err)
			}
			if len(results) != 0 {
				fmt.Fprintf(w, "%s\n", data)
			}
		}
	} else if value := query.Get("Finished"); len(value) > 0 {

		var finishedRides = map[string]Ride{}

		db, err := sql.Open("mysql", "user:password@tcp(127.0.0.1:3306)/my_db")
		if err != nil {
			panic(err.Error())
		}
		defer db.Close()

		results, err := db.Query(`Select * from Ride WHERE TripStatus="Finished"`)
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
			finishedRides[rideID] = scannedRide
		}

		data, _ := json.Marshal(finishedRides)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Fprintf(w, "%s\n", data)

	} else if value := query.Get("tripstatus"); len(value) > 0 {

		var inProgRideList map[string]Ride = map[string]Ride{}

		if value == "inprogress" {
			results, err := db.Query("SELECT * FROM Ride WHERE TripStatus='In Progress'; ")
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
				fmt.Println(scannedRide.PassengerName, scannedRide.Driver, scannedRide.TripStatus, scannedRide.StartDateTime, scannedRide.EndDateTime, scannedRide.PostalCodePickUp, scannedRide.PostalCodeDropOff)
				inProgRideList[rideID] = scannedRide
			}

			data, _ := json.Marshal(map[string]map[string]Ride{"Ride": inProgRideList})
			if err != nil {
				log.Fatal(err)
			}
			fmt.Fprintf(w, "%s\n", data)
		}

	} else if value := query.Get("Open"); len(value) > 0 {

		var openRides = map[string]Ride{}

		db, err := sql.Open("mysql", "user:password@tcp(127.0.0.1:3306)/my_db")
		if err != nil {
			panic(err.Error())
		}
		defer db.Close()

		results, err := db.Query(`Select * from Ride WHERE TripStatus="Open"`)
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
			openRides[rideID] = scannedRide
		}

		data, _ := json.Marshal(openRides)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Fprintf(w, "%s\n", data)

	}

}

func GetDB(db *sql.DB) {
	results, err := db.Query("SELECT * FROM Ride")
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
		fmt.Println(rideID, scannedRide.PassengerName, scannedRide.Driver, scannedRide.TripStatus, scannedRide.StartDateTime, scannedRide.EndDateTime, scannedRide.PostalCodePickUp, scannedRide.PostalCodeDropOff)
		rideList[rideID] = scannedRide
	}
}
