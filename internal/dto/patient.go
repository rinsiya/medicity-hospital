package dto

import "time"

type PatientMiniProfile struct{
PatientID      uint `json:"patient_id"` 
PatientName    string `json:"patient_name"`
ProfilePhoto   string `json:"profile_photo"`
}

type PatientProfile struct {
	
	PatientName string `json:"patient_name"`

	Gender string `json:"gender"`

	DOB time.Time `json:"dob"`

	Height float32 `json:"height"`

	Weight float32 `json:"weight"`
ProfilePhoto   string `json:"profile_photo"`


}
