package repository

import (
	"errors"
	"medicity/database"
	"medicity/internal/models"

	"gorm.io/gorm"
)

type UserRepository interface {
	FindByUsername(username string) (*models.User, error)
	FindByEmailAndPhone(email,phone string) (*models.User, error)
	FindByEmailOrPhone(email,phone string) (*models.User, error)
	FindByPhone(phone,role string) (*models.User, error)		
	PhoneExist(phone string) (bool, error)
	Create(tx *gorm.DB, user *models.User) error
	FetchUser(userID uint) (*models.User, error)
    Update(user *models.User) error
    GetPatientIDByUserID(userID uint) (uint, error)
    GetUserByID(userID uint)(*models.User,error)
}

  type userRepository struct{}

  func NewUserRepository() UserRepository {
  	return &userRepository{}
  }

func (r *userRepository) GetUserByID(userID uint)(*models.User,error){

var user models.User

	err := database.DB.Where("user_id = ?", userID).First(&user).Error

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *userRepository) GetPatientIDByUserID(userID uint) (uint, error) {

	var patient models.Patient

	err := database.DB.
		Select("patient_id").
		Where("user_id = ?", userID).
		First(&patient).Error

	if err != nil {
		return 0, err
	}

	return patient.PatientID, nil
}

func (r *userRepository) FindByEmailOrPhone(email, phone string) (*models.User, error) {
    var user models.User

    err := database.DB.
        Where("email = ? OR phone = ?", email, phone).
        First(&user).Error

    if errors.Is(err, gorm.ErrRecordNotFound) {
        return nil, nil
    }

    if err != nil {
        return nil, err
    }

    return &user, nil
}
func (r *userRepository) FindByUsername(username string) (*models.User, error) {
    var user models.User

    err := database.DB.
        Where("email = ? OR phone = ?", username, username).
        First(&user).Error

    if errors.Is(err, gorm.ErrRecordNotFound) {
        return nil, nil
    }

    if err != nil {
        return nil, err
    }

    return &user, nil
}

func (r *userRepository) Update(user *models.User) error {
	return database.DB.Save(user).Error
}
func (r *userRepository) FindByPhone(phone,role string) (*models.User, error) {
    var user models.User
    err := database.DB.Where("phone = ? AND role = ?", phone, role  ).First(&user).Error

if errors.Is(err, gorm.ErrRecordNotFound) {
        return nil, nil
    }

    if err != nil {
        return nil, err
    }

    return &user, nil

}

func (r *userRepository) PhoneExist(phone string) (bool, error) {
	var count int64

		err := database.DB.Model(&models.User{}).Where("phone = ?", phone).Count(&count).Error
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

func (r *userRepository) FindByEmailAndPhone(email, phone string) (*models.User, error) {
    var user models.User

    err := database.DB.
        Where("email = ? AND phone = ?", email, phone).
        First(&user).Error

    if errors.Is(err, gorm.ErrRecordNotFound) {
        return nil, nil
    }

    if err != nil {
        return nil, err
    }

    return &user, nil
}

func (r *userRepository) Create(tx *gorm.DB,user *models.User) error {

    return tx.Create(user).Error
}

func (r *userRepository) FetchUser(userID uint) (*models.User, error) {

	var user models.User

	err := database.DB.
		Where("user_id = ?", userID).
		First(&user).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return &user, nil
}