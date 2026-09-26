package repository

import (
	"gorm.io/gorm"

	"manadgasht/internal/models"
)

type AgentRepository interface {
	Create(agent *models.Agent) error
	FindByEmail(email string) (*models.Agent, error)
	FindByReferralCode(code string) (*models.Agent, error)
	FindByID(id uint) (*models.Agent, error)
	ExistsByNationalID(nationalID string) (bool, error)
	ExistsByReferralCode(code string) (bool, error)

	// IncreaseBalance و DecreaseBalance با db.Model(...).Update انجام
	// می‌شوند (نه Save روی یک struct از قبل خوانده‌شده) تا در برابر
	// درخواست‌های هم‌زمان هم درست جمع/کم بشن. tx همیشه باید همون
	// تراکنشی باشه که بقیه‌ی عملیات مرتبط (مثلاً ساخت سفارش) هم توشن.
	IncreaseBalance(tx *gorm.DB, agentID uint, amount float64) error
	DecreaseBalance(tx *gorm.DB, agentID uint, amount float64) error
}

type agentRepository struct {
	db *gorm.DB
}

func NewAgentRepository(db *gorm.DB) AgentRepository {
	return &agentRepository{db: db}
}

func (r *agentRepository) Create(agent *models.Agent) error {
	return r.db.Create(agent).Error
}

func (r *agentRepository) FindByEmail(email string) (*models.Agent, error) {
	var agent models.Agent
	if err := r.db.Where("email = ?", email).First(&agent).Error; err != nil {
		return nil, err
	}
	return &agent, nil
}

func (r *agentRepository) FindByReferralCode(code string) (*models.Agent, error) {
	var agent models.Agent
	if err := r.db.Where("referral_code = ?", code).First(&agent).Error; err != nil {
		return nil, err
	}
	return &agent, nil
}

func (r *agentRepository) FindByID(id uint) (*models.Agent, error) {
	var agent models.Agent
	if err := r.db.First(&agent, id).Error; err != nil {
		return nil, err
	}
	return &agent, nil
}

func (r *agentRepository) ExistsByNationalID(nationalID string) (bool, error) {
	var count int64
	err := r.db.Model(&models.Agent{}).Where("national_id = ?", nationalID).Count(&count).Error
	return count > 0, err
}

func (r *agentRepository) ExistsByReferralCode(code string) (bool, error) {
	var count int64
	err := r.db.Model(&models.Agent{}).Where("referral_code = ?", code).Count(&count).Error
	return count > 0, err
}

func (r *agentRepository) IncreaseBalance(tx *gorm.DB, agentID uint, amount float64) error {
	return tx.Model(&models.Agent{}).
		Where("id = ?", agentID).
		Update("balance", gorm.Expr("balance + ?", amount)).Error
}

func (r *agentRepository) DecreaseBalance(tx *gorm.DB, agentID uint, amount float64) error {
	result := tx.Model(&models.Agent{}).
		Where("id = ? AND balance >= ?", agentID, amount).
		Update("balance", gorm.Expr("balance - ?", amount))
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return ErrInsufficientBalance
	}
	return nil
}
