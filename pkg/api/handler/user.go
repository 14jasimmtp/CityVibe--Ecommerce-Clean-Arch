package handlers

import (
	"fmt"
	"net/http"

	"github.com/14jasimmtp/CityVibe-Project-Clean-Architecture/pkg/models"
	interfaceUsecase "github.com/14jasimmtp/CityVibe-Project-Clean-Architecture/pkg/usecase/interface"
	"github.com/14jasimmtp/CityVibe-Project-Clean-Architecture/pkg/utils"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
)

type UserHandler struct {
	UserUsecase interfaceUsecase.UserUseCase
}

func NewUserHandler(usecase interfaceUsecase.UserUseCase) *UserHandler {
	return &UserHandler{UserUsecase: usecase}
}

// @Summary		User Signup
// @Description	user can signup by giving their details
// @Tags			User Login/Signup
// @Accept			json
// @Produce		    json
// @Param			signup  body  models.UserSignUpDetails  true	"signup"
// @Success		200	{object}	string "message":"successfully signed up.Enter otp to login"
// @Failure		500	{object}	string "error":err.Error()
// @Router			/signup    [POST]
func (clean *UserHandler) UserSignup(c *gin.Context) {

	var User models.UserSignUpDetails

	if c.ShouldBindJSON(&User) != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Enter Details in correct format"})
		return
	}
	data, err := utils.Validation(User)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": data})
	}

	err = clean.UserUsecase.SignUp(User)
	if err != nil {

		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Successfully signed up.Enter otp to login."})

}

// @Summary		User Login
// @Description	user can login by giving their phone and password
// @Tags			User Login/Signup
// @Accept			json
// @Produce		    json
// @Param			Login  body  models.UserLoginDetails  true	"signup"
// @Success		200	{object}	string "message":"Enter otp to login"
// @Failure		500	{object}	string "error":err.Error()
// @Router			/login    [POST]
func (clean *UserHandler) UserLogin(c *gin.Context) {
	var User models.UserLoginDetails

	if c.ShouldBindJSON(&User) != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Enter Details in correct format"})
		return
	}
	Error, err := utils.Validation(User)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": Error})
		return
	}

	user, Tokenstring, err := clean.UserUsecase.UserLogin(User)
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.SetCookie("Authorisation", Tokenstring, 3600, "", "", false, true)
	c.JSON(http.StatusOK, gin.H{"message": "user successfully logged in", "user": user})
}

// @Summary		Verify OTP
// @Description	user can login by giving the otp send to the mobile number
// @Tags			User Login/Signup
// @Accept			json
// @Produce		    json
// @Param			Verify  body  models.OTP  true	"Verify"
// @Success		200	{object}	string "message":"user successfully logged in" "user":models.UserLoginResponse
// @Failure		500	{object}	string "error":err.Error()
// @Router			/verify    [POST]
func (clean *UserHandler) VerifyLoginOtp(c *gin.Context) {
	var otp models.OTP

	if c.Bind(&otp) != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Enter otp to login"})
		return
	}

	err := utils.CheckOtp(otp.Phone, otp.Otp)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid otp"})
		return
	}

	user, err := clean.UserUsecase.FindUserByPhone(otp.Phone)
	if err != nil {
		return
	}

	Tokenstring, err := utils.TokenGenerate(user, "user")
	if err != nil {
		return
	}
	c.SetCookie("Authorisation", Tokenstring, 3600, "", "", false, true)
	var ResUser models.UserLoginResponse
	copier.Copy(&ResUser, &user)
	c.JSON(http.StatusOK, gin.H{"message": "user successfully logged in", "user": ResUser})
}

// @Summary		User Logout
// @Description	user can logout by sending this request to server
// @Tags			User Login/Signup
// @Produce		    json
// @Success		200	{object}	string "message":"user logged out successfully"
// @Failure		500	{object}	string "error":err.Error()
// @Router			/logout    [POST]
func (clean *UserHandler) UserLogout(c *gin.Context) {
	c.SetCookie("Authorisation", "", -1, "", "", false, true)

	c.JSON(http.StatusOK, gin.H{"message": "user logged out successfully"})
	fmt.Println("cookie deleted")
}

// @Summary		Forgot password
// @Description	user can will get otp to change password if forgotted
// @Tags			User Login/Signup
// @Produce		    json
// @Param			forgotPassword  body  models.Phone  true	"Forgot password"
// @Success		200	{object}	string "message":"user logged out successfully"
// @Failure		500	{object}	string "error":err.Error()
// @Router			/password/forgot    [POST]
func (clean *UserHandler) ForgotPassword(c *gin.Context) {
	var forgotPassword models.Phone
	if c.ShouldBindJSON(&forgotPassword) != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Enter constraints correctly"})
	}

	Error, err := utils.Validation(forgotPassword)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": Error})
		return
	}

	err = clean.UserUsecase.ForgotPassword(forgotPassword.Phone)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Enter otp and new password"})

}

// @Summary		Reset Password
// @Description user can reset password by entering otp and new password
// @Tags			User Login/Signup
// @Produce		    json
// @Param			ResetPassword  body  models.ForgotPassword  true	"Reset Password"
// @Success		200	{object}	string "message":"user logged out successfully"
// @Failure		500	{object}	string "error":err.Error()
// @Router			/password/forgot/change    [POST]
func (clean *UserHandler) ResetForgottenPassword(c *gin.Context) {
	var Newpassword models.ForgotPassword

	if c.ShouldBindJSON(&Newpassword) != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Enter details in correct format"})
		return
	}

	Error, err := utils.Validation(Newpassword)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": Error})
		return
	}

	err = clean.UserUsecase.ResetForgottenPassword(Newpassword)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
	}

	c.JSON(http.StatusOK, gin.H{"message": "password changed successfully"})
}

