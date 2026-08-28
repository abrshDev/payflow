package entities

import (
	"testing"

	domainerrors "github.com/abrshDev/payment-service/internal/domain/errors"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestNewPayment_Valid(t *testing.T) {
	p, err := NewPayment(uuid.New(), 1000, "USD")
	assert.NoError(t, err)
	assert.Equal(t, StatusCreated, p.Status)
	assert.Len(t, p.History, 1)
}

func TestNewPayment_InvalidAmount(t *testing.T) {
	_, err := NewPayment(uuid.New(), 0, "USD")
	assert.ErrorIs(t, err, domainerrors.ErrInvalidAmount)

	_, err = NewPayment(uuid.New(), -500, "USD")
	assert.ErrorIs(t, err, domainerrors.ErrInvalidAmount)
}

func TestPayment_HappyPath(t *testing.T) {
	p, _ := NewPayment(uuid.New(), 1000, "USD")

	assert.NoError(t, p.Authorize())
	assert.Equal(t, StatusAuthorized, p.Status)

	assert.NoError(t, p.Capture())
	assert.Equal(t, StatusCaptured, p.Status)

	assert.NoError(t, p.Refund())
	assert.Equal(t, StatusRefunded, p.Status)

	// created -> authorized -> captured -> refunded = 4 history entries
	assert.Len(t, p.History, 4)
}

func TestPayment_CannotCaptureWithoutAuthorization(t *testing.T) {
	p, _ := NewPayment(uuid.New(), 1000, "USD")

	err := p.Capture()
	assert.ErrorIs(t, err, domainerrors.ErrPaymentNotAuthorized)
	assert.Equal(t, StatusCreated, p.Status) // unchanged
}

func TestPayment_CannotRefundWithoutCapture(t *testing.T) {
	p, _ := NewPayment(uuid.New(), 1000, "USD")
	_ = p.Authorize()

	err := p.Refund()
	assert.ErrorIs(t, err, domainerrors.ErrPaymentNotCaptured)
	assert.Equal(t, StatusAuthorized, p.Status) // unchanged
}

func TestPayment_CannotAuthorizeTwice(t *testing.T) {
	p, _ := NewPayment(uuid.New(), 1000, "USD")
	_ = p.Authorize()

	err := p.Authorize()
	assert.ErrorIs(t, err, domainerrors.ErrInvalidStateTransition)
}

func TestPayment_FailFromCreated(t *testing.T) {
	p, _ := NewPayment(uuid.New(), 1000, "USD")

	err := p.Fail("fraud check failed")
	assert.NoError(t, err)
	assert.Equal(t, StatusFailed, p.Status)
}

func TestPayment_CannotFailAfterCapture(t *testing.T) {
	p, _ := NewPayment(uuid.New(), 1000, "USD")
	_ = p.Authorize()
	_ = p.Capture()

	err := p.Fail("too late")
	assert.ErrorIs(t, err, domainerrors.ErrPaymentAlreadyTerminal)
	assert.Equal(t, StatusCaptured, p.Status) // unchanged
}

func TestPayment_CannotRefundTwice(t *testing.T) {
	p, _ := NewPayment(uuid.New(), 1000, "USD")
	_ = p.Authorize()
	_ = p.Capture()
	_ = p.Refund()

	err := p.Refund()
	assert.ErrorIs(t, err, domainerrors.ErrPaymentNotCaptured)
}
