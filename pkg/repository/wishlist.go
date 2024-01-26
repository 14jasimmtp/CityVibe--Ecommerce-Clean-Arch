package repository

import (
	"errors"
	"strconv"

	"github.com/14jasimmtp/CityVibe-Project-Clean-Architecture/pkg/db"
	"github.com/14jasimmtp/CityVibe-Project-Clean-Architecture/pkg/models"
	interfaceRepo "github.com/14jasimmtp/CityVibe-Project-Clean-Architecture/pkg/repository/interface"
	"gorm.io/gorm"
)

type WishlistRepo struct {
	DB *gorm.DB
}

func NewWishlistRepo(db *gorm.DB) interfaceRepo.WishlistRepo{
	return &WishlistRepo{DB: db}
}

func (clean *WishlistRepo) CheckExistInWishlist(userID uint, pid string) error {
	var product models.Product
	query := db.DB.Raw(`SELECT * FROM wishlists WHERE user_id = ? AND product_id = ?`, userID, pid).Scan(&product)
	if query.Error != nil {
		return errors.New(`something got wrong`)
	}
	if query.RowsAffected > 0 {
		return errors.New(`product already exist in wishlist`)
	}
	return nil
}

func (clean *WishlistRepo) AddProductToWishlist(pid string, userID uint) error {
	Pid, err := strconv.Atoi(pid)
	if err != nil {
		return err
	}
	query := db.DB.Exec(`INSERT INTO wishlists(user_id,product_id) VALUES(?,?)`, userID, Pid)
	if query.Error != nil {
		return errors.New(`something got wrong`)
	}
	return nil
}

func (clean *WishlistRepo) GetWishlistProducts(userID uint) ([]models.UpdateProduct, error) {
	var Products []models.UpdateProduct
	query := db.DB.Raw(
		`SELECT products.id,products.name,products.description,categories.category,sizes.size,products.stock,products.price,products.color
		 FROM products 
		 INNER JOIN wishlists ON products.id=wishlists.product_id
		 INNER JOIN categories ON products.category_id=categories.id
		 INNER JOIN sizes ON products.size_id=sizes.id
		 WHERE wishlists.user_id = ?`, userID,
	).Scan(&Products)
	if query.Error != nil {
		return []models.UpdateProduct{}, errors.New(`something went wrong`)
	}
	if query.RowsAffected < 1 {
		return []models.UpdateProduct{}, errors.New(`no products in wishlist`)
	}
	return Products, nil
}

func (clean *WishlistRepo) RemoveProductFromWishlist(pid string, userID uint) error {
	query := db.DB.Exec(`DELETE FROM wishlists WHERE user_id = ? AND product_id = ?`, userID, pid)
	if query.Error != nil {
		return errors.New(`something went wrong`)
	}
	return nil
}

func (clean *WishlistRepo) WishlistSingleProduct(id string) (models.Product, error) {
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
