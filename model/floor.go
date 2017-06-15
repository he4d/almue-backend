package model

import "errors"

//Floor represents a floor in the home
type Floor struct {
	ModelBase
	Description string     `json:"description"`
	Shutters    []Shutter  `json:"shutters"`
	Lightings   []Lighting `json:"lightings"`
}

//Validate validates the given floor
func (f Floor) Validate() error {
	return errors.New("not implemented")
}
