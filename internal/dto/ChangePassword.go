package dto
type ChangePasswordInput struct{
		Password       string `form:"password" json:"password" binding:"required,min=6"`
	ConfirmPassword string `form:"confirm_password" json:"confirm_password" binding:"required,eqfield=Password"`
}

type ForgotPasswordPhone struct{
	Phone string `form:"Phone" json:"Phone" binding:"required"`}