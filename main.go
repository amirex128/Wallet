package main

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"os"
	"os/signal"
	"syscall"
)

func main() {

	go func() {
		InitOrm()

		r := gin.Default()
		validate = validator.New()

		r.GET("/wallet/balance", walletBalance)

		r.POST("/wallet/balance/gift", walletBalanceGift)

		if err := r.Run(":8083"); err != nil {
			panic("failed to run server")
		}

	}()

	// Create a channel to listen for an interrupt or terminate signal from the OS.
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	<-ch
}
