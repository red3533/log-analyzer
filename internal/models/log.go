package models

import "time"

type Log struct {
	IP        string
	Timestamp time.Time
	Method    string
	URL       string
	Status    int
	SizeByte  int
}
