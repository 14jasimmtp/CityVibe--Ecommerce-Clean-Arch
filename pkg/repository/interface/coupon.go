package interfaceRepo

import (
	"github.com/14jasimmtp/CityVibe-Project-Clean-Architecture/pkg/domain"
	"github.com/14jasimmtp/CityVibe-Project-Clean-Architecture/pkg/models"
)

type CouponRepo interface {
	CheckCouponExist(coupon string) (bool, error)
	CheckCouponExistWithID(coupon int) (bool, error)
	CreateCoupon(coupon models.Coupon) (domain.Coupon, error)
	DisableCoupon(coupon uint) error
	EnableCoupon(coupon uint) error
	GetCouponsForAdmin() ([]domain.Coupon, error)
	UpdateCoupon(coupon models.Coupon, coupon_id string) (domain.Coupon, error)
	GetDiscountRate(coupon string) (float64, error)
	UpdateCouponUsage(userID uint, coupon string) error
	UpdateCouponCount(coupon string) error
	CheckCouponUsage(userId uint, coupon string) error
	ViewUserCoupons(userID uint) ([]models.Couponlist, error)
	CouponAppliedOrNot(userID uint) error
	RemoveCouponAndChangePrice(userID uint) error
}