// @Summary		View Addresses
// @Description user can view addresses that he registered
// @Tags			User Profile
// @Produce		    json
// @Success		200	{object}	string "message":"user address","Address":Address
// @Failure		500	{object}	string "error":err.Error()
// @Router			/address    [Get]
func (clean *UserHandler) ViewUserAddress(c *gin.Context) {
	Token, err := c.Cookie("Authorisation")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	Address, err := clean.UserUsecase.ViewUserAddress(Token)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User Address", "Address": Address})

}

// @Summary Add new address details
// @Description Add new address details for the authenticated user.
// @Tags User Profile
// @Accept json
// @Produce json
// @Param address_details body models.Address true "New address details to be added"
// @Success 200 {object} string "message": "Address added successfully", "Address": AddressRes
// @Failure 400 {object} string "error": "Bad Request"
// @Failure 401 {object} string "error": "Unauthorized"
// @Failure 500 {object} string "error": "Internal Server Error"
// @Router /address [post]
func (clean *UserHandler) AddNewAddressDetails(c *gin.Context) {
	var Address models.Address

	if c.ShouldBindJSON(&Address) != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Enter Details correctly"})
	}

	Error, err := utils.Validation(Address)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": Error})
		return
	}

	Token, err := c.Cookie("Authorisation")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	AddressRes, err := clean.UserUsecase.AddAddress(Address, Token)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	c.JSON(http.StatusOK, gin.H{"message": "Address added successfully", "Address": AddressRes})
}

// @Summary Edit user address
// @Description Edit the address for a user.
// @Tags User Profile
// @Accept json
// @Produce json
// @Param id query string true "Address ID to be updated"
// @Param address body models.Address true "Updated address details"
// @Success 200 {object} string "message": "Address updated successfully", "Address": UpdatedAddress}
// @Failure 400 {object} string "error": "Enter constraints correctly"
// @Failure 401 {object} string "error": "Unauthorized"
// @Failure 500 {object} string "error": "Internal Server Error"
// @Router /address [put]
func (clean *UserHandler) EditUserAddress(c *gin.Context) {
	var UpdateAddress models.Address

	if c.ShouldBindJSON(&UpdateAddress) != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Enter constraints correctly"})
	}

	Error, err := utils.Validation(UpdateAddress)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": Error})
		return
	}

	Aid := c.Query("id")
	Token, err := c.Cookie("Authorisation")
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"error": err.Error()})
		return
	}

	UpdatedAddress, err := clean.UserUsecase.UpdateAddress(Token, Aid, UpdateAddress)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Address updated successfully", "Address": UpdatedAddress})

}

// @Summary Remove user address
// @Description Remove the address associated with a user.
// @Tags User Profile
// @Accept json
// @Produce json
// @Param id query string true "Address ID to be removed"
// @Success 200 {object} string "message": "Address removed successfully"
// @Failure 400 {object} string "error": "Bad Request"
// @Failure 401 {object} string "error": "Unauthorized"
// @Failure 500 {object} string "error": "Internal Server Error"
// @Router /address [delete]
func (clean *UserHandler) RemoveUserAddress(c *gin.Context) {
	Aid := c.Query("id")
	Token, err := c.Cookie("Authorisation")
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"error": err.Error()})
		return
	}

	err = clean.UserUsecase.DeleteAddress(Token, Aid)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Address removed successfully"})
}

// @Summary Get user profile
// @Description Retrieve the profile details of the authenticated user.
// @Tags User Profile
// @Accept json
// @Produce json
// @Success 200 {object} string "message": "User Profile", "profile": UserDetails
// @Failure 400 {object} string "error": "Bad Request"
// @Failure 401 {object} string "error": "Unauthorized"
// @Failure 500 {object} string "error": "Internal Server Error"
// @Router /profile [get]
func (clean *UserHandler) UserProfile(c *gin.Context) {
	Token, err := c.Cookie("Authorisation")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	UserDetails, err := clean.UserUsecase.UserProfile(Token)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User Profile", "profile": UserDetails})

}

// @Summary Update user profile
// @Description Update the profile details of the authenticated user.
// @Tags User Profile
// @Accept json
// @Produce json
// @Param user_details body models.UserProfile true "Updated user profile details"
// @Success 200 {object} string "message": "Updated User Profile", "profile": updatedUserDetails
// @Failure 400 {object} string "error": "Bad Request"
// @Failure 401 {object} string "error": "Unauthorized"
// @Failure 500 {object} string "error": "Internal Server Error"
// @Router /profile [put]
func (clean *UserHandler) UpdateUserProfile(c *gin.Context) {
	var UserDetails models.UserProfile

	if c.ShouldBindJSON(&UserDetails) != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Enter Details correctly"})
		return
	}

	Error, err := utils.Validation(UserDetails)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": Error})
		return
	}

	Token, err := c.Cookie("Authorisation")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updatedUserDetails, err := clean.UserUsecase.UpdateUserProfile(Token, UserDetails)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Updated User Profile", "profile": updatedUserDetails})

}
