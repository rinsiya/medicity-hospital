package dto

type DoctorProfileSubmission struct {
	ProfessionalRole   string                     `json:"professional_role"`
	DepartmentID       uint                       `json:"department_id"`
	HospitalAddress    string                     `json:"hospital_address"`
	Experience         uint                       `json:"experience"`
	Biography          string                     `json:"biography"`

	Qualifications []DoctorQualificationInput      `json:"qualifications"`
	Experiences    []DoctorExperienceInput         `json:"experiences"`
}

type DoctorQualificationInput struct {
    Qualification string `json:"qualification"`
    Branch        string `json:"branch"`
    University    string `json:"university"`
    Year          uint16 `json:"year"`
}

type DoctorExperienceInput struct {
	Experience string `json:"experience"`
}
type DoctorBasicProfileResponse struct {
    Name  string `json:"name"`
    Email string `json:"email"`
    Phone string `json:"phone"`
}