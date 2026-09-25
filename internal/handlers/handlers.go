package handlers

import "manadgasht/internal/service"

// Handlers تمام هندلرهای HTTP سایت را نگه می‌دارد. هندلرها فقط با
// سرویس‌ها کار می‌کنند و هیچ‌وقت مستقیم به دیتابیس دسترسی ندارند؛
// این باعث می‌شه هم تست کردن ساده‌تر باشه، هم تغییر منطق تجاری بدون
// دست‌زدن به کد HTTP ممکن باشه.
type Handlers struct {
	Packages *service.PackageService
	Blog     *service.BlogService
}

func New(packages *service.PackageService, blog *service.BlogService) *Handlers {
	return &Handlers{
		Packages: packages,
		Blog:     blog,
	}
}
