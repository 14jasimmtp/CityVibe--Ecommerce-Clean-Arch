package interfaceUsecase

import (
	"github.com/14jasimmtp/CityVibe-Project-Clean-Architecture/pkg/domain"
	"github.com/14jasimmtp/CityVibe-Project-Clean-Architecture/pkg/models"
)

type UserUseCase interface {
	SignUp(User models.UserSignUpDetails) error
	UserLogin(user models.UserLoginDetails) error
	ForgotPassword(phone string) error
	ResetForgottenPassword(Newpassword models.ForgotPassword) error
	AddAddress(Address models.Address, Token string) (models.AddressRes, error)
	ViewUserAddress(Token string) ([]models.AddressRes, error)
	UpdateAddress(Token, aid string, NewAddress models.Address) (models.AddressRes, error)
	DeleteAddress(Token, aid string) error
	UserProfile(Token string) (models.UserProfile, error)
	UpdateUserProfile(Token string, User models.UserProfile) (models.UserProfile, error)
	FindUserByPhone(phone string) (*domain.User, error)
}
