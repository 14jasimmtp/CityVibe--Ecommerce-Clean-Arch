package interfaceRepo

import (
	"github.com/14jasimmtp/CityVibe-Project-Clean-Architecture/pkg/domain"
	"github.com/14jasimmtp/CityVibe-Project-Clean-Architecture/pkg/models"
)

type UserRepo interface {
	CheckUserExistsEmail(email string) (*domain.User, error)
	CheckUserExistsByPhone(phone string) (*domain.User, error)
	SignUpUser(user models.UserSignUpDetails) (*models.UserDetailsResponse, error)
	FindUserByPhone(phone string) (*domain.User, error)
	GetUserById(id int) (*models.UserDetailsResponse, error)
	ChangePassword(ResetUser models.ForgotPassword) error
	AddAddress(Address models.Address, UserId uint) (models.AddressRes, error)
	UpdateAddress(userid uint, aid string, Address models.Address) (models.AddressRes, error)
	ViewAddress(id uint) ([]models.AddressRes, error)
	RemoveAddress(Userid uint, aid string) error
	UserProfile(userid uint) (models.UserProfile, error)
	UpdateUserProfile(userid uint, user models.UserProfile) (models.UserProfile, error)
	CheckAddressExist(userid uint, address uint) bool
	UpdateWallet(wallet float64, userID uint) error
}
