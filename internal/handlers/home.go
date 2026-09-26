package handlers

import "net/http"

// Home صفحه‌ی اصلی سایت را نمایش می‌دهد: هیرو + جست‌وجو، چند تور/هتل
// ویژه، دسته‌های تور، و چند مطلب اخیر بلاگ.
func (h *Handlers) Home(w http.ResponseWriter, r *http.Request) {
	packages, err := h.Packages.Featured(6)
	if err != nil {
		logError("بارگذاری تورهای ویژه", err)
	}

	posts, err := h.Blog.Recent(3)
	if err != nil {
		logError("بارگذاری مطالب اخیر بلاگ", err)
	}

	render(w, r, "home.html", map[string]interface{}{
		"Title":    "خانه",
		"Packages": packages,
		"Posts":    posts,
	})
}
