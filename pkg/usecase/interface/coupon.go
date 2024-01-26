package interfaceUsecase

import (
	"github.com/14jasimmtp/CityVibe-Project-Clean-Architecture/pkg/domain"
	"github.com/14jasimmtp/CityVibe-Project-Clean-Architecture/pkg/models"
)

type CouponUsecase interface {
	CreateCoupon(coupon models.Coupon) (domain.Coupon, error)
	DisableCoupon(coupon uint) error
	EnableCoupon(coupon uint) error
	GetCouponsForAdmin() ([]domain.Coupon, error)
	UpdateCoupon(coupon models.Coupon, couponID string) (domain.Coupon, error)
	ViewCouponsUser(Token string) ([]models.Couponlist, error)
	RemoveCoupon(Token string) error
}
