package usecase

import (
	"errors"

	"github.com/14jasimmtp/CityVibe-Project-Clean-Architecture/pkg/models"
	interfaceRepo "github.com/14jasimmtp/CityVibe-Project-Clean-Architecture/pkg/repository/interface"
	interfaceUsecase "github.com/14jasimmtp/CityVibe-Project-Clean-Architecture/pkg/usecase/interface"
	"github.com/14jasimmtp/CityVibe-Project-Clean-Architecture/pkg/utils"
)

type WishlistUseCase struct {
	WishlistRepo interfaceRepo.WishlistRepo
}

func NewWishlistUsecase(repo interfaceRepo.WishlistRepo) interfaceUsecase.WishlistUsecase{
	return &WishlistUseCase{WishlistRepo: repo}
}

func (clean *WishlistUseCase) AddProductToWishlist(pid string, Token string) error {

	userID, err := utils.ExtractUserIdFromToken(Token)
	if err != nil {
		return err
	}
	_, err = clean.WishlistRepo.WishlistSingleProduct(pid)
	if err != nil {
		return err
	}

	err = clean.WishlistRepo.CheckExistInWishlist(userID, pid)
	if err != nil {
		return err
	}
	err = clean.WishlistRepo.AddProductToWishlist(pid, userID)
	if err != nil {
		return err
	}
	return nil
}

func (clean *WishlistUseCase) ViewUserWishlist(Token string) ([]models.UpdateProduct, error) {
	userID, err := utils.ExtractUserIdFromToken(Token)
	if err != nil {
		return []models.UpdateProduct{}, err
	}
	WishedProducts, err := clean.WishlistRepo.GetWishlistProducts(userID)
	if err != nil {
		return []models.UpdateProduct{}, err
	}
	return WishedProducts, nil
}

func (clean *WishlistUseCase) RemoveProductFromWishlist(pid, Token string) error {
	userID, err := utils.ExtractUserIdFromToken(Token)
	if err != nil {
		return err
	}

	err = clean.WishlistRepo.CheckExistInWishlist(userID, pid)
	if err == nil {
		return errors.New(`no product found in wishlist with this id`)
	}
	err = clean.WishlistRepo.RemoveProductFromWishlist(pid, userID)
	if err != nil {
		return err
	}
	return nil
}
