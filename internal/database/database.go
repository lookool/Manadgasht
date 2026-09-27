package database

import (
	"fmt"
	"log/slog"
	"os"

	"golang.org/x/crypto/bcrypt"
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
	seedDefaultAdmin(db, cfg)

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
		&models.Session{},
		&models.WithdrawalRequest{},
	)
}

// seedCommissionDefaults اگر تنظیمات کمیسیون هنوز در دیتابیس نباشد،
// مقادیر پیش‌فرض اولیه (طبق نمونه‌ای که خودتون گفتید: ۵٪ تور، ۲٪ هتل)
// را می‌سازد؛ این‌ها فقط نقطه‌ی شروع هستن و بعداً از پنل ادمین قابل
// تغییرن.
func seedCommissionDefaults(db *gorm.DB) {
	defaults := map[models.PackageType]float64{
		models.PackageTypeTour:      5,
		models.PackageTypeHotel:     2,
		models.PackageTypeTourHotel: 4,
	}

	for t, percentage := range defaults {
		var count int64
		db.Model(&models.CommissionSetting{}).Where("type = ?", t).Count(&count)
		if count == 0 {
			db.Create(&models.CommissionSetting{Type: t, Percentage: percentage})
		}
	}
}

// seedDefaultAdmin اگر هنوز هیچ ادمینی تو دیتابیس نباشه، یکی از روی
// ADMIN_USERNAME/ADMIN_PASSWORD تو .env می‌سازه. بعد از اولین ورود،
// حتماً رمزش رو از پنل عوض کن.
func seedDefaultAdmin(db *gorm.DB, cfg *config.Config) {
	var count int64
	db.Model(&models.Admin{}).Count(&count)
	if count > 0 {
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(cfg.AdminPassword), bcrypt.DefaultCost)
	if err != nil {
		slog.Error("ساخت رمز ادمین پیش‌فرض ناموفق بود", "error", err)
		return
	}

	admin := &models.Admin{
		FullName:     "مدیر سایت",
		Username:     cfg.AdminUsername,
		PasswordHash: string(hash),
	}
	if err := db.Create(admin).Error; err != nil {
		slog.Error("ساخت حساب ادمین پیش‌فرض ناموفق بود", "error", err)
		return
	}

	slog.Warn("یک حساب ادمین پیش‌فرض ساخته شد؛ حتماً بعد از ورود رمزش رو عوض کن", "username", cfg.AdminUsername)
}