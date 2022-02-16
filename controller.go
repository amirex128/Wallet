package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func walletBalanceGift(ctx *gin.Context) {
	var setGift SetGift
	if err := ctx.ShouldBindJSON(&setGift); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error in bind": err.Error(),
		})
		return
	}

	if err := validate.Struct(setGift); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error in validate": err.Error(),
		})
		return
	}
	var user User
	var gift Gift

	userResult := DB.Find(&user, "phone = ?", setGift.Phone)
	if userResult.RowsAffected == 0 {
		user.Phone = setGift.Phone
		DB.Create(&user)
	}
	transaction := DB.Begin()

	if transaction.First(&gift, "code= ?", setGift.Code).RowsAffected == 0 {
		transaction.Rollback()
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "gift code not found",
		})
		return
	}

	if gift.Count == 0 {
		transaction.Rollback()
		ctx.JSON(http.StatusOK, gin.H{
			"message": "all gift used",
		})
		return
	}

	if transaction.Find(&LogGift{}, "phone = ? and code = ?", setGift.Phone, setGift.Code).RowsAffected != 0 {
		transaction.Rollback()
		ctx.JSON(http.StatusOK, gin.H{
			"message": "gift code is used",
		})
		return
	}

	user.Balance += gift.Price
	if transaction.Model(&user).Where("phone", user.Phone).Update("balance", user.Balance).RowsAffected == 0 {
		transaction.Rollback()
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "increase balance failed",
		})
		return
	}

	gift.Count--
	if transaction.Model(&gift).Where("code", gift.Code).Update("count", gift.Count).RowsAffected == 0 {
		transaction.Rollback()
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "gift not save",
		})
		return
	}

	if transaction.Create(&LogGift{
		Phone: user.Phone,
		Code:  gift.Code,
	}).RowsAffected == 0 {
		transaction.Rollback()
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "create log failed",
		})
		return
	}

	if err := transaction.Commit().Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error in commit": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"current_balance": user.Balance,
		"amount_gift":     gift.Price,
		"gift_code":       gift.Code,
		"phone":           user.Phone,
	})
}

func walletBalance(c *gin.Context) {
	var user User
	resultUser := DB.Find(&user, "phone = ?", c.Query("phone"))
	if resultUser.RowsAffected == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "user not found",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"phone":   user.Phone,
		"balance": user.Balance,
	})
}
