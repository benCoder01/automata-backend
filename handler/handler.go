package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/benCoder01/automata-backend/control"
	"github.com/benCoder01/automata-backend/request"
	rpio "github.com/stianeikeland/go-rpio"
)

// GetConfig will return the whole .json config.
// However it will not read the .json file but will instead return the conrol.Config variable in json format.
// This guarantess that the program does not return old data.
// This could happen if the user requests the config file while this file gets updated.
func GetConfig(w http.ResponseWriter, r *http.Request) {
	json, err := json.Marshal(control.Config)

	if err != nil {
		// TODO: Replace send in error with http.Error(...)
		send(http.StatusInternalServerError, "Could not parse config to json.", w)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(json)
}

// AddControl will add a new control to the slice and will update the .json file.
func AddControl(w http.ResponseWriter, r *http.Request) {
	req, err := request.ParseAddControl(r.Body)

	if err != nil {
		send(http.StatusBadRequest, "Could not parse the request body.", w)
		return
	}

	var pin rpio.Pin
	pin = rpio.Pin(req.Pin)
	pin.Output()
	pin.High()

	// Create Control variable
	newControl := control.Control{Activated: false, Name: req.Name, Pin: pin, ID: generateID()}

	// add the new control to the existing slice of controls.
	config := control.GetConfig()
	config.AppendControl(&newControl)

	// TODO: update .json file
	json, err := json.Marshal(newControl)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(json)
}

// DeleteControl will delete a control from the controls slice.
// It will also upadet the config file
func DeleteControl(w http.ResponseWriter, r *http.Request) {
	req, err := request.ParseDeleteControl(r.Body)

	if err != nil {
		send(http.StatusBadRequest, "Could not parse the request body.", w)
		return
	}

	config := control.GetConfig()

	// Get the index of the element
	index := config.GetWithID(req.ID)

	if index == -1 {
		send(http.StatusBadRequest, "Variable is not valid", w)
		return
	}

	control := config.DeleteControl(index)

	// TODO: update .json file
	json, err := json.Marshal(control)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(json)

}

// UpdateControl handles the route for updating a certain control.
// Therefor it also has to update the .json file, which holds the config.
func UpdateControl(w http.ResponseWriter, r *http.Request) {
	req, err := request.ParseUpdateControl(r.Body)

	if err != nil {
		send(http.StatusBadRequest, "Could not parse the request body.", w)
		return
	}
	config := control.GetConfig()

	// Get the index of the element
	index := config.GetWithID(req.ID)

	control := config.UpdateControl(index, req.Name, req.Pin)

	// TODO: update .json file
	json, err := json.Marshal(control)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(json)

}

// Trigger activates and deactivates a control by using the gpio pin provided by the struct control.
func Trigger(w http.ResponseWriter, r *http.Request) {
	id, err := readIDFromURL(r)

	if err != nil {
		send(500, err.Error(), w)
		return
	}

	config := control.GetConfig()

	index := config.GetWithID(id)

	if index == -1 {
		send(http.StatusBadRequest, "Variable is not valid", w)
		return
	}

	err = config.Controls[index].Trigger()
	if err != nil {
		send(500, err.Error(), w)
		return
	}

	// TODO: Response as JSON response
	if config.Controls[index].Activated {
		fmt.Println(config.Controls[index].Name + " activated")
		send(200, "Control "+strconv.Itoa(id)+" - "+config.Controls[index].Name+" activated", w)
	} else {
		fmt.Println(config.Controls[index].Name + " deactivated")
		send(200, "Control "+strconv.Itoa(id)+" - "+config.Controls[index].Name+" deactivated", w)
	}
}

// send is used if you only need to send a small amount of data back to the client. For example if a request suceeded.
func send(status int, text string, w http.ResponseWriter) {
	w.WriteHeader(status)
	w.Write([]byte(text))
}

func generateID() int {

	config := control.GetConfig()

	var id int
	id = 0

	searching := true

	for searching {
		for _, control := range config.Controls {
			if control.ID == id {
				id++
				break
			}
		}
		searching = false
		return id
	}

	return 0
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
