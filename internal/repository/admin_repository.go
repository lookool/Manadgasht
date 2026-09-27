package repository

import (
	"gorm.io/gorm"

	"manadgasht/internal/models"
)

type AdminRepository interface {
	FindByUsername(username string) (*models.Admin, error)
	FindByID(id uint) (*models.Admin, error)
}

type adminRepository struct {
	db *gorm.DB
}

func NewAdminRepository(db *gorm.DB) AdminRepository {
	return &adminRepository{db: db}
}

func (r *adminRepository) FindByUsername(username string) (*models.Admin, error) {
	var admin models.Admin
	if err := r.db.Where("username = ?", username).First(&admin).Error; err != nil {
		return nil, err
	}
	return &admin, nil
}

func (r *adminRepository) FindByID(id uint) (*models.Admin, error) {
	var admin models.Admin
	if err := r.db.First(&admin, id).Error; err != nil {
		return nil, err
	}
	return &admin, nil
}