package repository

import (
	"gorm.io/gorm"

	"manadgasht/internal/models"
)

// PackageRepository تنها لایه‌ای است که مستقیم با جدول Package در
// دیتابیس صحبت می‌کند. هیچ منطق تجاری اینجا نیست، فقط کوئری خام.
type PackageRepository interface {
	Featured(limit int) ([]models.Package, error)
	List(region, packageType, query string) ([]models.Package, error)
	FindBySlug(slug string) (*models.Package, error)
}

type packageRepository struct {
	db *gorm.DB
}

func NewPackageRepository(db *gorm.DB) PackageRepository {
	return &packageRepository{db: db}
}

func (r *packageRepository) Featured(limit int) ([]models.Package, error) {
	var packages []models.Package
	err := r.db.
		Where("is_active = ?", true).
		Order("created_at desc").
		Limit(limit).
		Find(&packages).Error
	return packages, err
}

// List همه‌ی تور/هتل‌های فعال را برمی‌گرداند. region می‌تواند خالی،
// "domestic" یا "international" باشد؛ packageType می‌تواند خالی یا یکی
// از مقادیر models.PackageType باشد؛ query هم یک جست‌وجوی ساده روی
// عنوان است (برای کادر جست‌وجوی صفحه‌ی اصلی).
func (r *packageRepository) List(region, packageType, query string) ([]models.Package, error) {
	db := r.db.Where("is_active = ?", true)

	switch region {
	case "domestic":
		db = db.Where("is_international = ?", false)
	case "international":
		db = db.Where("is_international = ?", true)
	}

	if packageType != "" {
		db = db.Where("type = ?", packageType)
	}

	if query != "" {
		db = db.Where("title ILIKE ?", "%"+query+"%")
	}

	var packages []models.Package
	err := db.Order("created_at desc").Find(&packages).Error
	return packages, err
}

func (r *packageRepository) FindBySlug(slug string) (*models.Package, error) {
	var pkg models.Package
	err := r.db.Where("slug = ? AND is_active = ?", slug, true).First(&pkg).Error
	if err != nil {
		return nil, err
	}
	return &pkg, nil
}
