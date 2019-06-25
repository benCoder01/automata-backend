package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/benCoder01/automata-backend/parser"
	"github.com/stianeikeland/go-rpio"
)

type Control struct {
	id        int
	name      string
	activated bool
	pin       rpio.Pin
}

func (c *Control) trigger() error {
	c.pin.Toggle()
	c.activated = !c.activated
	return nil
}

var controls []Control

func initialize(filename string) error {
	config, err := parser.Parse(filename)

	if err != nil {
		return err
	}

	// parse config into array of controls
	for _, val := range config.Pins {
		var pin rpio.Pin
		pin = rpio.Pin(val.Pin)
		pin.Output()
		pin.High()

		controls = append(controls, Control{val.Id, val.Name, false, pin})

	}

	fmt.Println("Finished initialisation.")

	return nil
}

func readIDFromURL(r *http.Request) (int, error) {
	// read id from url variable
	keys, ok := r.URL.Query()["id"]

	if !ok || len(keys[0]) < 1 {
		return 0, errors.New("No variable passed")
	}

	// parse string id to integer

	id, err := strconv.Atoi(keys[0])

	if err != nil {
		return 0, errors.New("Could not convert input to integer")
	}

	return id, nil
}

func trigger(w http.ResponseWriter, r *http.Request) {
	id, err := readIDFromURL(r)

	if err != nil {
		send(500, err.Error(), w)
		return
	}

	if id >= 0 && len(controls) < id+1 {
		send(http.StatusBadRequest, "Variable is not valid", w)
		return
	}

	err = controls[id].trigger()
	if err != nil {
		send(500, err.Error(), w)
		return
	}

	if controls[id].activated {
		fmt.Println(controls[id].name + " activated")
		send(200, "Control "+strconv.Itoa(id)+" - "+controls[id].name+" activated", w)
	} else {
		fmt.Println(controls[id].name + " deactivated")
		send(200, "Control "+strconv.Itoa(id)+" - "+controls[id].name+" deactivated", w)
	}
}

func send(status int, text string, w http.ResponseWriter) {
	w.WriteHeader(status)
	w.Write([]byte(text))
}

func handleRequests() {
	fmt.Println("Serving...")
	http.HandleFunc("/activate", trigger)
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
