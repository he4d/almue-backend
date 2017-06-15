package model

import "time"

type ModelBase struct {
	ID       int64     `json:"id"`
	Created  time.Time `json:"created"`
	Modified time.Time `json:"modified"`
}
