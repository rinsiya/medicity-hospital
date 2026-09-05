package database

import (
	"medicity/logger"
	"medicity/internal/models"

	"go.uber.org/zap"
)

func MigrateDatabase() error {
		logger.Log.Info("Database migration initialized")

	err := DB.AutoMigrate(

		&models.User{},
		&models.Doctor{},

		&models.Department{},
		
		&models.Patient{},

		&models.Address{},

		&models.File{},
		&models.DoctorProfile{},
		&models.DoctorQualification{},
		&models.DoctorTimeSlot{},
		&models.BankAccount{},
		&models.Wallet{},
        &models.PendingUserSignup{},
		&models.Appointment{},
		&models.CancelledAppointment{},
		&models.Prescription{},
		&models.VitalData{},
		&models.Message{},
		&models.DoctorReview{},
		&models.Payment{},
		&models.Refund{},
		&models.Withdrawal{}, 
	&models.Notification{},
	&models.PasswordReset{},
	)

	if err != nil {
		logger.Log.Fatal("Migration failed", zap.Error(err))
		return err
	}

	logger.Log.Info("Database migrated successfully")

	return nil
}
