package repository

import (
	"medicity/database"
	"medicity/internal/models"

	"gorm.io/gorm"
)



type DoctorRepository interface {
	CreateDoctor(tx *gorm.DB, doctor *models.Doctor) error
	IsProfileComplete(doctorID uint) (bool, error)
	GetDoctorByUserID(userID uint) (*models.Doctor, error)
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
		Where("id = ?", doctorID).
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