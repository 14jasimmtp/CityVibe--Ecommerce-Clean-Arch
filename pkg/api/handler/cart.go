package handlers

import (
	"net/http"

	interfaceUsecase "github.com/14jasimmtp/CityVibe-Project-Clean-Architecture/pkg/usecase/interface"
	"github.com/gin-gonic/gin"
)

type CartHandler struct{
	CartUsecase interfaceUsecase.CartUsecase
}

func NewCartHandler(usecase interfaceUsecase.CartUsecase) *CartHandler{
	return &CartHandler{CartUsecase: usecase}
}

// AddToCart godoc
// @Summary Add product to user's cart
// @Description Add a product to the user's cart based on the provided product ID.
// @Tags Cart
// @Accept json
// @Produce json
// @Param product_id query string true "Product ID to add to the cart"
// @Success 200 {object} string "message": "Product added to cart successfully", "Cart": Cart
// @Failure 400 {object} string "error": "Bad Request"
// @Failure 401 {object} string "error": "Unauthorized"
// @Failure 500 {object} string "error": "Internal Server Error"
// @Router /cart [post]
func (clean *CartHandler) AddToCart(c *gin.Context) {
	pid := c.Query("product_id")

	Token, err := c.Cookie("Authorisation")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	Cart, err := clean.CartUsecase.AddToCart(pid, Token)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "product added to cart successfully", "Cart": Cart})
}

// ViewCart godoc
// @Summary View user's cart
// @Description Retrieve details of the user's cart.
// @Tags Cart
// @Accept json
// @Produce json
// @Success 200 {object} string "message": "Cart details", "Cart": UserCart
// @Failure 400 {object} string "error": "Bad Request"
// @Failure 401 {object} string "error": "Unauthorized"
// @Failure 500 {object} string "error": "Internal Server Error"
// @Router /cart [get]
func (clean *CartHandler) ViewCart(c *gin.Context) {
	Token, err := c.Cookie("Authorisation")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	UserCart, err := clean.CartUsecase.ViewCart(Token)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Cart details", "Cart": UserCart})

}

// RemoveProductsFromCart godoc
// @Summary Remove product from user's cart
// @Description Remove a product from the user's cart based on the provided product ID.
// @Tags Cart
// @Accept json
// @Produce json
// @Param product_id query string true "Product ID to remove from the cart"
// @Success 200 {object} string "message": "Product removed from cart successfully", "Cart": Cart
// @Failure 400 {object} string "error": "Bad Request"
// @Failure 401 {object} string "error": "Unauthorized"
// @Failure 500 {object} string "error": "Internal Server Error"
// @Router /cart [delete]
func (clean *CartHandler) RemoveProductsFromCart(c *gin.Context) {
	id := c.Query("product_id")
	Token, err := c.Cookie("Authorisation")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	Cart,err := clean.CartUsecase.RemoveProductsFromCart(id, Token)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "product removed from cart successfully","Cart":Cart})

}

// IncreaseQuantityUpdate godoc
// @Summary Increase quantity of a product in the user's cart
// @Description Increase the quantity of a product in the user's cart based on the provided product ID.
// @Tags Cart
// @Accept json
// @Produce json
// @Param product_id query string true "Product ID to increase quantity"
// @Success 200 {object} string "message": "Quantity added successfully"
// @Failure 400 {object} string "error": "Bad Request"
// @Failure 401 {object} string "error": "Unauthorized"
// @Failure 500 {object} string "error": "Internal Server Error"
// @Router /cart/add-quantity [put]
func (clean *CartHandler) IncreaseQuantityUpdate(c *gin.Context) {
	pid := c.Query("product_id")

	Token, err := c.Cookie("Authorisation")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err=clean.CartUsecase.UpdateQuantityIncrease(Token,pid)
	if err != nil{
		c.JSON(http.StatusBadRequest, gin.H{"error": "unable to add quantity"})
		return
	}

	err=clean.CartUsecase.UpdatePriceAdd(Token,pid)
	if err != nil{
		c.JSON(http.StatusBadRequest, gin.H{"error": "unable to add quantity"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "quantity added successfully"})

}

// DecreaseQuantityUpdate godoc
// @Summary Decrease quantity of a product in the user's cart
// @Description Decrease the quantity of a product in the user's cart based on the provided product ID.
// @Tags Cart
// @Accept json
// @Produce json
// @Param product_id query string true "Product ID to decrease quantity"
// @Success 200 {object} string "message": "Quantity decreased by 1 successfully"
// @Failure 400 {object} string "error": "Bad Request"
// @Failure 401 {object} string "error": "Unauthorized"
// @Failure 500 {object} string "error": "Internal Server Error"
// @Router /cart/reduce-quantity [put]
func (clean *CartHandler) DecreaseQuantityUpdate(c *gin.Context) {
	pid := c.Query("product_id")

	Token, err := c.Cookie("Authorisation")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err=clean.CartUsecase.UpdateQuantityDecrease(Token,pid)
	if err != nil{
		c.JSON(http.StatusBadRequest, gin.H{"error": "unable to decrease quantity"})
		return
	}

	err=clean.CartUsecase.UpdatePriceDecrease(Token,pid)
	if err != nil{
		c.JSON(http.StatusBadRequest, gin.H{"error": "unable to decrease quantity"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "quantity decreased by 1 successfully"})
}