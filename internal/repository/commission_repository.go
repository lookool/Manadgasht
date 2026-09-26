package repository

import (
	"gorm.io/gorm"

	"manadgasht/internal/models"
)

type CommissionRepository interface {
	PercentageFor(packageType models.PackageType) (float64, error)
}

type commissionRepository struct {
	db *gorm.DB
}

func NewCommissionRepository(db *gorm.DB) CommissionRepository {
	return &commissionRepository{db: db}
}

// PercentageFor درصد کمیسیون تنظیم‌شده برای یک نوع پکیج را برمی‌گرداند.
// اگر هنوز ردیفی برایش ساخته نشده باشد (حالتی که نباید پیش بیاید،
// چون در seedCommissionDefaults همه‌شان ساخته می‌شوند)، صفر برمی‌گردد.
func (r *commissionRepository) PercentageFor(packageType models.PackageType) (float64, error) {
	var setting models.CommissionSetting
	err := r.db.Where("type = ?", packageType).First(&setting).Error
	if err != nil {
		return 0, err
	}
	return setting.Percentage, nil
}
