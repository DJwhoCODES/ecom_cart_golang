package database

import "errors"

var (
	ErrCantFindProduct    = errors.New("can't find the product")
	ErrCantDecodeProducts = errors.New("can't decode products")
	ErrUserIdIsNotValid   = errors.New("user ID is not valid")
	ErrCantUpdateUser     = errors.New("can't update user")
	ErrCantRemoveItemCart = errors.New("can't remove item from cart")
	ErrCantGetItem        = errors.New("can't get item")
	ErrCantBuyCartItem    = errors.New("can't buy cart item")
)

func AddProductToCart() {

}

func RemoveCartItem() {

}

func BuyItemFromCart() {

}

func InstantBuy() {

}
