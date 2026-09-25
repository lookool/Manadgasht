package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"manadgasht/internal/config"
	"manadgasht/internal/database"
	"manadgasht/internal/handlers"
	"manadgasht/internal/repository"
	"manadgasht/internal/router"
	"manadgasht/internal/service"
)

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	cfg := config.Load()
	db := database.Connect(cfg)

	// وایرینگ لایه‌ها: repository (دیتابیس خام) → service (منطق تجاری)
	// → handlers (HTTP). هر لایه فقط لایه‌ی زیر خودش رو می‌شناسه، و این
	// فایل تنها جایی است که همه رو به هم وصل می‌کنه.
	packageRepo := repository.NewPackageRepository(db)
	blogRepo := repository.NewBlogRepository(db)

	packageService := service.NewPackageService(packageRepo)
	blogService := service.NewBlogService(blogRepo)

	h := handlers.New(packageService, blogService)
	r := router.New(h)

	srv := &http.Server{
		Addr:         ":" + cfg.AppPort,
		Handler:      r,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		slog.Info("سرور مانادگشت در حال اجراست", "address", "http://localhost"+srv.Addr)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			slog.Error("سرور با خطا متوقف شد", "error", err)
			os.Exit(1)
		}
	}()

	// خاموشی آرام: با Ctrl+C یا سیگنال kill، درخواست‌های در حال پردازش
	// قبل از بستن کامل سرور فرصت تموم‌شدن پیدا می‌کنن.
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop

	slog.Info("در حال خاموش کردن سرور ...")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		slog.Error("خاموشی آرام ناموفق بود", "error", err)
	}
}
