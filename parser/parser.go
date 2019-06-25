package parser

import (
	"encoding/json"
	"io/ioutil"
)

// PinInformation saves the information - which was provided by the .json file - for each pin.
type PinInformation struct {
	Name string
	Id   int
	Pin  int
}

// Configuration contains the complete result from the parsing process of the .json file.
type Configuration struct {
	Pins []PinInformation
}

// Parse takes a json file and returns a struct.
func Parse(filename string) (Configuration, error) {
	// res contains a byte array with the content of the file.
	res, err := ioutil.ReadFile(filename)

	if err != nil {
		return Configuration{}, err
	}

	var config Configuration
	err = json.Unmarshal(res, &config)

	if err != nil {
		return Configuration{}, err
	}

	return config, nil
}
