package request

import "io"

// AddControl describes the body for the specific request.
type AddControl struct {
	Name string
	Pin  int
}

// ParseAddControl returns the specific request parsed from the body of a http request.
func ParseAddControl(body io.ReadCloser) (AddControl, error) {
	return AddControl{"Test", 1}, nil
}

// DeleteControl describes the body for the specific request.
type DeleteControl struct {
	ID int
}

// ParseDeleteControl returns the specific request parsed from the body of a http request.
func ParseDeleteControl() (AddControl, error) {
	return AddControl{"Test", 1}, nil

}

// UpdateControl describes the body for the specific request.
type UpdateControl struct {
	Name string
	Pin  int
	ID   int
}

// ParseUpdateControl returns the specific request parsed from the body of a http request.
func ParseUpdateControl() (AddControl, error) {
	return AddControl{"Test", 1}, nil

}
