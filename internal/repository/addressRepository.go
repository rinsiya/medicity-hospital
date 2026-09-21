package repository

import (
	"medicity/database"
	"medicity/internal/models"

)

type AddressRepository interface {
	CreateAddress(address *models.Address) error
	GetAddressByPatientID(patientID uint) (*models.Address, error)
	UpdateAddress(address *models.Address) error
}

type addressRepository struct {}

func NewAddressRepository() AddressRepository {
	return &addressRepository{}
}

func (r *addressRepository) CreateAddress(address *models.Address) error {
	return database.DB.Create(address).Error
}

func (r *addressRepository) GetAddressByPatientID(patientID uint) (*models.Address, error) {
	var address models.Address

	err := database.DB.Where("patient_id = ?", patientID).
		First(&address).Error

	if err != nil {
		return nil, err
	}

	return &address, nil
}

func (r *addressRepository) UpdateAddress(address *models.Address) error {
	return database.DB.Save(address).Error
}