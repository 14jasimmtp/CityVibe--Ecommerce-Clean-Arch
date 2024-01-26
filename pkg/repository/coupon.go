package repository

import (
	"errors"
	"fmt"

	"github.com/14jasimmtp/CityVibe-Project-Clean-Architecture/pkg/db"
	"github.com/14jasimmtp/CityVibe-Project-Clean-Architecture/pkg/domain"
	"github.com/14jasimmtp/CityVibe-Project-Clean-Architecture/pkg/models"
	interfaceRepo "github.com/14jasimmtp/CityVibe-Project-Clean-Architecture/pkg/repository/interface"
	"gorm.io/gorm"
)

type CouponRepo struct {
	DB *gorm.DB
}

func NewCouponRepo(db *gorm.DB) interfaceRepo.CouponRepo{
	return &CouponRepo{DB: db}
}

func (clean *CouponRepo) CheckCouponExist(coupon string) (bool, error) {
	var couponvalid domain.Coupon
	query := `SELECT * FROM coupons WHERE coupon = ?`

	db := db.DB.Raw(query, coupon).Scan(&couponvalid)
	if db.Error != nil {
		return true, errors.New(`something went wrong`)
	}

	if db.RowsAffected > 0 {
		return true, errors.New(`already exist`)
	}

	return false, nil

}

func (clean *CouponRepo) CheckCouponExistWithID(coupon int) (bool, error) {
	var couponvalid domain.Coupon
	query := `SELECT * FROM coupons WHERE id = ?`

	db := db.DB.Raw(query, coupon).Scan(&couponvalid)
	if db.Error != nil {
		return true, errors.New(`something went wrong`)
	}

	if db.RowsAffected > 0 {
		return true, nil
	}

	return false, nil
}

func (clean *CouponRepo) CreateCoupon(coupon models.Coupon) (domain.Coupon, error) {
	var Coupons domain.Coupon
	query := db.DB.Raw(`INSERT INTO coupons (coupon,discount_percentage,usage_limit) VALUES (?,?,?) RETURNING id,coupon,discount_percentage,usage_limit,active`, coupon.Coupon, coupon.DiscoutPercentage, coupon.UsageLimit).Scan(&Coupons)
	if query.Error != nil {
		return domain.Coupon{}, errors.New(`something went wrong`)
	}
	return Coupons, nil
}

func (clean *CouponRepo) DisableCoupon(coupon uint) error {
	query := db.DB.Exec(`UPDATE coupons SET active = false WHERE id = ?`, coupon)
	if query.RowsAffected < 1 {
		return errors.New(`no coupons found with this id`)
	}
	if query.Error != nil {
		return errors.New(`something went wrong`)
	}
	return nil
}

func (clean *CouponRepo) EnableCoupon(coupon uint) error {
	query := db.DB.Exec(`UPDATE coupons SET active = true WHERE id = ?`, coupon)
	if query.RowsAffected < 1 {
		return errors.New(`no coupons found with this id`)
	}
	if query.Error != nil {
		return errors.New(`something went wrong`)
	}
	return nil
}

func (clean *CouponRepo) GetCouponsForAdmin() ([]domain.Coupon, error) {
	var Coupons []domain.Coupon
	query := db.DB.Raw(`SELECT * FROM coupons`).Scan(&Coupons)
	if query.Error != nil {
		return []domain.Coupon{}, errors.New(`something went wrong`)
	}
	if query.RowsAffected < 1 {
		return []domain.Coupon{}, errors.New(`no coupons added.Add a coupon to view`)
	}
	return Coupons, nil
}

func (clean *CouponRepo) UpdateCoupon(coupon models.Coupon, coupon_id string) (domain.Coupon, error) {
	var coupons domain.Coupon
	query := db.DB.Raw(`UPDATE coupons SET coupon = ? ,discount_percentage = ? RETURNING id,coupon,discount_percentage,valid`, coupon.Coupon, coupon.DiscoutPercentage).Scan(&coupons)
	if query.Error != nil {
		return domain.Coupon{}, errors.New(`something went wrong`)
	}
	return coupons, nil
}

func (clean *CouponRepo) GetDiscountRate(coupon string) (float64, error) {
	var discountRate float64
	query := db.DB.Raw(`SELECT discount_percentage from coupons WHERE coupon = ? AND active = true`, coupon).Scan(&discountRate)
	if query.Error != nil {
		return 0.0, errors.New(`something went wrong`)
	}
	if query.RowsAffected < 1 {
		return 0.0, errors.New(`no coupons found`)
	}
	return discountRate, nil
}

func (clean *CouponRepo) UpdateCouponUsage(userID uint, coupon string) error {
	query := db.DB.Exec(`insert into used_coupons (user_id,coupon) values(?,?)`, userID, coupon)
	if query.Error != nil {
		return errors.New(`something went wrong`)
	}
	return nil
}

func (clean *CouponRepo) UpdateCouponCount(coupon string) error {
	query := db.DB.Exec(`UPDATE coupons SET usage_limit = usage_limit - 1`)
	if query.Error != nil {
		return errors.New(`something went wrong`)
	}
	return nil
}

func (clean *CouponRepo) CheckCouponUsage(userId uint, coupon string) error {
	var count int
	query := db.DB.Raw(`SELECT count(*) from used_coupons where user_id = ? AND coupon = ?`, userId, coupon).Scan(&count)
	if query.Error != nil {
		return errors.New(`something went wrong`)
	}
	if count >= 1 {
		return errors.New(`coupon already used`)
	}
	return nil
}

func (clean *CouponRepo) ViewUserCoupons(userID uint) ([]models.Couponlist, error) {
	var coupons []models.Couponlist
	query := db.DB.Raw(`SELECT coupon,discount_percentage FROM coupons `).Scan(&coupons)
	if query.Error != nil {
		return []models.Couponlist{}, errors.New(`something went wrong`)
	}
	if query.RowsAffected == 0 {
		return []models.Couponlist{}, errors.New(`no coupons found`)
	}

	return coupons, nil
}

func (clean *CouponRepo) CouponAppliedOrNot(userID uint) error {
	var final_price float64
	query := db.DB.Raw(`select final_price from carts where user_id = ? LIMIT 1`, userID).Scan(&final_price)
	if query.Error != nil {
		fmt.Println(query.Error)
		return errors.New(`something went wrong`)
	}
	if query.RowsAffected == 0 {
		return errors.New(`no products found`)
	}

	if final_price == 0.0 {
		return errors.New(`coupon is not applied to remove`)
	}
	return nil
}

func (clean *CouponRepo) RemoveCouponAndChangePrice(userID uint) error {
	query := db.DB.Exec(`update carts set final_price = 0 where user_id = ?`, userID)
	if query.Error != nil {
		fmt.Println(query.Error)
		return errors.New(`something went wrong`)

	}
	return nil
}
