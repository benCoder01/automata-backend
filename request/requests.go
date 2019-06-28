package request

import (
	"encoding/json"
	"io"
)

// AddControl describes the body for the specific request.
type AddControl struct {
	Name string
	Pin  int
}

// ParseAddControl returns the specific request parsed from the body of a http request.
func ParseAddControl(body io.ReadCloser) (AddControl, error) {
	decoder := json.NewDecoder(body)

	var addControl AddControl

	err := decoder.Decode(&addControl)

	if err != nil {
		return AddControl{}, err
	}

	return addControl, nil
}

// DeleteControl describes the body for the specific request.
type DeleteControl struct {
	ID int
}

// ParseDeleteControl returns the specific request parsed from the body of a http request.
func ParseDeleteControl(body io.ReadCloser) (DeleteControl, error) {
	decoder := json.NewDecoder(body)

	var deleteControl DeleteControl

	err := decoder.Decode(&deleteControl)

	if err != nil {
		return DeleteControl{}, err
	}

	return deleteControl, nil

}

// UpdateControl describes the body for the specific request.
type UpdateControl struct {
	Name string
	Pin  int
	ID   int
}

// ParseUpdateControl returns the specific request parsed from the body of a http request.
func ParseUpdateControl(body io.ReadCloser) (UpdateControl, error) {
	decoder := json.NewDecoder(body)

	var updateControl UpdateControl

	err := decoder.Decode(&updateControl)

	if err != nil {
		return UpdateControl{}, err
	}

	return updateControl, nil
}
