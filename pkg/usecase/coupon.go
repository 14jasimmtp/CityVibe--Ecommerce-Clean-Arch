package usecase

import (
	"errors"
	"strconv"

	"github.com/14jasimmtp/CityVibe-Project-Clean-Architecture/pkg/domain"
	"github.com/14jasimmtp/CityVibe-Project-Clean-Architecture/pkg/models"
	interfaceRepo "github.com/14jasimmtp/CityVibe-Project-Clean-Architecture/pkg/repository/interface"
	interfaceUsecase "github.com/14jasimmtp/CityVibe-Project-Clean-Architecture/pkg/usecase/interface"
	"github.com/14jasimmtp/CityVibe-Project-Clean-Architecture/pkg/utils"
)

type CouponUseCase struct {
	CouponRepo interfaceRepo.CouponRepo
}

func NewCouponUsecase(repo interfaceRepo.CouponRepo) interfaceUsecase.CouponUsecase{
	return &CouponUseCase{CouponRepo: repo}
}

func (clean *CouponUseCase) CreateCoupon(coupon models.Coupon) (domain.Coupon, error) {
	CheckCouponExist, err := clean.CouponRepo.CheckCouponExist(coupon.Coupon)
	if CheckCouponExist {
		return domain.Coupon{}, errors.New(`coupon already exist`)
	}
	if err != nil {
		return domain.Coupon{}, err
	}

	Coupon, err := clean.CouponRepo.CreateCoupon(coupon)
	if err != nil {
		return domain.Coupon{}, err
	}
	return Coupon, nil
}

func (clean *CouponUseCase) DisableCoupon(coupon uint) error {
	err := clean.CouponRepo.DisableCoupon(coupon)
	if err != nil {
		return err
	}
	return nil
}

func (clean *CouponUseCase) EnableCoupon(coupon uint) error {
	err := clean.CouponRepo.EnableCoupon(coupon)
	if err != nil {
		return err
	}
	return nil
}

func (clean *CouponUseCase) GetCouponsForAdmin() ([]domain.Coupon, error) {
	Coupons, err := clean.CouponRepo.GetCouponsForAdmin()
	if err != nil {
		return []domain.Coupon{}, err
	}
	return Coupons, nil
}

func (clean *CouponUseCase) UpdateCoupon(coupon models.Coupon, coupon_id string) (domain.Coupon, error) {
	cid, err := strconv.Atoi(coupon_id)
	if err != nil {
		return domain.Coupon{}, err
	}
	CheckCoupon, err := clean.CouponRepo.CheckCouponExistWithID(cid)

	if !CheckCoupon {
		return domain.Coupon{}, errors.New(`no coupon found with this id`)
	}
	if err != nil {
		return domain.Coupon{}, err
	}
	Coupon, err := clean.CouponRepo.UpdateCoupon(coupon, coupon_id)
	if err != nil {
		return domain.Coupon{}, err
	}
	return Coupon, nil
}

func (clean *CouponUseCase) ViewCouponsUser(Token string) ([]models.Couponlist, error) {
	UserID, err := utils.ExtractUserIdFromToken(Token)
	if err != nil {
		return []models.Couponlist{}, errors.New(`error in token`)
	}

	Coupons, err := clean.CouponRepo.ViewUserCoupons(UserID)
	if err != nil {
		return []models.Couponlist{}, err
	}

	return Coupons, nil
}

func (clean *CouponUseCase) RemoveCoupon(Token string) error {
	userID, err := utils.ExtractUserIdFromToken(Token)
	if err != nil {
		return err
	}
	if err := clean.CouponRepo.CouponAppliedOrNot(userID); err != nil {
		return err
	}
	if err := clean.CouponRepo.RemoveCouponAndChangePrice(userID); err != nil {
		return err
	}
	return nil
}
