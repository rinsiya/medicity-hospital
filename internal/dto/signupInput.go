package dto

type SignupInput struct {
	FirstName      string `form:"first_name" json:"first_name" binding:"required"`
	LastName       string `form:"last_name" json:"lastname" binding:"required"`
	Email          string `form:"email" json:"email" binding:"required,email"`
	Phone          string `form:"phone" json:"phone" binding:"required"`
	Password       string `form:"password" json:"password" binding:"required,min=6"`
	ConfirmPassword string `form:"confirm_password" json:"confirm_password" binding:"required,eqfield=Password"`
}

