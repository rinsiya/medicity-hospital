package service

import (
	"errors"
	"time"

	"medicity/internal/dto"
	"medicity/internal/models"
	"medicity/internal/repository"
	"medicity/logger"
	"medicity/pkg/utils"

	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrUserNotFound        = errors.New("user not found")
	ErrFailedPasswordReset = errors.New("Password reset failed")
)

type UserService interface {
	Login(input dto.LoginInput, role string) (*models.User, error)
	UserExistByPhoneAndRole(phone, role string) error
	//ForgotPasswordOTPVerification(phone string) error
	ResendOTP(phone string) error
	FindPasswordResetUserByPhone(phone string) (*models.PasswordReset, error)
	SaveOTPPasswordResetUser(otpHash, phone string) error
	ValidatePasswordResetOTP(phone string, otp string) (time.Time, error)
	VerifyPasswordResetOTP(phone string, otp string) (time.Time, error)
	UpdatePassword(phone string, role string, newPassword string) error
	GetPatientIDByUserID(userID uint) uint
	GetUserByID(userID uint)(*models.User,error)
}

// signupService := &SignupService{}
type userService struct {
	repo              repository.UserRepository
	passwordResetRepo repository.PasswordResetRepository
	signupService     SignupService
}

func NewUserService(r repository.UserRepository, signupService SignupService, passwordResetRepo repository.PasswordResetRepository) UserService {
	return &userService{repo: r, signupService: signupService, passwordResetRepo: passwordResetRepo}
}

//	func (s *userService) ForgotPasswordOTPVerification(phone string) error {
//		return nil
//	}

func (s  *userService) GetUserByID(userID uint)(*models.User,error){

	return s.repo.GetUserByID(userID)

}
func (s *userService) GetPatientIDByUserID(userID uint) uint {
	patientID, err := s.repo.GetPatientIDByUserID(userID)
	if err != nil {
		logger.Log.Error("Failed to find patient ID by user ID", zap.Uint("userID", userID), zap.Error(err))
		return 0
	}
	return patientID
}	

func (s *userService) UpdatePassword(phone string, role string, newPassword string) error {
user,err:=s.repo.FindByPhone(phone,role)

if err != nil{
	logger.Log.Error("Failed to find user by phone and role",zap.String("phone", phone),
zap.String("role", role),
zap.Error(err),
)
return err
}
if user==nil{
	logger.Log.Warn("User not found for password update",zap.String("phone", phone),
zap.String("role", role),
)
return ErrUserNotFound	
}
	hashedPassword, err := utils.HashPassword(newPassword)
if err != nil {
	logger.Log.Error("Failed to hash password", zap.String("phone", phone), zap.String("role", role), zap.Error(err))
	return err
}

user.Password = hashedPassword

err = s.repo.Update(user)
if err != nil {
	logger.Log.Error("Failed to update password", zap.String("phone", phone), zap.String("role", role), zap.Error(err))
	return err
}

logger.Log.Info("Password updated successfully", zap.String("phone", phone), zap.String("role", role))
return nil

}

// func (s *userService) UpdatePassword(phone string, newPassword string) error {}
func (s *userService) VerifyPasswordResetOTP(phone string, otp string) (time.Time, error) {

	// Find pending user
	passwordResetUser, err := s.passwordResetRepo.FindByPhone(phone)

	if err != nil {
		return time.Time{}, err
	}

	if passwordResetUser == nil {
		return time.Time{}, errors.New(
			"password reset user not found",
		)
	}

	// Validate OTP
	otpExpiresAt, err := s.ValidatePasswordResetOTP(phone, otp)

	if err != nil {
		return otpExpiresAt, err
	}
	return otpExpiresAt, nil
}

func (s *userService) ValidatePasswordResetOTP(phone string, otp string) (time.Time, error) {

	passwordResetUser, err := s.passwordResetRepo.FindByPhone(phone)

	if err != nil {
		logger.Log.Error(
			"Failed to find user by phone",
			zap.String("phone", phone),
			zap.Error(err),
		)

		return time.Time{}, err
	}

	if passwordResetUser == nil {
		return time.Time{}, errors.New("password reset user not found")
	}

	// Keep the expiry time so the handler can send it back to UI
	expiresAt := passwordResetUser.OTPExpiresAt

	// Check expiry
	if time.Now().After(expiresAt) {

		logger.Log.Warn(
			"OTP expired",
			zap.String("phone", phone),
		)

		return expiresAt, errors.New("OTP expired")
	}

	// Compare OTP
	err = bcrypt.CompareHashAndPassword(
		[]byte(passwordResetUser.OTPHash),
		[]byte(otp),
	)

	if err != nil {

		logger.Log.Warn(
			"Invalid OTP",
			zap.String("phone", phone),
		)

		return expiresAt, errors.New("Invalid OTP")
	}
	logger.Log.Info(
		"OTP verified successfully",
		zap.String("phone", phone),
	)
		logger.Log.Info(
		"deleting entry in password reset table after successful verification",
		zap.String("phone", phone),
	)
		 err = s.passwordResetRepo.Delete(passwordResetUser.PasswordResetID);
		 if err != nil {
				logger.Log.Info(
		"deleting entry in password reset table failed",
		zap.String("phone", phone),
	)
	return expiresAt, errors.New("deleting entry in password reset table failed")
		 }

	logger.Log.Info(
		"deleted entry in password reset table",
		zap.String("phone", phone),
	)

	return expiresAt, nil
}

