package service

import (
	"manadgasht/internal/models"
	"manadgasht/internal/repository"
)

// PackageService منطق تجاری مربوط به تور/هتل را نگه می‌دارد. هندلرها
// فقط با این سرویس صحبت می‌کنند، نه مستقیم با دیتابیس یا Repository.
type PackageService struct {
	repo repository.PackageRepository
}

func NewPackageService(repo repository.PackageRepository) *PackageService {
	return &PackageService{repo: repo}
}

// Featured چند تور/هتل ویژه برای نمایش در صفحه‌ی اصلی برمی‌گرداند.
func (s *PackageService) Featured(limit int) ([]models.Package, error) {
	return s.repo.Featured(limit)
}

// List لیست تور/هتل‌ها را با فیلتر منطقه (داخلی/خارجی)، نوع
// (تور/هتل/ترکیبی) و جست‌وجوی متنی برمی‌گرداند؛ همان چیزی که پشت
// منوی «تورهای داخلی/خارجی»، کاشی‌های دسته‌بندی و کادر جست‌وجو است.
func (s *PackageService) List(region, packageType, query string) ([]models.Package, error) {
	return s.repo.List(region, packageType, query)
}

func (s *PackageService) FindBySlug(slug string) (*models.Package, error) {
	return s.repo.FindBySlug(slug)
}
