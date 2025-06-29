package model

import "time"

type URL struct {
	Code       string     `json:"code"`
	LongURL    string     `json:"long_url"`
	CreatedAt  time.Time  `json:"created_at"`
	Visibility string     `json:"visibility"` // "public" or "private"
	ExpiresAt  *time.Time `json:"expires_at"` // 
}

type ShortenRequest struct {
	LongURL    string     `json:"long_url"`
	CustomCode string     `json:"custom_code"` // Optional
	Visibility string     `json:"visibility"`  // public / private
	ExpiresAt  *time.Time `json:"expires_at"`  // Optional 
}