func (s *userService) SaveOTPPasswordResetUser(otpHash, phone string) error {
	passwordResetUser, err := s.passwordResetRepo.FindByPhone(phone)

	if err != nil {

		logger.Log.Error(
			"Failed to find user by phone",
			zap.String("phone", phone),
			zap.Error(err),
		)

		return err
	}
	passwordResetUser.OTPHash = otpHash
	passwordResetUser.OTPExpiresAt = time.Now().Add(2 * time.Minute)
	logger.Log.Info("Password reset update otp hash and expiry",
		zap.String("phone", phone))
	err = s.passwordResetRepo.Update(passwordResetUser)
	if err != nil {
		logger.Log.Error(
			"Failed to save otp to user",
			zap.String("phone", phone),
			zap.Error(err),
		)
		return err
	}
	logger.Log.Info("Password reset pending OTP verification",
		zap.String("phone", phone),
	//zap.String("role", string(passwordResetUser)),
	)

	return nil

}

func (s *userService) ResendOTP(phone string) error {
	otpHash, err := s.signupService.GenerateAndSendOTP(phone)
	s.SaveOTPPasswordResetUser(otpHash, phone)
	if err != nil {
		return err
	}
	logger.Log.Info(
		"OTP resent successfully",
		zap.String("phone", phone),
	)

	return nil
}
func (s *userService) FindPasswordResetUserByPhone(phone string) (*models.PasswordReset, error) {
	passwordResetUser, err := s.passwordResetRepo.FindByPhone(phone)
	if err != nil {

		logger.Log.Error(
			"Failed to find password reset user by phone",
			zap.String("phone", phone),
			zap.Error(err),
		)

		return nil, err
	}
	return passwordResetUser, err

}
func (s *userService) UserExistByPhoneAndRole(phone, role string) error {
	logger.Log.Info("Forgot password request",
		zap.String("phone", phone),
		zap.String("role", role),
	)

	user, err := s.repo.FindByPhone(phone, role)
	if err != nil {
		logger.Log.Error("Failed to find user for forgot password request",
			zap.String("phone", phone),
			zap.String("role", role),
			zap.Error(err),
		)
		return err
	}
	if user == nil {
		logger.Log.Warn("Forgot password request: user not found",
			zap.String("phone", phone),
			zap.String("role", role),
		)
		return ErrUserNotFound
	}

	logger.Log.Info("Forgot password request processed",
		zap.String("phone", phone),
		//zap.String("role", role),
	)

	otpHash, err := s.signupService.GenerateAndSendOTP(phone)
	if err != nil {

		logger.Log.Error(
			"Failed to create password reset OTP",
			zap.String("phone", phone),
			zap.Error(err),
		)

		return err
	}
	passwordResetUser := &models.PasswordReset{
		Phone:        phone,
		OTPHash:      otpHash,
		OTPExpiresAt: time.Now().Add(2 * time.Minute),
		Verified:     false,
	}

	err = s.passwordResetRepo.Create(passwordResetUser)
	if err != nil {
		logger.Log.Error(
			"Failed to create password reset record",
			zap.String("phone", phone),
			zap.Error(err),
		)

		return err
	}

	logger.Log.Info(
		"password reset OTP created successfully",
		zap.String("phone", phone),
	)

	return nil
}

func (s *userService) Login(input dto.LoginInput, role string) (*models.User, error) {

	logger.Log.Info("Login attempt",
		zap.String("username", input.Username),
		zap.String("requested_role", role),
	)

	// Find user by email or phone
	user, err := s.repo.FindByUsername(input.Username)

	if err != nil {
		logger.Log.Error("Failed to find user during login",
			zap.String("username", input.Username),
			zap.Error(err),
		)
		return nil, err
	}

	// User not found
	if user == nil {
		logger.Log.Warn("Login failed: user not found",
			zap.String("username", input.Username),
		)
		return nil, ErrInvalidCredentials
	}

	// Check password
	if !utils.CheckPassword(user.Password, input.Password) {
		logger.Log.Warn("Login failed: invalid password",
			zap.String("username", input.Username),
		)
		return nil, ErrInvalidCredentials
	}

	// Check requested role
	if string(user.Role) != role {
		logger.Log.Warn("Login failed: role mismatch",
			zap.String("username", input.Username),
			zap.String("user_role", string(user.Role)),
			zap.String("requested_role", role),
		)
		return nil, ErrInvalidCredentials
	}

	// Check account status
	if user.Status != models.UserActive {
		logger.Log.Warn("Login failed: account is not active",
			zap.String("username", input.Username),
			zap.String("status", string(user.Status)),
		)
		return nil, errors.New("user account is not active")
	}

	logger.Log.Info("Login successful",
		zap.Int("user_id", int(user.UserID)),
		zap.String("role", string(user.Role)),
	)

	return user, nil
}
