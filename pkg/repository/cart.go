package repository

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/14jasimmtp/CityVibe-Project-Clean-Architecture/pkg/db"
	"github.com/14jasimmtp/CityVibe-Project-Clean-Architecture/pkg/models"
	interfaceRepo "github.com/14jasimmtp/CityVibe-Project-Clean-Architecture/pkg/repository/interface"
	"gorm.io/gorm"
)

type CartRepo struct {
	DB *gorm.DB
}

func NewCartRepo(db *gorm.DB) interfaceRepo.CartRepo {
	return &CartRepo{DB: db}
}

func (clean *CartRepo) AddToCart(pid int, userid uint, productAmount float64) error {
	query := db.DB.Exec(`INSERT INTO carts (user_id,product_id,quantity,price) VALUES (?,?,?,?)`, userid, pid, 1, productAmount)
	if query.Error != nil {
		return query.Error
	}

	return nil
}

func (clean *CartRepo) DisplayCart(userid uint) ([]models.Cart, error) {

	var count int
	if err := db.DB.Raw("SELECT COUNT(*) FROM carts WHERE user_id = ? ", userid).First(&count).Error; err != nil {
		return []models.Cart{}, err
	}

	if count == 0 {
		return []models.Cart{}, nil
	}

	var Cart []models.Cart

	if err := db.DB.Raw("SELECT carts.user_id,users.firstname as user_name,carts.product_id,products.name as product_name,categories.category as category,carts.quantity,carts.price,carts.final_price FROM carts inner join users on carts.user_id = users.id inner join products on carts.product_id = products.id inner join categories on categories.id = products.category_id where user_id = ?", userid).First(&Cart).Error; err != nil {
		return []models.Cart{}, err
	}

	return Cart, nil
}

func (clean *CartRepo) RemoveProductFromCart(pid int, userid uint) error {
	query := db.DB.Exec(`DELETE FROM carts WHERE product_id = ? AND user_id = ?`, pid, userid)
	if query.Error != nil {
		return query.Error
	}
	if query.RowsAffected == 0 {
		return errors.New(`no products found in cart`)
	}

	return nil
}

func (clean *CartRepo) CheckProductExistInCart(userId uint, pid string) (bool, error) {
	var count int
	query := db.DB.Raw(`SELECT COUNT(*) FROM carts WHERE user_id = ? AND product_id = ?`, userId, pid).Scan(&count)
	if query.Error != nil {
		return false, errors.New(`something went wrong`)
	}
	fmt.Println(count)

	if count > 0 {
		return true, nil
	}

	return false, nil
}

func (clean *CartRepo) CheckSingleProduct(id string) (models.Product, error) {
	var product models.Product
	idint, err := strconv.Atoi(id)
	if err != nil {
		return models.Product{}, errors.New("error while converting id to int")
	}

	query := db.DB.Raw("SELECT products.id as id,name,description,categories.category,sizes.size,stock,color,price FROM products INNER JOIN categories ON categories.id = products.category_id INNER JOIN sizes ON sizes.id=products.size_id WHERE products.id = ?", idint).Scan(&product)
	if product.Name == "" {
		return models.Product{}, errors.New("no products found with this id")
	}

	if query.Error != nil {
		return models.Product{}, errors.New("something went wrong")
	}

	return product, nil
}

func (clean *CartRepo) GetCartProductAmountFromID(pid string) (float64, error) {
	var price struct {
		Price      float64
		OfferPrize float64
	}

	if err := db.DB.Raw("select price,offer_prize from products where id = ?", pid).Scan(&price).Error; err != nil {
		return 0.0, err
	}
	if price.OfferPrize != 0 {
		return price.OfferPrize, nil
	}
	return price.Price, nil
}

