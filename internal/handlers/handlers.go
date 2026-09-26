package handlers

import "manadgasht/internal/service"

// Handlers تمام هندلرهای HTTP سایت را نگه می‌دارد. هندلرها فقط با
// سرویس‌ها کار می‌کنند و هیچ‌وقت مستقیم به دیتابیس دسترسی ندارند.
type Handlers struct {
	Packages *service.PackageService
	Blog     *service.BlogService
	Auth     *service.AuthService
	Orders   *service.OrderService
	Wallet   *service.WalletService
}

func New(
	packages *service.PackageService,
	blog *service.BlogService,
	auth *service.AuthService,
	orders *service.OrderService,
	wallet *service.WalletService,
) *Handlers {
	return &Handlers{
		Packages: packages,
		Blog:     blog,
		Auth:     auth,
		Orders:   orders,
		Wallet:   wallet,
	}
}
