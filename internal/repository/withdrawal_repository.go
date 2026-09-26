package repository

import (
	"gorm.io/gorm"

	"manadgasht/internal/models"
)

type WithdrawalRepository interface {
	Create(tx *gorm.DB, w *models.WithdrawalRequest) error
	ListByAgent(agentID uint) ([]models.WithdrawalRequest, error)
}

type withdrawalRepository struct {
	db *gorm.DB
}

func NewWithdrawalRepository(db *gorm.DB) WithdrawalRepository {
	return &withdrawalRepository{db: db}
}

func (r *withdrawalRepository) Create(tx *gorm.DB, w *models.WithdrawalRequest) error {
	return tx.Create(w).Error
}

func (r *withdrawalRepository) ListByAgent(agentID uint) ([]models.WithdrawalRequest, error) {
	var list []models.WithdrawalRequest
	err := r.db.Where("agent_id = ?", agentID).Order("created_at desc").Find(&list).Error
	return list, err
}
