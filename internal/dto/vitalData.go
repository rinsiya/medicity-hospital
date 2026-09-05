package dto

type Vital struct {
	Name       string `json:"name"`
	Value      string `json:"value"`
	Unit       string `json:"unit"`
	RecordedAt string `json:"recorded_at"`
}

type SavePatientVitalsRequest struct {
	Vitals []Vital `json:"vitals" binding:"required"`
}