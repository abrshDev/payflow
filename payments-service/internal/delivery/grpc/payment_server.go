package grpc

import (
	"context"
	"errors"

	paymentsv1 "github.com/abrshDev/payflow/proto/payments/v1"
	"github.com/abrshDev/payment-service/internal/app/payment/commands"
	"github.com/abrshDev/payment-service/internal/app/payment/queries"
	"github.com/abrshDev/payment-service/internal/domain/entities"
	domainerrors "github.com/abrshDev/payment-service/internal/domain/errors"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type PaymentServer struct {
	paymentsv1.UnimplementedPaymentServiceServer

	createHandler    *commands.CreatePaymentHandler
	authorizeHandler *commands.AuthorizePaymentHandler
	captureHandler   *commands.CapturePaymentHandler
	refundHandler    *commands.RefundPaymentHandler
	getHandler       *queries.GetPaymentHandler
}

func NewPaymentServer(
	create *commands.CreatePaymentHandler,
	authorize *commands.AuthorizePaymentHandler,
	capture *commands.CapturePaymentHandler,
	refund *commands.RefundPaymentHandler,
	get *queries.GetPaymentHandler,
) *PaymentServer {
	return &PaymentServer{
		createHandler:    create,
		authorizeHandler: authorize,
		captureHandler:   capture,
		refundHandler:    refund,
		getHandler:       get,
	}
}

func (s *PaymentServer) CreatePayment(ctx context.Context, req *paymentsv1.CreatePaymentRequest) (*paymentsv1.Payment, error) {
	merchantID, err := uuid.Parse(req.GetMerchantId())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid merchant_id")
	}

	payment, err := s.createHandler.Handle(ctx, commands.CreatePaymentCommand{
		MerchantID: merchantID,
		Amount:     req.GetAmount(),
		Currency:   req.GetCurrency(),
	})
	if err != nil {
		return nil, toStatus(err)
	}
	return toProto(payment), nil
}

func (s *PaymentServer) GetPayment(ctx context.Context, req *paymentsv1.GetPaymentRequest) (*paymentsv1.Payment, error) {
	id, err := parseID(req.GetId())
	if err != nil {
		return nil, err
	}
	return s.get(ctx, id)
}

func (s *PaymentServer) AuthorizePayment(ctx context.Context, req *paymentsv1.AuthorizePaymentRequest) (*paymentsv1.Payment, error) {
	id, err := parseID(req.GetId())
	if err != nil {
		return nil, err
	}
	if err := s.authorizeHandler.Handle(ctx, commands.AuthorizePaymentCommand{PaymentID: id}); err != nil {
		return nil, toStatus(err)
	}
	return s.get(ctx, id)
}

func (s *PaymentServer) CapturePayment(ctx context.Context, req *paymentsv1.CapturePaymentRequest) (*paymentsv1.Payment, error) {
	id, err := parseID(req.GetId())
	if err != nil {
		return nil, err
	}
	if err := s.captureHandler.Handle(ctx, commands.CapturePaymentCommand{PaymentID: id}); err != nil {
		return nil, toStatus(err)
	}
	return s.get(ctx, id)
}

func (s *PaymentServer) RefundPayment(ctx context.Context, req *paymentsv1.RefundPaymentRequest) (*paymentsv1.Payment, error) {
	id, err := parseID(req.GetId())
	if err != nil {
		return nil, err
	}
	if err := s.refundHandler.Handle(ctx, commands.RefundPaymentCommand{PaymentID: id}); err != nil {
		return nil, toStatus(err)
	}
	return s.get(ctx, id)
}

// helpers

func (s *PaymentServer) get(ctx context.Context, id uuid.UUID) (*paymentsv1.Payment, error) {
	payment, err := s.getHandler.Handle(ctx, queries.GetPaymentQuery{PaymentID: id})
	if err != nil {
		return nil, status.Error(codes.NotFound, err.Error())
	}
	return toProto(payment), nil
}

func parseID(raw string) (uuid.UUID, error) {
	id, err := uuid.Parse(raw)
	if err != nil {
		return uuid.Nil, status.Error(codes.InvalidArgument, "invalid payment id")
	}
	return id, nil
}

func toStatus(err error) error {
	switch {
	case errors.Is(err, domainerrors.ErrInvalidAmount),
		errors.Is(err, domainerrors.ErrInvalidCurrency):
		return status.Error(codes.InvalidArgument, err.Error())
	case errors.Is(err, domainerrors.ErrInvalidStateTransition),
		errors.Is(err, domainerrors.ErrPaymentNotAuthorized),
		errors.Is(err, domainerrors.ErrPaymentNotCaptured),
		errors.Is(err, domainerrors.ErrPaymentAlreadyTerminal):
		return status.Error(codes.FailedPrecondition, err.Error())
	default:
		return status.Error(codes.Internal, err.Error())
	}
}

func toProto(p *entities.Payment) *paymentsv1.Payment {
	return &paymentsv1.Payment{
		Id:         p.ID.String(),
		MerchantId: p.MerchantID.String(),
		Amount:     p.Amount.Amount,
		Currency:   p.Amount.Currency,
		Status:     string(p.Status),
	}
}
