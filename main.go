package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/benCoder01/automata-backend/control"
	"github.com/benCoder01/automata-backend/handler"
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
	// TODO: Test if http methods are correct
	http.HandleFunc("/activate", handler.Trigger)
	http.HandleFunc("/add", handler.AddControl)
	http.HandleFunc("/update", handler.UpdateControl)
	http.HandleFunc("/delete", handler.DeleteControl)
	http.HandleFunc("/config", handler.GetConfig)

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
