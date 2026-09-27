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
	"manadgasht/internal/middleware"
	"manadgasht/internal/repository"
	"manadgasht/internal/router"
	"manadgasht/internal/service"
)

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	cfg := config.Load()
	db := database.Connect(cfg)

	packageRepo := repository.NewPackageRepository(db)
	blogRepo := repository.NewBlogRepository(db)
	userRepo := repository.NewUserRepository(db)
	agentRepo := repository.NewAgentRepository(db)
	adminRepo := repository.NewAdminRepository(db)
	sessionRepo := repository.NewSessionRepository(db)
	orderRepo := repository.NewOrderRepository(db)
	withdrawalRepo := repository.NewWithdrawalRepository(db)
	commissionRepo := repository.NewCommissionRepository(db)

	packageService := service.NewPackageService(packageRepo)
	blogService := service.NewBlogService(blogRepo)
	authService := service.NewAuthService(userRepo, agentRepo, adminRepo, sessionRepo)
	orderService := service.NewOrderService(db, packageRepo, orderRepo, agentRepo, userRepo, commissionRepo)
	walletService := service.NewWalletService(db, agentRepo, withdrawalRepo)

	h := handlers.New(packageService, blogService, authService, orderService, walletService)
	authMiddleware := middleware.NewAuth(authService)
	r := router.New(h, authMiddleware)

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