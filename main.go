package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/benCoder01/automata-backend/control"
	"github.com/stianeikeland/go-rpio"
)

var config *control.Configuration

func initialize(filename string) error {
	err := control.ParseFromJSON(filename)

	if err != nil {
		return err
	}

	fmt.Println("Finished initialisation.")

	return nil
}

func handleRequests() {
	http.HandleFunc("/activate", trigger)
	http.HandleFunc("/add", addControl)
	http.HandleFunc("/update", updateControl)
	http.HandleFunc("/delete", deleteControl)

	fmt.Println("Serving...")

	log.Fatal(http.ListenAndServe(":3000", nil))
}

func main() {
	// get quantity from command line arguments
	if len(os.Args) != 2 {
		fmt.Println("Please specify configuration file.")
		return
	}

	err := rpio.Open()

	if err != nil {
		panic(fmt.Sprint("unable to open gpio", err.Error()))
	}

	defer rpio.Close()

	err = initialize(os.Args[1])

	if err != nil {
		panic(err)
	}

	handleRequests()
}
