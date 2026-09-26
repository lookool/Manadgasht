package repository

import (
	"time"

	"gorm.io/gorm"

	"manadgasht/internal/models"
)

type SessionRepository interface {
	Create(session *models.Session) error
	FindValidByToken(token string) (*models.Session, error)
	Delete(token string) error
}

type sessionRepository struct {
	db *gorm.DB
}

func NewSessionRepository(db *gorm.DB) SessionRepository {
	return &sessionRepository{db: db}
}

func (r *sessionRepository) Create(session *models.Session) error {
	return r.db.Create(session).Error
}

// FindValidByToken فقط نشستی را برمی‌گرداند که هنوز منقضی نشده باشد.
func (r *sessionRepository) FindValidByToken(token string) (*models.Session, error) {
	var session models.Session
	err := r.db.Where("token = ? AND expires_at > ?", token, time.Now()).First(&session).Error
	if err != nil {
		return nil, err
	}
	return &session, nil
}

func (r *sessionRepository) Delete(token string) error {
	return r.db.Delete(&models.Session{}, "token = ?", token).Error
}
