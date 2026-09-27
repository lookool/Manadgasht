package service

import (
	"strings"

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

// ---------- پنل ادمین ----------

type BlogAdminInput struct {
	Title       string
	Slug        string
	Summary     string
	Content     string
	IsPublished bool
	CoverImage  string
}

func (in BlogAdminInput) validate() *ValidationError {
	verr := newValidationError()
	if strings.TrimSpace(in.Title) == "" {
		verr.add("Title", "عنوان را وارد کنید")
	}
	if strings.TrimSpace(in.Slug) == "" {
		verr.add("Slug", "نامک (Slug) را وارد کنید")
	}
	if strings.TrimSpace(in.Content) == "" {
		verr.add("Content", "متن پست را وارد کنید")
	}
	return verr
}

func (s *BlogService) ListAll() ([]models.BlogPost, error) {
	return s.repo.ListAll()
}

func (s *BlogService) FindByID(id uint) (*models.BlogPost, error) {
	return s.repo.FindByID(id)
}

func (s *BlogService) CreatePost(in BlogAdminInput) (*models.BlogPost, error) {
	verr := in.validate()
	if verr.hasErrors() {
		return nil, verr
	}

	post := &models.BlogPost{
		Title:       strings.TrimSpace(in.Title),
		Slug:        strings.TrimSpace(in.Slug),
		Summary:     in.Summary,
		Content:     in.Content,
		IsPublished: in.IsPublished,
		CoverImage:  in.CoverImage,
	}

	if err := s.repo.Create(post); err != nil {
		if isDuplicateKeyError(err) {
			verr.add("Slug", "این نامک (Slug) قبلاً استفاده شده")
			return nil, verr
		}
		return nil, err
	}
	return post, nil
}

func (s *BlogService) UpdatePost(id uint, in BlogAdminInput) (*models.BlogPost, error) {
	verr := in.validate()
	if verr.hasErrors() {
		return nil, verr
	}

	post, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}

	post.Title = strings.TrimSpace(in.Title)
	post.Slug = strings.TrimSpace(in.Slug)
	post.Summary = in.Summary
	post.Content = in.Content
	post.IsPublished = in.IsPublished
	if in.CoverImage != "" {
		post.CoverImage = in.CoverImage
	}

	if err := s.repo.Update(post); err != nil {
		if isDuplicateKeyError(err) {
			verr.add("Slug", "این نامک (Slug) قبلاً استفاده شده")
			return nil, verr
		}
		return nil, err
	}
	return post, nil
}

func (s *BlogService) DeletePost(id uint) error {
	return s.repo.Delete(id)
}