package dto

import "time"

type PatientMiniProfile struct {
	PatientID    uint   `json:"patient_id"`
	PatientName  string `json:"patient_name"`
	ProfilePhoto string `json:"profile_photo"`
}

type PatientProfile struct {
	PatientID         uint   `json:"patient_id"`
	PatientName       string `json:"patient_name"`
	IsProfileComplete bool   `json:"is_profile_complete"`
	Gender            string `json:"gender"`

	DOB time.Time `json:"dob"`

	Height float32 `json:"height"`

	Weight       float32 `json:"weight"`
	ProfilePhoto string  `json:"profile_photo"`
}
