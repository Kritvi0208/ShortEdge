package model

import "time"

type Visit struct {
	Code      string    `json:"code"`
	Timestamp time.Time `json:"timestamp"`
	IP        string    `json:"ip"`
	Country   string    `json:"country"`
	Browser   string    `json:"browser"`
	Device    string    `json:"device"`
}
