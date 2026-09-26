package service

import (
	"errors"
	"strings"

	"gorm.io/gorm"

	"manadgasht/internal/models"
	"manadgasht/internal/repository"
)

type RequestWithdrawalInput struct {
	AgentID    uint
	Amount     float64
	CardNumber string
	IBAN       string
}

type WalletService struct {
	db          *gorm.DB
	agents      repository.AgentRepository
	withdrawals repository.WithdrawalRepository
}

func NewWalletService(db *gorm.DB, agents repository.AgentRepository, withdrawals repository.WithdrawalRepository) *WalletService {
	return &WalletService{db: db, agents: agents, withdrawals: withdrawals}
}

// RequestWithdrawal یک درخواست برداشت برای همکار ثبت می‌کند و بلافاصله
// مبلغ را از موجودی او کسر می‌کند (تا امکان خرج دوباره‌ی همان مبلغ
// وجود نداشته باشد). بررسی و واریز واقعی، بعداً از پنل ادمین انجام
// می‌شود؛ اگر رد بشه، باید مبلغ به موجودی برگرده (فاز پنل ادمین).
func (s *WalletService) RequestWithdrawal(in RequestWithdrawalInput) (*models.WithdrawalRequest, error) {
	verr := newValidationError()
	in.CardNumber = strings.TrimSpace(in.CardNumber)
	in.IBAN = strings.TrimSpace(strings.ToUpper(in.IBAN))

	if in.Amount <= 0 {
		verr.add("Amount", "مبلغ باید بزرگ‌تر از صفر باشد")
	}
	if in.CardNumber == "" {
		verr.add("CardNumber", "شماره کارت را وارد کنید")
	}
	if in.IBAN == "" {
		verr.add("IBAN", "شماره شبا را وارد کنید")
	}
	if verr.hasErrors() {
		return nil, verr
	}

	request := &models.WithdrawalRequest{
		AgentID:    in.AgentID,
		Amount:     in.Amount,
		CardNumber: in.CardNumber,
		IBAN:       in.IBAN,
		Status:     models.WithdrawalPending,
	}

	err := s.db.Transaction(func(tx *gorm.DB) error {
		if err := s.agents.DecreaseBalance(tx, in.AgentID, in.Amount); err != nil {
			return err
		}
		return s.withdrawals.Create(tx, request)
	})
	if err != nil {
		if errors.Is(err, repository.ErrInsufficientBalance) {
			verr.add("Amount", "موجودی کافی نیست")
			return nil, verr
		}
		return nil, err
	}

	return request, nil
}

func (s *WalletService) ListByAgent(agentID uint) ([]models.WithdrawalRequest, error) {
	return s.withdrawals.ListByAgent(agentID)
}
