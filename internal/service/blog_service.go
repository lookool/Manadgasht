package service

import (
	"manadgasht/internal/models"
	"manadgasht/internal/repository"
)

type BlogService struct {
	repo repository.BlogRepository
}

func NewBlogService(repo repository.BlogRepository) *BlogService {
	return &BlogService{repo: repo}
}

func (s *BlogService) Recent(limit int) ([]models.BlogPost, error) {
	return s.repo.Recent(limit)
}

func (s *BlogService) List() ([]models.BlogPost, error) {
	return s.repo.List()
}

func (s *BlogService) FindBySlug(slug string) (*models.BlogPost, error) {
	return s.repo.FindBySlug(slug)
}
