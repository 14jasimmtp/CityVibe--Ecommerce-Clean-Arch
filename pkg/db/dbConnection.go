package db

import (
	"fmt"

	"github.com/14jasimmtp/CityVibe-Project-Clean-Architecture/pkg/config"
	"github.com/14jasimmtp/CityVibe-Project-Clean-Architecture/pkg/domain"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func DBInitialise(cfg config.Config) *gorm.DB {
	var err error

	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Shanghai",
		cfg.DBHost,
		cfg.DBUser,
		cfg.DBPassword,
		cfg.DBName,
		cfg.DBPort,
	)
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		fmt.Println(dsn)
		panic("error while connecting database")
	}
	DB.AutoMigrate(&domain.User{})
	DB.AutoMigrate(&domain.Admin{})
	DB.AutoMigrate(&domain.Product{})
	DB.AutoMigrate(&domain.Category{})
	DB.AutoMigrate(&domain.Size{})
	DB.AutoMigrate(&domain.Address{})
	DB.AutoMigrate(&domain.Cart{})
	DB.AutoMigrate(&domain.Order{})
	DB.AutoMigrate(&domain.OrderItem{})
	DB.AutoMigrate(&domain.Wishlist{})
	DB.AutoMigrate(&domain.Wallet{})
	DB.AutoMigrate(&domain.Coupon{})
	DB.AutoMigrate(&domain.PaymentMethod{})
	DB.AutoMigrate(&domain.RazorPay{})
	DB.AutoMigrate(&domain.UsedCoupon{})
	DB.AutoMigrate(&domain.Offer{})

	createAdmin(DB)

	return DB
}

func createAdmin(db *gorm.DB) {
	email := "admin@cityvibe.com"
	password := "Admin@123"

	hashed, _ := bcrypt.GenerateFromPassword([]byte(password), 10)

	err := db.FirstOrCreate(&domain.Admin{Email: email, Password: string(hashed)}).Error
	if err != nil {
		fmt.Print("failed to create admin")
	}
}
