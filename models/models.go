package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID              primitive.ObjectID `json:"_id" bson:"_id,omitempty"`
	First_Name      *string            `json:"first_name" bson:"first_name,omitempty" validate:"required,min=2,max=100"`
	Last_Name       *string            `json:"last_name" bson:"last_name,omitempty" validate:"required,min=2,max=100"`
	Password        *string            `json:"password" bson:"password,omitempty" validate:"required,min=6"`
	Email           *string            `json:"email" bson:"email,omitempty" validate:"email,required"`
	Phone           *string            `json:"phone" bson:"phone,omitempty" validate:"required"`
	Access_Token    *string            `json:"access_token" bson:"access_token,omitempty"`
	Refresh_Token   *string            `json:"refresh_token" bson:"refresh_token,omitempty"`
	Created_At      time.Time          `json:"created_at" bson:"created_at"`
	Updated_At      time.Time          `json:"updated_at" bson:"updated_at"`
	User_Id         string             `json:"user_id" bson:"user_id"`
	User_Cart       []ProductUser      `json:"user_cart" bson:"user_cart"`
	Address_Details []Address          `json:"address_details" bson:"address_details"`
	Order_Status    []Order            `json:"order_status" bson:"order_status"`
}

type Product struct {
	Product_Id   primitive.ObjectID `json:"product_id" bson:"product_id"`
	Product_Name *string            `json:"product_name" bson:"product_name" validate:"required"`
	Price        *uint32            `json:"price" bson:"price" validate:"required,gte=0"`
	Rating       *uint8             `json:"rating" bson:"rating,omitempty"`
	Image        *string            `json:"image" bson:"image,omitempty"`
}

type ProductUser struct {
	Product_Id   primitive.ObjectID `json:"product_id" bson:"product_id"`
	Product_Name *string            `json:"product_name" bson:"product_name"`
	Price        *uint32            `json:"price" bson:"price"`
	Rating       *uint8             `json:"rating" bson:"rating"`
	Image        *string            `json:"image" bson:"image"`
}

type Address struct {
	Address_Id primitive.ObjectID `json:"address_id" bson:"address_id"`
	House      *string            `json:"house" bson:"house" validate:"required"`
	Street     *string            `json:"street" bson:"street" validate:"required"`
	City       *string            `json:"city" bson:"city" validate:"required"`
	Pincode    *string            `json:"pincode" bson:"pincode" validate:"required,len=6"`
}

type Order struct {
	Order_Id       primitive.ObjectID `json:"order_id" bson:"order_id"`
	Order_Cart     []ProductUser      `json:"order_cart" bson:"order_cart"`
	Ordered_At     time.Time          `json:"ordered_at" bson:"ordered_at"`
	Price          *uint32            `json:"price" bson:"price"`
	Discount       *uint8             `json:"discount" bson:"discount,omitempty"`
	Payment_Method Payment            `json:"payment_method" bson:"payment_method"`
}

type Payment struct {
	Digital bool `json:"digital" bson:"digital"`
	COD     bool `json:"cod" bson:"cod"`
}
