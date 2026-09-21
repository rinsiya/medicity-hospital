package service

import (
	"errors"
	"strings"
	"time"

	"medicity/database"
	"medicity/internal/dto"
	"medicity/internal/models"
	"medicity/internal/repository"
	"medicity/logger"
	"medicity/pkg/utils"
    "go.uber.org/zap"
	"gorm.io/gorm"
)

var (
	ErrEmailAndPhoneAlreadyExists       = errors.New("email and phone already registered")
	ErrPhoneAlreadyExists       = errors.New("phone already registered")
	ErrEmailOrPhoneAlreadyExists       = errors.New("email or phone already registered. Check the data entered carefully and try to login")
ErrEmailAndPhoneMismatch   = errors.New("Email Or phone registered ")
	ErrPhonePendingVerification = errors.New("phone has pending verification")
)

type SignupService interface {
	Signup(input *dto.SignupInput, role models.UserRole) (string,error)
	ResendOTP(phone string) error
	CreateUser(phone string) (*models.User,error)
	FindPendingUserByPhone(phone string)(*models.PendingUserSignup,error)
	UpdatePhone(oldPhone string,newPhone string) (models.UserRole,string,error)
	//CreateUser(phone string,otp string) (time.Time, error)
SaveOTPPendingUser(otpHash,phone string) error
}

type signupService struct {
	pendingUserRepo repository.PendingUserSignupRepository
	userRepo        repository.UserRepository
	patientRepo     repository.PatientRepository
	doctorRepo      repository.DoctorRepository
	otpService      OTPService
}


func NewSignupService(
    otpservice OTPService,
    pendingUserRepo repository.PendingUserSignupRepository,
    userRepo repository.UserRepository,
    patientRepo repository.PatientRepository,
    doctorRepo repository.DoctorRepository,
) SignupService {

    return &signupService{
        otpService:      otpservice,
        pendingUserRepo: pendingUserRepo,
        userRepo:        userRepo,
        patientRepo:     patientRepo,
        doctorRepo:      doctorRepo,
    }
}
func (s *signupService) Signup(input *dto.SignupInput, role models.UserRole) (string,error) {

	email := strings.ToLower(strings.TrimSpace(input.Email))
	phone := strings.TrimSpace(input.Phone)

	logger.Log.Info("Signup service started",
		zap.String("role", string(role)),
		zap.String("email", email),
	)

	// Check existing user by email

	user, err := s.userRepo.FindByEmailOrPhone(email,phone)
	if err != nil {
		logger.Log.Error(
			"Failed to check email and phone existence",
			zap.String("email", email),
			zap.String("phone", phone),
			zap.Error(err),
		)
		return phone,err
	}
	if user!=nil {
				logger.Log.Warn(
			"Signup rejected: email or phone already registered.",
			zap.String("email", email),
			zap.String("phone", phone),
			zap.String("role", string(role)),
		)
if user.Phone ==phone {
	if user.Email==email{
		return user.Phone,ErrEmailAndPhoneAlreadyExists
	}
			return user.Phone,ErrEmailOrPhoneAlreadyExists

}
if user.Email==email{
			return user.Phone,ErrEmailOrPhoneAlreadyExists

}


		return user.Phone,ErrEmailAndPhoneAlreadyExists
	}

	//Check pending signup by email

	pendingUser, err := s.pendingUserRepo.FindByEmailOrPhone(email,phone)
	if err != nil {

		logger.Log.Error(
			"Failed to check pending signup by email and phone",
			zap.String("email", email),
			zap.String("phone", phone),
			zap.Error(err),
		)

		return phone,err
	}

	if pendingUser != nil {
						logger.Log.Warn(
			"Signup rejected: email or phone already registered.phone verification pending.",
			zap.String("email", email),
			zap.String("phone", phone),
			zap.String("role", string(role)),
		)
if pendingUser.Phone ==phone {
					logger.Log.Warn(
			"Signup rejected:phone already registered.phone verification pending.",
					)
	if pendingUser.Email==email{
		logger.Log.Warn(
			"Signup rejected:email and phone already registered.phone verification pending.",
					)
					logger.Log.Warn(
			"Pending user exist .complete phone verification.",
					)
		   SessionID,err := s.otpService.SendOTP(phone)
s.SaveOTPPendingUser(SessionID,phone)
    if err != nil {
        return pendingUser.Phone, err
    }
    logger.Log.Info(
        "OTP for pending user send successfully",
        zap.String("phone", phone),
        
    )
        return pendingUser.Phone, nil

    
}
		return pendingUser.Phone,ErrEmailAndPhoneMismatch
	}
		
if pendingUser.Email==email{
			return pendingUser.Phone,ErrEmailAndPhoneMismatch

}
	}

	

	// Hash password
	hashedPassword, err := utils.HashPassword(input.Password)
	if err != nil {

		logger.Log.Error(
			"Failed to hash signup password",
			zap.String("email", email),
			zap.Error(err),
		)

		return phone,err
	}

	
	sessionID, err := s.otpService.SendOTP(phone)
	//GenerateAndSendOTP(phone)
	if err != nil {

		logger.Log.Error(
			"Failed to send signup OTP",
			zap.String("email", email),
			zap.Error(err),
		)

		return phone,err
		}
	pendingSignup := &models.PendingUserSignup{
		Role:         role,
		FirstName:    input.FirstName,
		LastName:     input.LastName,
		Email:        email,
		Phone:        phone,
		Password:     hashedPassword,
		SessionID:      sessionID,
		OTPExpiresAt: time.Now().Add(2 * time.Minute),
		CreatedAt:    time.Now(),
	}

	err = s.pendingUserRepo.Create(pendingSignup)
	if err != nil {

		logger.Log.Error(
			"Failed to create pending signup",
			zap.String("email", email),
			zap.String("role", string(role)),
			zap.Error(err),
		)

		return phone,err
	}
	logger.Log.Info(
		"Signup pending OTP verification",
		zap.String("email", email),
		zap.String("role", string(role)),
	)

	return phone,nil
}

