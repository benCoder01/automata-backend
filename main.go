package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/benCoder01/automata-backend/parser"
	"github.com/benCoder01/automata-backend/request"
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

func generateID() int {
	var id int
	id = 0

	searching := true

	for searching {
		for _, control := range controls {
			if control.id == id {
				id++
				break
			}
		}
		searching = false
		return id
	}

	return 0
}

// getConfig will return the whole .json config
func getConfig(w http.ResponseWriter, r *http.Request) {

}

// addControl will add a new control to the slice and will update the .json file.
func addControl(w http.ResponseWriter, r *http.Request) {
	req, err := request.ParseAddControl(r.Body)

	if err != nil {
		// TODO: Error
	}

	var pin rpio.Pin
	pin = rpio.Pin(req.Pin)
	pin.Output()
	pin.High()

	// Create Control variable
	control := Control{activated: false, name: req.Name, pin: pin, id: generateID()}

	controls = append(controls, control)

	// TODO: update .json file
}

// deleteControl will delete a control from the controls slice.
// It will also upadet the config file
func deleteControl(w http.ResponseWriter, r *http.Request) {

}

// updateControl handles the route for updating a certain control.
// Therefor it also has to update the .json file, which holds the config.
func updateControl(w http.ResponseWriter, r *http.Request) {

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

// send is used if you only need to send a small amount of data back to the client. For example if a request suceeded.
func send(status int, text string, w http.ResponseWriter) {
	w.WriteHeader(status)
	w.Write([]byte(text))
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
