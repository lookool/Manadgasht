package database

import (
	"fmt"
	"log/slog"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"manadgasht/internal/config"
	"manadgasht/internal/models"
)

// Connect به پستگرس وصل می‌شود و جدول‌ها را (در صورت نبودن) می‌سازد.
func Connect(cfg *config.Config) *gorm.DB {
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName, cfg.DBSSLMode,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		slog.Error("اتصال به دیتابیس ناموفق بود", "error", err)
		os.Exit(1)
	}

	if err := autoMigrate(db); err != nil {
		slog.Error("مایگریشن دیتابیس ناموفق بود", "error", err)
		os.Exit(1)
	}

	seedCommissionDefaults(db)

	slog.Info("اتصال به دیتابیس با موفقیت انجام شد")
	return db
}

func autoMigrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&models.User{},
		&models.Admin{},
		&models.Agent{},
		&models.Package{},
		&models.CommissionSetting{},
		&models.Order{},
		&models.BlogPost{},
	)
}

// seedCommissionDefaults اگر تنظیمات کمیسیون هنوز در دیتابیس نباشد،
// مقادیر پیش‌فرض (صفر درصد) را می‌سازد تا ادمین بعداً از پنل مقدارشان
// را عوض کند.
func seedCommissionDefaults(db *gorm.DB) {
	defaults := []models.PackageType{
		models.PackageTypeTour,
		models.PackageTypeHotel,
		models.PackageTypeTourHotel,
	}

	for _, t := range defaults {
		var count int64
		db.Model(&models.CommissionSetting{}).Where("type = ?", t).Count(&count)
		if count == 0 {
			db.Create(&models.CommissionSetting{Type: t, Percentage: 0})
		}
	}
}
