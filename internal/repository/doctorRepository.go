package repository

import (
	"medicity/database"
	"medicity/internal/dto"
	"medicity/internal/models"
	

	"gorm.io/gorm"
)



type DoctorRepository interface {
	CreateDoctor(tx *gorm.DB, doctor *models.Doctor) error
	IsProfileComplete(doctorID uint) (bool, error)
	GetDoctorByUserID(userID uint) (*models.Doctor, error)
    SubmitDoctorProfile(doctorID uint,input dto.DoctorProfileSubmission) error 
     GetDoctorBasicProfile(userID uint) (*dto.DoctorBasicProfileResponse, error)
}

  type doctorRepository struct{}

  func NewDoctorRepository() DoctorRepository {
  	return &doctorRepository{}
  }


  func(r *doctorRepository) CreateDoctor(tx *gorm.DB, doctor *models.Doctor) error{

		return tx.Create(doctor).Error;


  }
  func (r *doctorRepository) IsProfileComplete(doctorID uint) (bool,error) {
	var doctor models.Doctor

	err := database.DB.
		Where("doctor_id = ?", doctorID).
		First(&doctor).Error

	if err != nil {
		return false,err
	}
if doctor.DepartmentID !=nil{
	return false,nil
}
	return true,nil
}

func (r *doctorRepository) GetDoctorByUserID(userID uint) (*models.Doctor, error) {
	var doctor models.Doctor

	err := database.DB.
		Where("user_id = ?", userID).
		First(&doctor).Error

	if err != nil {
		return nil, err
	}

	return &doctor, nil
}

func (r *doctorRepository) SubmitDoctorProfile(doctorID uint,input dto.DoctorProfileSubmission) error {

    return database.DB.Transaction(func(tx *gorm.DB) error {

        // 1. Update Doctor

        err := tx.Model(&models.Doctor{}).
            Where("doctor_id = ?", doctorID).
            Updates(map[string]interface{}{
                "professional_role": input.ProfessionalRole,
                "department_id":     input.DepartmentID,
            }).Error

        if err != nil {
            return err
        }

        // 2. Create/update DoctorProfile

        profile := models.DoctorProfile{
            DoctorID:            doctorID,
            HospitalAddress:     input.HospitalAddress,
            Experience:          input.Experience,
            Biography:           input.Biography,
        }

        // Create or update profile
        err = tx.
            Where("doctor_id = ?", doctorID).
            Assign(profile).
            FirstOrCreate(&profile).Error

        if err != nil {
            return err
        }

        // 3. Save qualifications

        for _, qualification := range input.Qualifications {

            doctorQualification := models.DoctorQualification{
                DoctorID:    doctorID,
                Qualification: qualification.Qualification,
                University:   qualification.University,
                Year:          qualification.Year,
            }

            if err := tx.Create(&doctorQualification).Error; err != nil {
                return err
            }
        }

        // 4. Replace professional experiences
if err := tx.
    Where("doctor_id = ?", doctorID).
    Delete(&models.ProfessionalExperience{}).Error; err != nil {
    return err
}

for _, experience := range input.Experiences {

    professionalExperience := models.ProfessionalExperience{
        DoctorID: doctorID,
        Details:  experience.Experience,
    }

    if err := tx.Create(&professionalExperience).Error; err != nil {
        return err
    }
}

        // 4. Change verification status. 

        err = tx.Model(&models.Doctor{}). 
            Where("doctor_id = ?", doctorID).
            Update(
                "verification_status",
                models.PendingVerification,
            ).Error  

        if err != nil { 
            return err  
        }

        return nil
    })
}

func (r *doctorRepository) GetDoctorBasicProfile(userID uint) (*dto.DoctorBasicProfileResponse, error) {

    var profile dto.DoctorBasicProfileResponse
 
    err := database.DB.
        Table("users").
        Select(`
            CONCAT(doctors.first_name, ' ', doctors.last_name) AS name,
            users.email,
            users.phone
        `).
        Joins("JOIN doctors ON doctors.user_id = users.user_id").
        Where("users.user_id = ?", userID).
        Scan(&profile).Error

    if err != nil {
        return nil, err
    }

    return &profile, nil
}