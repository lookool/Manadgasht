package repository

import (
	"gorm.io/gorm"

	"manadgasht/internal/models"
)

type OrderRepository interface {
	Create(tx *gorm.DB, order *models.Order) error
	ListByUser(userID uint) ([]models.Order, error)
	ListByReferringAgent(agentID uint) ([]models.Order, error)
}

type orderRepository struct {
	db *gorm.DB
}

func NewOrderRepository(db *gorm.DB) OrderRepository {
	return &orderRepository{db: db}
}

func (r *orderRepository) Create(tx *gorm.DB, order *models.Order) error {
	return tx.Create(order).Error
}

func (r *orderRepository) ListByUser(userID uint) ([]models.Order, error) {
	var orders []models.Order
	err := r.db.Where("user_id = ?", userID).Order("created_at desc").Find(&orders).Error
	return orders, err
}

func (r *orderRepository) ListByReferringAgent(agentID uint) ([]models.Order, error) {
	var orders []models.Order
	err := r.db.Where("referring_agent_id = ?", agentID).Order("created_at desc").Find(&orders).Error
	return orders, err
}
