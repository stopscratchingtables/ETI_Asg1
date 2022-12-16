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

var x Passenger
var url = "http:/localhost:6969/api/v1/passenger"

var passList map[string]Passenger = map[string]Passenger{}

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
	router.HandleFunc("/api/v1/passenger", passengerFilter)
	router.HandleFunc("/api/v1/passenger/{passengerID}", passenger).Methods("GET", "DELETE", "PUT", "POST")
	fmt.Println("Listening at port 6969")
	log.Fatal(http.ListenAndServe(":6969", router))

}

func allPassenger(w http.ResponseWriter, r *http.Request) {
	db, err := sql.Open("mysql", "user:password@tcp(127.0.0.1:3306)/my_db")
	if err != nil {
		panic(err.Error())
	}

	defer db.Close()

	GetDB(db)

	data, err := json.Marshal(map[string]map[string]Passenger{"Passenger": passList})
	if len(passList) != 0 {
		fmt.Fprintf(w, "%s\n", data)
	}

}

func passenger(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-type", "application/json")

	params := mux.Vars(r)
	passVal := params["passengerID"]

	db, err := sql.Open("mysql", "user:password@tcp(127.0.0.1:3306)/my_db")
	// handle error
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	GetDB(db)
	e, err := json.Marshal(passList)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(e))

	if passDetails, exists := passList[passVal]; exists {
		if r.Method == "GET" { // GET METHOD

			data, _ := json.Marshal(passDetails)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Fprintf(w, "%s\n", data)

		} else if r.Method == "DELETE" {

			_, err := db.Exec("DELETE FROM PASSENGER WHERE ID=?", passVal)
			if err != nil {
				panic(err.Error())
			}
			fmt.Fprintf(w, passVal, "Deleted")

		} else if r.Method == "POST" {

			w.WriteHeader(http.StatusConflict)

		} else if r.Method == "PUT" {
			updatedPass := Passenger{}
			reqBody, _ := ioutil.ReadAll(r.Body)
			json.Unmarshal(reqBody, &updatedPass)

			_, err := db.Exec("UPDATE Passenger SET FirstName=?, LastName=?, MobileNum=?, EmailAddress=? WHERE ID=?", updatedPass.FirstName, updatedPass.LastName, updatedPass.MobileNum, updatedPass.EmailAddr, passVal)
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

			newPass := Passenger{}
			reqBody, _ := ioutil.ReadAll(r.Body)
			json.Unmarshal(reqBody, &newPass)
			_, err := db.Exec("INSERT into Passenger (ID, FirstName, LastName, MobileNum, EmailAddress) value(?, ?, ?, ?, ?, ?)", passVal, newPass.FirstName, newPass.LastName, newPass.MobileNum, newPass.EmailAddr)
			if err != nil {
				panic(err.Error())
			}
			passList[passVal] = newPass
			w.WriteHeader(http.StatusAccepted)

		} else if r.Method == "PUT" {
			// error message
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprintf(w, "Passenger ID Does Not Exist!")

		}

	}
}

// If there is no query, list all passengers.
// If there is query, list "passengerFilter"
func passengerFilter(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-type", "application/json")

	db, err := sql.Open("mysql", "user:password@tcp(127.0.0.1:3306)/my_db")
	// handle error
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()
	GetDB(db)

	query := r.URL.Query()
	results := map[string]Passenger{}

	if value := query.Get("q"); len(value) > 0 {
		for k, v := range passList { // e.g. k="IT", v="{"Diplpma in IT, ......, ...., ....."}"
			if strings.Contains(strings.ToLower(v.FirstName), strings.ToLower(value)) {
				results[k] = v
			}
		}
		if len(results) == 0 {
			w.WriteHeader(http.StatusNotFound)
		} else {
			// prints results
			data, err := json.Marshal(map[string]map[string]Passenger{"Passenger": results})
			if err != nil {
				log.Fatal(err)
			}
			if len(results) != 0 {
				fmt.Fprintf(w, "%s\n", data)
			}
		}
	} else if value := query.Get("history"); len(value) > 0 {

		var inProgRideList map[string]Ride = map[string]Ride{}

		if value != "" {
			results, err := db.Query("SELECT * FROM Ride WHERE PassengerID = '" + value + "' ORDER BY StartDateTime DESC")
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

	} else {

		data, err := json.Marshal(map[string]map[string]Passenger{"Passenger": passList}) // prints "Driver" + json of the drivers
		if err != nil {
			log.Fatal(err)
		}
		if len(passList) != 0 {
			fmt.Fprintf(w, "%s\n", data)
		}

	}

}

func GetDB(db *sql.DB) {
	results, err := db.Query("SELECT * FROM Passenger")
	if err != nil {
		fmt.Println(err)
	}
	for results.Next() {
		var passID string
		var scannedPass Passenger
		err = results.Scan(&passID, &scannedPass.FirstName, &scannedPass.LastName, &scannedPass.MobileNum, &scannedPass.EmailAddr)
		if err != nil {
			fmt.Println(err)
		}
		passList[passID] = scannedPass
	}
}
