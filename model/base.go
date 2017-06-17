package model

import "time"

//Base is a basemodel for all database models
type Base struct {
	ID       int64     `json:"id"`
	Created  time.Time `json:"created"`
	Modified time.Time `json:"modified"`
}
