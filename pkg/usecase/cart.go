package usecase

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/14jasimmtp/CityVibe-Project-Clean-Architecture/pkg/models"
	interfaceRepo "github.com/14jasimmtp/CityVibe-Project-Clean-Architecture/pkg/repository/interface"
	interfaceUsecase "github.com/14jasimmtp/CityVibe-Project-Clean-Architecture/pkg/usecase/interface"
	"github.com/14jasimmtp/CityVibe-Project-Clean-Architecture/pkg/utils"
)

type CartUseCase struct {
	CartRepo interfaceRepo.CartRepo
}

func NewCartUsecase(repo interfaceRepo.CartRepo) interfaceUsecase.CartUsecase{
	return &CartUseCase{CartRepo: repo}
}

func (clean *CartUseCase) ViewCart(Token string) (models.CartResponse, error) {
	UserId, err := utils.ExtractUserIdFromToken(Token)
	if err != nil {
		return models.CartResponse{}, err
	}

	Cart, err := clean.CartRepo.DisplayCart(UserId)
	if err != nil {
		return models.CartResponse{}, err
	}

	cartTotal, err := clean.CartRepo.CartTotalAmount(UserId)
	if err != nil {
		return models.CartResponse{}, err
	}

	return models.CartResponse{
		TotalPrice: cartTotal,
		Cart:       Cart,
	}, nil

}

func (clean *CartUseCase) AddToCart(pid, Token string) (models.CartResponse, error) {

	_, err := clean.CartRepo.CheckSingleProduct(pid)
	if err != nil {
		return models.CartResponse{}, errors.New("product doesn't exist")
	}

	UserId, err := utils.ExtractUserIdFromToken(Token)
	if err != nil {
		return models.CartResponse{}, err
	}

	ProId, err := strconv.Atoi(pid)
	if err != nil {
		return models.CartResponse{}, err
	}

	productPrize, err := clean.CartRepo.GetCartProductAmountFromID(pid)
	if err != nil {
		return models.CartResponse{}, err
	}
	true, err := clean.CartRepo.CheckProductExistInCart(UserId, pid)
	if err != nil {
		return models.CartResponse{}, err
	}
	fmt.Println(true)
	if true {
		TotalProductAmount, err := clean.CartRepo.TotalPrizeOfProductInCart(UserId, pid)
		if err != nil {
			return models.CartResponse{}, err
		}

		err = clean.CartRepo.UpdateCart(1, TotalProductAmount+productPrize, UserId, pid)
		if err != nil {
			return models.CartResponse{}, err
		}
	} else {
		if err := clean.CartRepo.CheckCartStock(ProId); err != nil {
			return models.CartResponse{}, err
		}
		err := clean.CartRepo.AddToCart(ProId, UserId, productPrize)
		if err != nil {
			return models.CartResponse{}, err
		}
	}

	CartDetails, err := clean.CartRepo.DisplayCart(UserId)
	if err != nil {
		return models.CartResponse{}, err
	}

	cartTotalAmount, err := clean.CartRepo.CartTotalAmount(UserId)
	if err != nil {
		return models.CartResponse{}, err
	}

	return models.CartResponse{
		TotalPrice: cartTotalAmount,
		Cart:       CartDetails,
	}, nil
}

func (clean *CartUseCase) RemoveProductsFromCart(pid, Token string) (models.CartResponse, error) {
	ProId, err := strconv.Atoi(pid)
	if err != nil {
		return models.CartResponse{}, err
	}

	UserId, err := utils.ExtractUserIdFromToken(Token)
	if err != nil {
		return models.CartResponse{}, err
	}

	err = clean.CartRepo.RemoveProductFromCart(ProId, UserId)
	if err != nil {
		return models.CartResponse{}, err
	}

	updatedCart, err := clean.CartRepo.DisplayCart(UserId)
	if err != nil {
		return models.CartResponse{}, err
	}
	cartTotal, err := clean.CartRepo.CartTotalAmount(UserId)
	if err != nil {
		return models.CartResponse{}, err
	}
	return models.CartResponse{
		TotalPrice: cartTotal,
		Cart:       updatedCart,
	}, nil
}

func (clean *CartUseCase) UpdateQuantityFromCart(Token, pid, quantity string) ([]models.Cart, error) {
	UserId, err := utils.ExtractUserIdFromToken(Token)
	if err != nil {
		return []models.Cart{}, err
	}

	updatedCart, err := clean.CartRepo.UpdateQuantity(UserId, pid, quantity)
	if err != nil {
		return []models.Cart{}, err
	}

	return updatedCart, nil
}

func (clean *CartUseCase) UpdateQuantityIncrease(Token, pid string) error {
	userID, err := utils.ExtractUserIdFromToken(Token)
	if err != nil {
		return err
	}

	check, err := clean.CartRepo.CheckProductExistInCart(userID, pid)
	if err != nil {
		return err
	}
	if !check {
		return errors.New(`no products found in cart with this id`)
	}

	err = clean.CartRepo.UpdateQuantityAdd(userID, pid)
	if err != nil {
		return err
	}

	return nil
}

func (clean *CartUseCase) UpdatePriceAdd(Token, pid string) error {
	userID, err := utils.ExtractUserIdFromToken(Token)
	if err != nil {
		return err
	}
	err = clean.CartRepo.UpdateTotalPrice(userID, pid)
	if err != nil {
		return err
	}
	return nil
}

func (clean *CartUseCase) UpdateQuantityDecrease(Token, pid string) error {
	userID, err := utils.ExtractUserIdFromToken(Token)
	if err != nil {
		return err
	}

	check, err := clean.CartRepo.CheckProductExistInCart(userID, pid)
	if err != nil {
		return err
	}

	if !check {
		return errors.New(`no products found in cart with this id`)
	}

	quantity, err := clean.CartRepo.ProductQuantityCart(userID, pid)
	if err != nil {
		return err
	}

	if quantity == 1 {
		return errors.New(`quantity is 1 .can't reduce anymore`)
	}

	err = clean.CartRepo.UpdateQuantityless(userID, pid)
	if err != nil {
		return err
	}
	return nil
}

func (clean *CartUseCase) UpdatePriceDecrease(Token, pid string) error {
	userID, err := utils.ExtractUserIdFromToken(Token)
	if err != nil {
		return err
	}
	err = clean.CartRepo.UpdateTotalPrice(userID, pid)
	if err != nil {
		return err
	}
	return nil
}

func (clean *CartUseCase) EraseCart(Token string) (models.CartResponse, error) {
	userID, err := utils.ExtractUserIdFromToken(Token)
	if err != nil {
		return models.CartResponse{}, err
	}
	ok, err := clean.CartRepo.CartExist(userID)
	if err != nil {
		return models.CartResponse{}, err
	}
	if !ok {
		return models.CartResponse{}, errors.New("cart already empty")
	}
	if err := clean.CartRepo.EmptyCart(userID); err != nil {
		return models.CartResponse{}, err
	}

	cartTotal, err := clean.CartRepo.CartTotalAmount(userID)

	if err != nil {
		return models.CartResponse{}, err
	}

	cartResponse := models.CartResponse{
		TotalPrice: cartTotal,
		Cart:       []models.Cart{},
	}

	return cartResponse, nil
}
