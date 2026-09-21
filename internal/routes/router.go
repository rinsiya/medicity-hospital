package routes

import (
	"log"
	"medicity/internal/handler"
	"medicity/internal/middleware"
	"medicity/internal/repository"
	"medicity/internal/service"
	"medicity/pkg/storage"

	"net/http"
	"os"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func SetupRoutes(router *gin.Engine) {
	// Repository
	userRepo := repository.NewUserRepository()
	pendingUserRepo := repository.NewPendingUserSignupRepository()
	patientRepo := repository.NewPatientRepository()

	passwordResetRepo := repository.NewPasswordResetRepository()
	fileRepo := repository.NewFileRepository()
	
	doctorRepo := repository.NewDoctorRepository()
addressRepo := repository.NewAddressRepository()
	// Service

	err := godotenv.Load()
if err != nil {
    log.Println("Warning: .env file not found")
}

apiKey := os.Getenv("TWO_FACTOR_API_KEY")
otpService :=service.NewOTPService(apiKey)
cld, err := cloudinary.New()
if err != nil {
	log.Fatal("Failed to initialize Cloudinary:", err)
}

cld.Config.URL.Secure = true

cloudinaryStorage := storage.NewCloudinaryStorage(cld)    
	signupService := service.NewSignupService(otpService,pendingUserRepo, userRepo, patientRepo, doctorRepo)
	userService := service.NewUserService(userRepo, otpService,signupService, passwordResetRepo)
addressService := service.NewAddressService(addressRepo)
	patientService := service.NewPatientService(patientRepo,fileRepo)
	doctorService := service.NewDoctorService(doctorRepo)
	fileService := service.NewFileService(
	fileRepo,
	cloudinaryStorage,
)
vitalService := service.NewVitalService(repository.NewVitalRepository())
	// Handler
	patientHandler := handler.NewPatientHandler(patientService,userService,fileService)
	doctorHandler := handler.NewDoctorHandler(doctorService)



	userHandler := handler.NewUserHandler(userService,otpService,doctorService)
	signupHandler := handler.NewSignupHandler(signupService, otpService)
	vitalDataHandler := handler.NewVitalHandler(patientService, vitalService)
	appointmentHandler := handler.NewAppointmentHandler(service.NewAppointmentService(repository.NewAppointmentRepository()), patientService)
    fileHandler := handler.NewFileHandler(fileService)
	prescriptionHandler := handler.NewPrescriptionHandler(service.NewPrescriptionService(repository.NewPrescriptionRepository()))
//addressHandler := handler.NewAddressHandler(addressService)
addressHandler := handler.NewAddressHandler(addressService,patientService)

	router.GET("/", userHandler.Home)
	router.GET("/home", userHandler.Home)
	router.POST("/:role/login", userHandler.Login)
	router.POST("/signup/:role", signupHandler.Signup)

	router.GET("/:role/verify-otp", signupHandler.ShowOTPPage)
	router.POST("/:role/verify-otp", signupHandler.ValidateOTP)
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
     router.POST("/:role/auth/forgot-password", userHandler.SendForgotPasswordOTP)
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

	patient.GET("/appointments/all",appointmentHandler.GetPatientAppointments)
patient.GET("/appointments", func(c *gin.Context) {
    c.HTML(http.StatusOK, "appointments.html", nil)
})
	//patient.GET("/files",fileHandler.GetUserFiles)
	patient.GET("/prescriptions",prescriptionHandler.GetPatientPrescriptions)

patient.POST("/complete-address", addressHandler.AddAddress)
//patient.POST("/medical-reports",fileHandler.UploadMedicalFiles)
patient.POST("/upload-file",fileHandler.UploadFile)
//patient.GET("/miniProfile", patientHandler.GetMiniProfile)
	}

	doctor := router.Group("/doctor")
	doctor.Use(middleware.JWTAuth("doctor"))
	doctor.Use(middleware.RequireRole("doctor"))
	{
		doctor.GET("/dashboard", doctorHandler.Dashboard)
		//doctor.GET("/change-password", doctorHandler.ChangePassword)
		//doctor.POST("/change-password", doctorHandler.UpdatePassword)
		doctor.GET("/complete-profile", doctorHandler.CompleteProfile)
		doctor.GET("/verification-penidng",doctorHandler.VerificationPending)
		doctor.GET("/profile-rejected",doctorHandler.ProfileRejected)

	}

}
