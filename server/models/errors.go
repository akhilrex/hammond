package models

import "fmt"

type VehicleAlreadyExistsError struct {
	Registration string
}

func (e *VehicleAlreadyExistsError) Error() string {
	return fmt.Sprintf("Vehicle with this url already exists")
}
