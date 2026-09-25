package repository

import (
	"gorm.io/gorm"

	"manadgasht/internal/models"
)

type BlogRepository interface {
	Recent(limit int) ([]models.BlogPost, error)
	List() ([]models.BlogPost, error)
	FindBySlug(slug string) (*models.BlogPost, error)
}

type blogRepository struct {
	db *gorm.DB
}

func NewBlogRepository(db *gorm.DB) BlogRepository {
	return &blogRepository{db: db}
}

func (r *blogRepository) Recent(limit int) ([]models.BlogPost, error) {
	var posts []models.BlogPost
	err := r.db.
		Where("is_published = ?", true).
		Order("created_at desc").
		Limit(limit).
		Find(&posts).Error
	return posts, err
}

func (r *blogRepository) List() ([]models.BlogPost, error) {
	var posts []models.BlogPost
	err := r.db.
		Where("is_published = ?", true).
		Order("created_at desc").
		Find(&posts).Error
	return posts, err
}

func (r *blogRepository) FindBySlug(slug string) (*models.BlogPost, error) {
	var post models.BlogPost
	err := r.db.Where("slug = ? AND is_published = ?", slug, true).First(&post).Error
	if err != nil {
		return nil, err
	}
	return &post, nil
}
