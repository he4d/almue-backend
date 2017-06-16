package model

//Floor represents the database object of a floor
type Floor struct {
	Base        `valid:"-"`
	Description string      `json:"description" valid:"alphanum,required"`
	Shutters    []*Shutter  `json:"shutters" valid:"-"`
	Lightings   []*Lighting `json:"lightings" valid:"-"`
}
