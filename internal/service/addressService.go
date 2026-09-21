package service


import (
	"errors"

	"medicity/internal/models"
	"medicity/internal/repository"

	"gorm.io/gorm"
)

type AddressService interface {
	AddOrUpdateAddress(patientID uint, address, place, country string) error
}

type addressService struct {
	addressRepository repository.AddressRepository
}

func NewAddressService(
	addressRepository repository.AddressRepository,
) AddressService {
	return &addressService{
		addressRepository: addressRepository,
	}
}

func (s *addressService) AddOrUpdateAddress(patientID uint,address string,place string,country string) error {

	// Address is optional.
	// If all fields are empty, do nothing.
	if address == "" && place == "" && country == "" {
		return nil
	}

	existingAddress, err := s.addressRepository.GetAddressByPatientID(patientID)

	if err != nil {

		// Address does not exist → create new address
		if errors.Is(err, gorm.ErrRecordNotFound) {

			newAddress := &models.Address{
				PatientID: patientID,
				Address:  address,
				Place:    place,
				Country:  country,
			}

			return s.addressRepository.CreateAddress(newAddress)
		}

		return err
	}

	// Address already exists → update it
	existingAddress.Address = address
	existingAddress.Place = place
	existingAddress.Country = country

	return s.addressRepository.UpdateAddress(existingAddress)
}