func (clean *CartRepo) UpdateQuantity(userid uint, pid, quantity string) ([]models.Cart, error) {
	query := db.DB.Raw(`UPDATE carts SET quantity = ? WHERE user_id = ? AND product_id = ?`, quantity, userid, pid)
	if query.Error != nil {
		return []models.Cart{}, query.Error
	}
	if query.RowsAffected == 0 {
		return []models.Cart{}, errors.New(`no products found to update in cart`)
	}

	var Cart []models.Cart
	if err := db.DB.Raw("SELECT carts.user_id,users.firstname as user_name,carts.product_id,products.name as product_name,carts.quantity,carts.price FROM carts inner join users on carts.user_id = users.id inner join products on carts.product_id = products.id where user_id = ?", userid).First(&Cart).Error; err != nil {
		return []models.Cart{}, err
	}

	return Cart, nil
}

func (clean *CartRepo) CartTotalAmount(userid uint) (float64, error) {
	var Amount float64
	err := db.DB.Raw(`SELECT SUM(price) FROM carts WHERE user_id = ?`, userid).Scan(&Amount).Error

	if err != nil {
		return 0.0, nil
	}
	return Amount, nil
}

func (clean *CartRepo) CartFinalPrice(userid uint) (float64, error) {
	var Amount float64
	err := db.DB.Raw(`SELECT SUM(final_price) FROM carts WHERE user_id = ?`, userid).Scan(&Amount).Error

	if err != nil {
		return 0.0, nil
	}
	return Amount, nil
}

func (clean *CartRepo) CheckCartExist(userID uint) bool {
	var count int
	if err := db.DB.Raw("SELECT COUNT(*) FROM carts WHERE  user_id = ?", userID).Scan(&count).Error; err != nil {
		return false
	}
	return count > 0
}

func (clean *CartRepo) UpdateCart(quantity int, price float64, userID uint, product_id string) error {

	if err := db.DB.Exec("update carts set quantity = quantity + $1, price = $2 where user_id = $3 and product_id = $4", quantity, price, userID, product_id).Error; err != nil {
		return err
	}

	return nil

}

func (clean *CartRepo) CheckCartStock(pid int) error {
	var stock int
	db.DB.Raw(`SELECT stock from products WHERE id = ?`, pid).Scan(&stock)
	if stock < 1 {
		return errors.New("product out of stock")
	}
	return nil
}

func (clean *CartRepo) TotalPrizeOfProductInCart(userID uint, productID string) (float64, error) {

	var totalPrice float64
	if err := db.DB.Raw("select sum(price) as total_price from carts where user_id = ? and product_id = ?", userID, productID).Scan(&totalPrice).Error; err != nil {
		return 0.0, err
	}
	return totalPrice, nil
}

func (clean *CartRepo) UpdateQuantityAdd(id uint, prdt_id string) error {
	err := db.DB.Exec("UPDATE Carts SET quantity = quantity + 1 WHERE user_id=$1 AND product_id = $2 ", id, prdt_id).Error
	if err != nil {
		return err
	}
	return nil
}

func (clean *CartRepo) UpdateTotalPrice(id uint, product_id string) error {
	err := db.DB.Exec("UPDATE carts SET price = carts.quantity * products.price FROM products  WHERE carts.product_id = products.id AND carts.user_id = $1 AND carts.product_id = $2", id, product_id).Error
	if err != nil {
		return err
	}
	return nil
}

func (clean *CartRepo) UpdateQuantityless(id uint, prdt_id string) error {
	err := db.DB.Exec("UPDATE Carts SET quantity = quantity - 1 WHERE user_id=$1 AND product_id = $2 ", id, prdt_id).Error
	if err != nil {
		return err
	}
	return nil
}

func (clean *CartRepo) CartExist(userID uint) (bool, error) {
	var count int
	if err := db.DB.Raw("SELECT COUNT(*) FROM carts WHERE user_id = ? ", userID).Scan(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil

}

func (clean *CartRepo) EmptyCart(userID uint) error {

	if err := db.DB.Exec("DELETE FROM carts WHERE user_id = ? ", userID).Error; err != nil {
		return err
	}

	return nil

}

func (clean *CartRepo) ProductQuantityCart(userID uint, pid string) (int, error) {
	var quantity int
	query := db.DB.Raw(`SELECT quantity FROM carts WHERE user_id = ? AND product_id = ?`, userID, pid).Scan(&quantity).Error
	if query != nil {
		return 0, errors.New(`something went wrong`)
	}
	return quantity, nil
}
