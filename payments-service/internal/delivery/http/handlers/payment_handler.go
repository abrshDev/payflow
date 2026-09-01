package handlers

import (
	"net/http"

	"github.com/abrshDev/payment-service/internal/app/payment/commands"
	"github.com/abrshDev/payment-service/internal/app/payment/queries"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type PaymentHandler struct {
	createHandler    *commands.CreatePaymentHandler
	authorizeHandler *commands.AuthorizePaymentHandler
	captureHandler   *commands.CapturePaymentHandler
	refundHandler    *commands.RefundPaymentHandler
	getHandler       *queries.GetPaymentHandler
}

func NewPaymentHandler(
	create *commands.CreatePaymentHandler,
	authorize *commands.AuthorizePaymentHandler,
	capture *commands.CapturePaymentHandler,
	refund *commands.RefundPaymentHandler,
	get *queries.GetPaymentHandler,
) *PaymentHandler {
	return &PaymentHandler{
		createHandler:    create,
		authorizeHandler: authorize,
		captureHandler:   capture,
		refundHandler:    refund,
		getHandler:       get,
	}
}

type createPaymentRequest struct {
	MerchantID uuid.UUID `json:"merchant_id" binding:"required"`
	Amount     int64     `json:"amount" binding:"required"`
	Currency   string    `json:"currency" binding:"required"`
}

func (h *PaymentHandler) Create(c *gin.Context) {
	var req createPaymentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	payment, err := h.createHandler.Handle(c.Request.Context(), commands.CreatePaymentCommand{
		MerchantID: req.MerchantID,
		Amount:     req.Amount,
		Currency:   req.Currency,
	})
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, payment)
}

func (h *PaymentHandler) Authorize(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid payment id"})
		return
	}

	if err := h.authorizeHandler.Handle(c.Request.Context(), commands.AuthorizePaymentCommand{PaymentID: id}); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "authorized"})
}

func (h *PaymentHandler) Capture(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid payment id"})
		return
	}

	if err := h.captureHandler.Handle(c.Request.Context(), commands.CapturePaymentCommand{PaymentID: id}); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "captured"})
}

func (h *PaymentHandler) Refund(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid payment id"})
		return
	}

	if err := h.refundHandler.Handle(c.Request.Context(), commands.RefundPaymentCommand{PaymentID: id}); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "refunded"})
}

func (h *PaymentHandler) Get(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid payment id"})
		return
	}

	payment, err := h.getHandler.Handle(c.Request.Context(), queries.GetPaymentQuery{PaymentID: id})
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, payment)
}
