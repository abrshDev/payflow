package http

import (
	"github.com/abrshDev/payment-service/internal/delivery/http/handlers"
	"github.com/gin-gonic/gin"
)

func NewRouter(paymentHandler *handlers.PaymentHandler) *gin.Engine {
	r := gin.Default()

	payments := r.Group("/payments")
	{
		payments.POST("", paymentHandler.Create)
		payments.GET("/:id", paymentHandler.Get)
		payments.POST("/:id/authorize", paymentHandler.Authorize)
		payments.POST("/:id/capture", paymentHandler.Capture)
		payments.POST("/:id/refund", paymentHandler.Refund)

	}

	return r
}
