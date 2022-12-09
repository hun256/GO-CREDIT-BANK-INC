package sql

import (
	"context"
	"fmt"

	"github.com/imrenagi/go-payment"
	"github.com/imrenagi/go-payment/gateway/midtrans"
	"gorm.io/gorm"
	"github.com/rs/zerolog"
)

func NewMidtransTransactionRepository(db *gorm.DB) *MidtransTransactionRepository {
	return &MidtransTransactionRepository{
		DB: db,
	}
}

type MidtransTransactionRepository struct {
	DB *gorm.DB
}

// Save will update the notification stored in sql database
func (m *MidtransTransactionRepository) Save(ctx context.Context, status *midtrans.TransactionStatus) error {
	log := zerolog.Ctx(ctx).With().Str("function", "MidtransTransactionRepository.Save").Logger()

	if err := m.DB.Save(status).Find(&status).Error; err != nil {
		log.Error().Err(err).Msgf("cant save midtrans transaction status")
		return payment.ErrDatabase
	}
	return nil
}

// FindByOrderID fetch a transaction status for a given orderID
func (m *MidtransTransactionRepository) FindByOrderID(ctx context.Context, orderID string) (*midtrans.TransactionStatus, error) {
	log := zerolog.Ctx(ctx).With().Str("function", "MidtransTransactionRepository.FindByOrderID").Logger()

	var status midtrans.TransactionStatus
	req := m.DB.
		Where("order_id = ?", orderID).
		First(&status)

	if req.Error == gorm.ErrRecordNotFound {
		return nil, fmt.Errorf("payment status for order %s %w", orderID, payment.ErrNotFound)
	}
	if req.Error != nil {
		log.Error().Err(req.Error).Msg("cant find midtrans transaction status")
		return nil, payment.ErrDatabase
	}

	return &status, nil

}
