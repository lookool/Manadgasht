package service

import (
	"strings"

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

// List لیست تور/هتل‌ها را با فیلتر منطقه، نوع و جست‌وجوی متنی برمی‌گرداند.
func (s *PackageService) List(region, packageType, query string) ([]models.Package, error) {
	return s.repo.List(region, packageType, query)
}

func (s *PackageService) FindBySlug(slug string) (*models.Package, error) {
	return s.repo.FindBySlug(slug)
}

// ---------- پنل ادمین ----------

// PackageAdminInput داده‌ی خام فرم ادمین برای ساخت/ویرایش یک تور/هتل.
type PackageAdminInput struct {
	Title           string
	Slug            string
	Type            models.PackageType
	Description     string
	Location        string
	DurationDays    int
	Price           float64
	IsInternational bool
	IsActive        bool
	CoverImage      string // مسیر عکسی که آپلود شده (یا خالی = بدون تغییر)
	Gallery         string
}

func (s *PackageService) ListAll() ([]models.Package, error) {
	return s.repo.ListAll()
}

func (s *PackageService) FindByID(id uint) (*models.Package, error) {
	return s.repo.FindByID(id)
}

func (in PackageAdminInput) validate() *ValidationError {
	verr := newValidationError()
	if strings.TrimSpace(in.Title) == "" {
		verr.add("Title", "عنوان را وارد کنید")
	}
	if strings.TrimSpace(in.Slug) == "" {
		verr.add("Slug", "نامک (Slug) را وارد کنید")
	}
	switch in.Type {
	case models.PackageTypeTour, models.PackageTypeHotel, models.PackageTypeTourHotel:
	default:
		verr.add("Type", "نوع پکیج معتبر نیست")
	}
	if in.Price <= 0 {
		verr.add("Price", "قیمت باید بزرگ‌تر از صفر باشد")
	}
	return verr
}

func (s *PackageService) CreatePackage(in PackageAdminInput) (*models.Package, error) {
	verr := in.validate()
	if verr.hasErrors() {
		return nil, verr
	}

	pkg := &models.Package{
		Title:           strings.TrimSpace(in.Title),
		Slug:            strings.TrimSpace(in.Slug),
		Type:            in.Type,
		Description:     in.Description,
		Location:        in.Location,
		DurationDays:    in.DurationDays,
		Price:           in.Price,
		IsInternational: in.IsInternational,
		IsActive:        in.IsActive,
		CoverImage:      in.CoverImage,
		Gallery:         in.Gallery,
	}

	if err := s.repo.Create(pkg); err != nil {
		if isDuplicateKeyError(err) {
			verr.add("Slug", "این نامک (Slug) قبلاً استفاده شده")
			return nil, verr
		}
		return nil, err
	}
	return pkg, nil
}

func (s *PackageService) UpdatePackage(id uint, in PackageAdminInput) (*models.Package, error) {
	verr := in.validate()
	if verr.hasErrors() {
		return nil, verr
	}

	pkg, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}

	pkg.Title = strings.TrimSpace(in.Title)
	pkg.Slug = strings.TrimSpace(in.Slug)
	pkg.Type = in.Type
	pkg.Description = in.Description
	pkg.Location = in.Location
	pkg.DurationDays = in.DurationDays
	pkg.Price = in.Price
	pkg.IsInternational = in.IsInternational
	pkg.IsActive = in.IsActive
	pkg.Gallery = in.Gallery
	if in.CoverImage != "" {
		pkg.CoverImage = in.CoverImage
	}

	if err := s.repo.Update(pkg); err != nil {
		if isDuplicateKeyError(err) {
			verr.add("Slug", "این نامک (Slug) قبلاً استفاده شده")
			return nil, verr
		}
		return nil, err
	}
	return pkg, nil
}

func (s *PackageService) DeletePackage(id uint) error {
	return s.repo.Delete(id)
}