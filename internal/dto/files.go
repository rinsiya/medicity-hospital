package dto

import "time"

type PatientFiles struct{
			PploadedAt    time.Time `json:"uploaded_at"`
			FileName      string    `json:"file_name"`
			FileCategory  string    `json:"file_category"`
			FileType      string    `json:"file_type"`

}