func (s *signupService) ResendOTP(phone string) error {

    SessionID,err := s.otpService.SendOTP(phone)
s.SaveOTPPendingUser(SessionID,phone)
    if err != nil {
        return err
    }
    logger.Log.Info(
        "OTP resent successfully",
        zap.String("phone", phone),
    )
    return nil
}


func (s *signupService) SaveOTPPendingUser(sessionID, phone string) (error) {
	pendingUser, err := s.pendingUserRepo.FindByPhone(phone)

	if err != nil {

		logger.Log.Error(
			"Failed to find pending user by phone",
			zap.String("phone", phone),
			zap.Error(err),
		)

		return err
	}
	pendingUser.SessionID= sessionID
	pendingUser.OTPExpiresAt = time.Now().Add(2 * time.Minute)

	err = s.pendingUserRepo.Update(pendingUser)
if err!=nil{
			logger.Log.Error(
			"Failed to save otp to pendinguser",
			zap.String("phone", phone),
			zap.Error(err),
		)
	return err
}
	logger.Log.Info("Signup pending OTP verification",
	zap.String("phone", phone),
	zap.String("role", string(pendingUser.Role)),
	)

	return nil

}
func (s *signupService) FindPendingUserByPhone(phone string) (*models.PendingUserSignup,error){
	pendingUser, err := s.pendingUserRepo.FindByPhone(phone)
	if err != nil {

		logger.Log.Error(
			"Failed to find pending user by phone",
			zap.String("phone", phone),
			zap.Error(err),
		)

		return nil,err
	}
	return pendingUser,err
}


func (s *signupService) CreateUser(phone string) (*models.User,error) {

	pendingUser, err := s.pendingUserRepo.FindByPhone(phone)


    if err != nil {
        logger.Log.Error(
            "Failed to fetch signup data",
            zap.String("phone", phone),
            zap.Error(err),
        )
        return nil,err
    }

    if pendingUser == nil {
        return nil,errors.New("pending signup not found")
    }
	var userID uint
	err=database.DB.Transaction(func(tx *gorm.DB) error {

		user := &models.User{
			Role:     pendingUser.Role,
			Email:    pendingUser.Email,
			Phone:    pendingUser.Phone,
			Password: pendingUser.Password,
		}

		if err := s.userRepo.Create(tx, user); err != nil {
			return err
		}
			userID = user.UserID
		if pendingUser.Role == models.RolePatient {

			patient := &models.Patient{
				UserID:    user.UserID,
				FirstName: pendingUser.FirstName,
				LastName:  pendingUser.LastName,
			}
			if err := s.patientRepo.CreatePatient(tx, patient); err != nil {
				return err
			}

		} else if pendingUser.Role == models.RoleDoctor {

			doctor := &models.Doctor{
				UserID:    user.UserID,
				FirstName: pendingUser.FirstName,
				LastName:  pendingUser.LastName,
			}

			if err := s.doctorRepo.CreateDoctor(tx, doctor); err != nil {
				return err
			}
		}

		// Remove temporary signup record

		if err = s.pendingUserRepo.Delete(tx, pendingUser.PendingID); err != nil {
			return err
		}

		return nil

	})

if err!=nil{
	return nil,err
}
	user, err := s.userRepo.FetchUser(userID)

	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *signupService) UpdatePhone(oldPhone string,newPhone string) (models.UserRole,string,error) {
	
	exists, err := s.userRepo.PhoneExist(newPhone)
	if err != nil {
logger.Log.Error("Error occur while changing phone",
zap.String("oldPhone",oldPhone),
zap.String("newPhone",newPhone),
)

		return "","",err
	}

	if exists {
		return "","",ErrPhoneAlreadyExists
	}

	// Find pending user using old phone
	pendingUser, err := s.pendingUserRepo.FindByPhone(oldPhone)

	if err != nil {
		return "","",err
	}
pendingUser.Phone=newPhone
	err = s.pendingUserRepo.Update(pendingUser)
if err != nil {
	logger.Log.Error(
        "Error updating pending user phone",
        zap.String("phone", pendingUser.Phone),
        zap.String("old phone", oldPhone),
        
    )
        return "","",err
    }
	// Generate OTP
	otpHash,err := s.otpService.SendOTP(pendingUser.Phone)
	    if err != nil {
		logger.Log.Error(
        "Error while sending otp to new number",
        zap.String("phone", pendingUser.Phone),
        
    )
        return pendingUser.Role,pendingUser.Phone,err
    }
	err = s.SaveOTPPendingUser(otpHash,pendingUser.Phone)
    if err != nil {
		logger.Log.Error(
        "Error while update phone number",
        zap.String("phone", pendingUser.Phone),
        
    )
        return pendingUser.Role,pendingUser.Phone,err
    }
    logger.Log.Info(
        "OTP resent successfully",
        zap.String("phone", pendingUser.Phone),
        
    )
return pendingUser.Role,pendingUser.Phone,nil

}

