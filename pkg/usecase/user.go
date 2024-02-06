package usecase

import (
	"errors"
	"fmt"

	"github.com/14jasimmtp/CityVibe-Project-Clean-Architecture/pkg/db"
	"github.com/14jasimmtp/CityVibe-Project-Clean-Architecture/pkg/domain"
	"github.com/14jasimmtp/CityVibe-Project-Clean-Architecture/pkg/models"
	interfaceRepo "github.com/14jasimmtp/CityVibe-Project-Clean-Architecture/pkg/repository/interface"
	interfaceUsecase "github.com/14jasimmtp/CityVibe-Project-Clean-Architecture/pkg/usecase/interface"
	"github.com/14jasimmtp/CityVibe-Project-Clean-Architecture/pkg/utils"
	"github.com/jinzhu/copier"
	"golang.org/x/crypto/bcrypt"
)

type UserUseCase struct {
	UserRepo interfaceRepo.UserRepo
}

func NewUserUsecase(repo interfaceRepo.UserRepo) interfaceUsecase.UserUseCase{
	return &UserUseCase{UserRepo: repo}
}

func (clean *UserUseCase) SignUp(User models.UserSignUpDetails) error {
	CheckEmail, err := clean.UserRepo.CheckUserExistsEmail(User.Email)
	if err != nil {
		fmt.Println("server error")
		return errors.New("server error")
	}
	if CheckEmail != nil {
		fmt.Println("user already exist")
		return errors.New("user already exist with this email")
	}

	CheckPhone, err := clean.UserRepo.CheckUserExistsByPhone(User.Phone)
	if err != nil {
		fmt.Println("server error")
		return errors.New("server error")
	}
	if CheckPhone != nil {
		fmt.Println("user already exist with this number")
		return errors.New("user already exist with this number")
	}

	if User.Password != User.ConfirmPassword {
		fmt.Println("passwords doesn't match")
		return errors.New("paswords doesn't match")
	}

	HashedPassword, err := bcrypt.GenerateFromPassword([]byte(User.Password), 10)
	if err != nil {
		fmt.Println("error while hashing ")
		return errors.New("server error occured(password hashing)")
	}
	User.Password = string(HashedPassword)

	sentOtp := utils.SendOtp(User.Phone)
	if sentOtp != nil {
		fmt.Println("error gen otp")
		return errors.New("error occured generating otp")
	}
	var Userdt domain.User
	err = copier.Copy(&Userdt, &User)
	if err != nil {
		return err
	}
	db.DB.Create(&Userdt)
	return nil
}

func (clean *UserUseCase) UserLogin(user models.UserLoginDetails) (*models.UserLoginResponse,string,error) {
	CheckPhone, err := clean.UserRepo.CheckUserExistsByPhone(user.Phone)
	if err != nil {
		return nil,"",errors.New("error with server")
	}
	if CheckPhone == nil {
		return nil,"",errors.New("phone number doesn't exist")
	}
	userdetails, err := clean.UserRepo.FindUserByPhone(user.Phone)
	fmt.Println(userdetails, user.Password)
	if err != nil {
		return nil,"",err
	}

	if userdetails.Blocked {
		return nil,"",errors.New("user is blocked")
	}

	err = bcrypt.CompareHashAndPassword([]byte(userdetails.Password), []byte(user.Password))
	if err != nil {
		return nil,"",errors.New("password not matching")
	}
	Tokenstring, err := utils.TokenGenerate(userdetails, "user")
	if err != nil {
		fmt.Println(err)
		return nil,"",errors.New("error generating token")
	}
	var ResUser models.UserLoginResponse
	copier.Copy(&ResUser, &userdetails)

	return &ResUser,Tokenstring,nil
}

func (clean *UserUseCase) ForgotPassword(phone string) error {
	user, err := clean.UserRepo.FindUserByPhone(phone)
	if err != nil {
		return errors.New("user doesn't found with this number")
	}

	if user.Blocked {
		return errors.New("user is blocked")
	}
	err = utils.SendOtp(phone)
	if err != nil {
		return errors.New("error generating otp ")
	}

	return nil
}

func (clean *UserUseCase) ResetForgottenPassword(Newpassword models.ForgotPassword) error {
	err := utils.CheckOtp(Newpassword.Phone, Newpassword.OTP)
	if err != nil {
		return err
	}
	hashed, err := bcrypt.GenerateFromPassword([]byte(Newpassword.NewPassword), 10)
	if err != nil {
		return errors.New("error while hashing password")
	}
	Newpassword.NewPassword = string(hashed)

	err = clean.UserRepo.ChangePassword(Newpassword)
	if err != nil {
		return err
	}

	return nil
}

func (clean *UserUseCase) FindUserByPhone(phone string) (*domain.User, error){
	return clean.UserRepo.FindUserByPhone(phone)
}

func (clean *UserUseCase) AddAddress(Address models.Address, Token string) (models.AddressRes, error) {
	fmt.Println(Token)
	UserId, err := utils.ExtractUserIdFromToken(Token)
	if err != nil {
		return models.AddressRes{}, err
	}

	AddressRes, err := clean.UserRepo.AddAddress(Address, UserId)
	if err != nil {
		return models.AddressRes{}, err
	}

	return AddressRes, nil
}

func (clean *UserUseCase) ViewUserAddress(Token string) ([]models.AddressRes, error) {
	UserId, err := utils.ExtractUserIdFromToken(Token)
	if err != nil {
		return []models.AddressRes{}, err
	}

	Address, err := clean.UserRepo.ViewAddress(UserId)
	if err != nil {
		return []models.AddressRes{}, err
	}

	return Address, nil
}

func (clean *UserUseCase) UpdateAddress(Token, aid string, NewAddress models.Address) (models.AddressRes, error) {
	UserId, err := utils.ExtractUserIdFromToken(Token)
	if err != nil {
		return models.AddressRes{}, err
	}

	UpdatedAddress, err := clean.UserRepo.UpdateAddress(UserId, aid, NewAddress)
	if err != nil {
		return models.AddressRes{}, err
	}

	return UpdatedAddress, nil

}

func (clean *UserUseCase) DeleteAddress(Token, aid string) error {
	UserId, err := utils.ExtractUserIdFromToken(Token)
	if err != nil {
		return err
	}

	err = clean.UserRepo.RemoveAddress(UserId, aid)
	if err != nil {
		return err
	}

	return nil
}

func (clean *UserUseCase) UserProfile(Token string) (models.UserProfile, error) {
	UserId, err := utils.ExtractUserIdFromToken(Token)
	if err != nil {
		return models.UserProfile{}, err
	}

	User, err := clean.UserRepo.UserProfile(UserId)
	if err != nil {
		return models.UserProfile{}, err
	}

	return User, nil
}

func (clean *UserUseCase) UpdateUserProfile(Token string, User models.UserProfile) (models.UserProfile, error) {
	UserId, err := utils.ExtractUserIdFromToken(Token)
	if err != nil {
		return models.UserProfile{}, err
	}

	UpdatedUser, err := clean.UserRepo.UpdateUserProfile(UserId, User)
	if err != nil {
		return models.UserProfile{}, err
	}

	return UpdatedUser, nil
}
