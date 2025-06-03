package models

import "time"

type LogParsed struct {
	IP        string
	Timestamp time.Time
	Method    string
	Status    int
}
