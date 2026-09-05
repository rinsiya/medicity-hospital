package routes

import (
	"medicity/internal/handler"
	"medicity/internal/middleware"
	"medicity/internal/repository"
	"medicity/internal/service"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.Engine) {
	// Repository
	userRepo := repository.NewUserRepository()
	pendingUserRepo := repository.NewPendingUserSignupRepository()
	patientRepo := repository.NewPatientRepository()
	doctorRepo := repository.NewDoctorRepository()
	// Service
	signupService := service.NewSignupService(pendingUserRepo, userRepo, patientRepo, doctorRepo)
	passwordResetRepo := repository.NewPasswordResetRepository()
	userService := service.NewUserService(userRepo, signupService, passwordResetRepo)

	patientService := service.NewPatientService(patientRepo)
	doctorService := service.NewDoctorService(doctorRepo)
vitalService := service.NewVitalService(repository.NewVitalRepository())
	// Handler
	patientHandler := handler.NewPatientHandler(patientService,userService)
	doctorHandler := handler.NewDoctorHandler(doctorService)

	userHandler := handler.NewUserHandler(userService)
	signupHandler := handler.NewSignupHandler(signupService)

	vitalDataHandler := handler.NewVitalHandler(patientService, vitalService)
	appointmentHandler := handler.NewAppointmentHandler(service.NewAppointmentService(repository.NewAppointmentRepository()), patientService)
    fileHandler := handler.NewFileHandler(service.NewFileService(repository.NewFileRepository()))
	prescriptionHandler := handler.NewPrescriptionHandler(service.NewPrescriptionService(repository.NewPrescriptionRepository()))
//addressHandler := handler.NewAddressHandler(addressService)


	router.GET("/", userHandler.Home)
	router.GET("/home", userHandler.Home)
	router.POST("/:role/login", userHandler.Login)
	router.POST("/signup/:role", signupHandler.Signup)

	router.GET("/:role/verify-otp", signupHandler.ShowOTPPage)
	router.POST("/verify-otp", signupHandler.ValidateOTP)

	router.GET("/:role/verify-otp-pending", signupHandler.ShowOTPPendingPage)
	router.POST("/resend-otp", signupHandler.ResendOTP)
	router.GET("/change-phone", signupHandler.ChangePhone)
	router.POST("/change-phone", signupHandler.UpdatePhone)

	router.GET("/patient/verification-success", signupHandler.PatientVerificationSuccess)
	router.GET("/doctor/verification-success", signupHandler.DoctorVerificationSuccess)
	router.GET("/patient/login", userHandler.PatientLogin)
	router.GET("/doctor/login", userHandler.DoctorLogin)
	router.GET("/admin/login", userHandler.AdminLogin)

	router.GET("/patient/signup", patientHandler.PatientSignup)
	router.GET("/doctor/signup", doctorHandler.DoctorSignup)
	router.POST("/patient/logout", patientHandler.Logout)
	//forgot password routes
	router.GET("/:role/forgot-password", userHandler.ForgotPasswordForm)
     router.POST("/:role/auth/forgot-password", userHandler.ForgotPasswordPhoneVerification)
     router.GET("/:role/forgot-password/otp-verification", userHandler.ForgotPasswordOTPVerification)
	router.POST("/:role/verify-password-reset-otp", userHandler.ForgotPasswordOTPVerificationRequest)
	router.POST("/:role/change-forgot-password", userHandler.UpdateForgotPassword)
	
	router.POST("/forgot-password-resend-otp", userHandler.ResendOTP)

	//router.POST("/patient/change-password", patientHandler.UpdatePassword)
	//router.GET("/patient/change-password", patientHandler.ChangePassword)

	patient := router.Group("/patient")
	patient.Use(middleware.JWTAuth("patient"))
	patient.Use(middleware.RequireRole("patient"))
	{
		patient.GET("/home", patientHandler.Home)
		patient.GET("/vitals", vitalDataHandler.GetPatientVitals)
		patient.POST("/vitals", vitalDataHandler.SavePatientVitals)
		patient.GET("/appointments/recent", appointmentHandler.GetRecentAppointments)
		patient.GET("/files/recent", fileHandler.GetRecentMedicalFiles)
		patient.GET("/prescriptions/recent", prescriptionHandler.GetRecentPrescriptions)
		patient.GET("/profile/complete", patientHandler.CheckProfileComplete)
		patient.GET("/complete-profile", patientHandler.CompleteProfile)
		patient.POST("/complete-profile", patientHandler.CompleteProfileRequest)
	//	patient.GET("/complete-address", addressHandler.CompleteAddress)
	//	patient.POST("/complete-address", addressHandler.CompleteAddressRequest)

	}

	doctor := router.Group("/doctor")
	doctor.Use(middleware.JWTAuth("doctor"))
	doctor.Use(middleware.RequireRole("doctor"))
	{
		doctor.GET("/dashboard", doctorHandler.Dashboard)
		doctor.GET("/change-password", doctorHandler.ChangePassword)
		doctor.POST("/change-password", doctorHandler.UpdatePassword)

	}

}
