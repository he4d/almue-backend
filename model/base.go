package model

import "time"

//Base is a basemodel for all database models
type Base struct {
	ID       int64     `json:"id" valid:"-"`
	Created  time.Time `json:"created" valid:"-"`
	Modified time.Time `json:"modified" valid:"-"`
}
