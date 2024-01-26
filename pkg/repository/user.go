package repository

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/14jasimmtp/CityVibe-Project-Clean-Architecture/pkg/db"
	"github.com/14jasimmtp/CityVibe-Project-Clean-Architecture/pkg/domain"
	"github.com/14jasimmtp/CityVibe-Project-Clean-Architecture/pkg/models"
	interfaceRepo "github.com/14jasimmtp/CityVibe-Project-Clean-Architecture/pkg/repository/interface"
	"gorm.io/gorm"
)

type UserRepo struct {
	DB *gorm.DB
}

func NewUserRepo (db *gorm.DB) interfaceRepo.UserRepo{
	return &UserRepo{DB: db}
}

func (clean *UserRepo) CheckUserExistsEmail(email string) (*domain.User, error) {
	var user domain.User
	result := db.DB.Where(&domain.User{Email: email}).First(&user)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, result.Error
	}
	return &user, nil

}

func (clean *UserRepo) CheckUserExistsByPhone(phone string) (*domain.User, error) {
	var user domain.User
	result := db.DB.Where(&domain.User{Phone: phone}).First(&user)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, result.Error
	}
	return &user, nil
}

func (clean *UserRepo) SignUpUser(user models.UserSignUpDetails) (*models.UserDetailsResponse, error) {
	var User models.UserDetailsResponse

	result := db.DB.Raw("INSERT INTO users(firstname,lastname,email,phone,password) VALUES(?,?,?,?,?)", user.FirstName, user.LastName, user.Email, user.Phone, user.Password).Scan(&User)
	if result.Error != nil {
		return nil, result.Error
	}
	return &User, nil
}

func (clean *UserRepo) FindUserByPhone(phone string) (*domain.User, error) {
	var user domain.User
	result := db.DB.Raw("SELECT * FROM users WHERE phone = ?", phone).Scan(&user)
	if result.Error != nil {
		return &domain.User{}, result.Error
	}
	return &user, nil
}

func (clean *UserRepo) GetUserById(id int) (models.UserDetailsResponse, error) {
	var user models.UserDetailsResponse

	result := db.DB.Raw("SELECT * FROM users WHERE id = ? ", id).Scan(&user)
	if result.Error != nil {
		fmt.Println("error fetching user")
		return models.UserDetailsResponse{}, result.Error
	}
	if result.RowsAffected == 0 {
		return models.UserDetailsResponse{}, errors.New(`no users found with this id`)
	}
	return user, nil
}

func (clean *UserRepo) ChangePassword(ResetUser models.ForgotPassword) error {
	query := db.DB.Exec(`UPDATE users SET password = ? WHERE phone = ?`, ResetUser.NewPassword, ResetUser.Phone)
	if query.Error != nil {
		return query.Error
	}
	return nil
}

func (clean *UserRepo) AddAddress(Address models.Address, UserId uint) (models.AddressRes, error) {
	var AddressRes models.AddressRes
	query := db.DB.Raw(`INSERT INTO addresses(user_id,name,house_name,phone,street,city,state,pin) VALUES (?,?,?,?,?,?,?,?) RETURNING id,name,house_name,phone,street,city,state,pin`, UserId, Address.Name, Address.Housename, Address.Phone, Address.Street, Address.City, Address.State, Address.Pin).Scan(&AddressRes)
	if query.Error != nil {
		return models.AddressRes{}, query.Error
	}
	return AddressRes, nil
}

func (clean *UserRepo) UpdateAddress(userid uint, aid string, Address models.Address) (models.AddressRes, error) {
	var AddressRes models.AddressRes
	query := db.DB.Raw(`UPDATE addresses SET name = ?,phone = ?,house_name = ?,street = ?,city = ?,state = ?,pin=? WHERE id = ? AND user_id = ? RETURNING id,name,phone,house_name AS housename,street,city,state,pin`, Address.Name, Address.Phone, Address.Housename, Address.Street, Address.City, Address.State, Address.Pin, aid, userid).Scan(&AddressRes)
	if query.Error != nil {
		return models.AddressRes{}, query.Error
	}
	if query.RowsAffected == 0 {
		return models.AddressRes{}, errors.New(`no address found to update with this id`)
	}
	return AddressRes, nil
}

func (clean *UserRepo) ViewAddress(id uint) ([]models.AddressRes, error) {
	var Address []models.AddressRes
	query := db.DB.Raw(`SELECT * FROM addresses WHERE user_id = ?`, id).Scan(&Address)
	if query.Error != nil {
		return []models.AddressRes{}, query.Error
	}

	if query.RowsAffected < 1 {
		return []models.AddressRes{}, errors.New("no address found. add new address")
	}

	return Address, nil
}

func (clean *UserRepo) RemoveAddress(Userid uint, aid string) error {
	query := db.DB.Exec(`DELETE FROM addresses WHERE id = ? && user_id = ?`, aid, Userid)
	if query.Error != nil {
		return query.Error
	}

	return nil
}

func (clean *UserRepo) UserProfile(userid uint) (models.UserProfile, error) {
	var err error
	var User models.UserProfile
	query := db.DB.Raw(`SELECT * FROM users WHERE id = ?`, userid).Scan(&User)
	if query.Error != nil {
		return models.UserProfile{}, query.Error
	}

	if query.RowsAffected < 1 {
		return models.UserProfile{}, errors.New(`no user profile found`)
	}
	formattedValue := fmt.Sprintf("%.2f", User.Wallet)
	User.Wallet, err = strconv.ParseFloat(formattedValue, 64)
	if err != nil {
		fmt.Println("Error:", err)
		return models.UserProfile{}, err
	}
	return User, nil
}

func (clean *UserRepo) UpdateUserProfile(userid uint, user models.UserProfile) (models.UserProfile, error) {
	var UpdatedUser models.UserProfile
	query := db.DB.Raw(`UPDATE users SET firstname = ?,lastname = ?,email = ?,phone = ? WHERE id = ? RETURNING firstname,lastname,email,phone`, user.Firstname, user.Lastname, user.Email, user.Phone, userid).Scan(&UpdatedUser)
	if query.Error != nil {
		return models.UserProfile{}, query.Error
	}

	return UpdatedUser, nil
}

func (clean *UserRepo) CheckAddressExist(userid uint, address uint) bool {
	var count int
	if err := db.DB.Raw("SELECT COUNT(*) FROM addresses WHERE id = ? AND user_id = ?", address, userid).Scan(&count).Error; err != nil {
		return false
	}
	return count > 0
}

func (clean *UserRepo) UpdateWallet(wallet float64, userID uint) error {
	query := db.DB.Exec(`UPDATE users SET wallet = ? WHERE id = ?`, wallet, userID)
	if query.Error != nil {
		return errors.New(`something went wrong`)
	}
	return nil
}
