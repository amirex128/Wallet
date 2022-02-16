package main

import (
	"github.com/go-playground/validator/v10"
	"time"
)

var validate *validator.Validate

type SetGift struct {
	Phone string `from:"phone" json:"phone" validate:"required,len=11,numeric" `
	Code  string `from:"code" json:"code" validate:"required,min=1,max=16"`
}
type Gift struct {
	Code  string `from:"code" json:"code" gorm:"primary_key"`
	Price uint   `from:"price" json:"price"`
	Count uint   `from:"count" json:"count"`
}

type LogGift struct {
	ID        uint `gorm:"primarykey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	Phone     string `from:"phone" json:"phone" `
	Code      string `from:"code" json:"code"`
}

type User struct {
	Phone   string `from:"phone" json:"phone" gorm:"primary_key"`
	Balance uint   `from:"balance" json:"balance"`
}
