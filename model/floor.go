package model

//Floor represents the database object of a floor
type Floor struct {
	Base
	Description *string `json:"description"`
}
