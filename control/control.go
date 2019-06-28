package control

import (
	"encoding/json"
	"io/ioutil"

	rpio "github.com/stianeikeland/go-rpio"
)

// Control contains information for each switch and the activator for the specific gpio pin.
type Control struct {
	ID        int
	Name      string
	Activated bool
	Pin       rpio.Pin
}

// Trigger activates a control on the gpio pin.
// If the pin is already activated, the method deactivates the pin.
func (c *Control) Trigger() error {
	c.Pin.Toggle()
	c.Activated = !c.Activated
	return nil
}

// Configuration contains each control.
type Configuration struct {
	Controls []Control
}

// AppendControl adds a new control field at the end of the slice.
func (config *Configuration) AppendControl(control *Control) {
	config.Controls = append(config.Controls, *control)
}

// GetWithID returns the index of the control object with the provided id.
// If the object does not exist, the method returns -1.
func (config *Configuration) GetWithID(id int) int {
	for index, control := range config.Controls {
		if control.ID == id {
			return index
		}
	}

	return -1
}

// structs needed to parse from .json file.

type pinInformationJSON struct {
	name string
	id   int
	pin  int
}

type configurationJSON struct {
	pins []pinInformationJSON
}

// Config holds the current configuration.
var Config *Configuration

// ParseFromJSON takes a string as a filename and returns a pointer to a Configuration struct.
func ParseFromJSON(filename string) error {

	// res contains a byte array with the content of the file.
	res, err := ioutil.ReadFile(filename)

	if err != nil {
		return err
	}

	var configJSON configurationJSON
	err = json.Unmarshal(res, &configJSON)

	if err != nil {
		return err
	}

	var controls []Control

	for _, val := range configJSON.pins {
		var pin rpio.Pin
		pin = rpio.Pin(val.pin)
		pin.Output()
		pin.High()

		controls = append(controls, Control{ID: val.id, Name: val.name, Activated: false, Pin: pin})

	}

	Config = &Configuration{Controls: controls}

	return nil
}

// GetConfig returns the config object.
func GetConfig() *Configuration {
	return Config
}
