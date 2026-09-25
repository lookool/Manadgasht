package handlers

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

// PackagesList همه‌ی تور/هتل‌های فعال را نشان می‌دهد؛ با پارامترهای
// اختیاری ?region=domestic|international (برای دو منوی هدر) و
// ?q= (برای کادر جست‌وجوی صفحه‌ی اصلی).
func (h *Handlers) PackagesList(w http.ResponseWriter, r *http.Request) {
	region := r.URL.Query().Get("region")
	packageType := r.URL.Query().Get("type")
	query := r.URL.Query().Get("q")

	packages, err := h.Packages.List(region, packageType, query)
	if err != nil {
		logError("بارگذاری لیست تور/هتل", err)
		http.Error(w, "خطای داخلی سرور", http.StatusInternalServerError)
		return
	}

	title := "تور و هتل"
	switch region {
	case "domestic":
		title = "تورهای داخلی"
	case "international":
		title = "تورهای خارجی"
	}
	switch packageType {
	case "hotel":
		title = "هتل"
	case "tour_hotel":
		title = "پکیج تور + هتل"
	}

	render(w, "packages.html", map[string]interface{}{
		"Title":    title,
		"Packages": packages,
		"Region":   region,
		"Query":    query,
	})
}

// PackageDetail یک تور/هتل خاص را بر اساس Slug نشان می‌دهد.
func (h *Handlers) PackageDetail(w http.ResponseWriter, r *http.Request) {
	slug := chi.URLParam(r, "slug")

	pkg, err := h.Packages.FindBySlug(slug)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	render(w, "package_detail.html", map[string]interface{}{
		"Title":   pkg.Title,
		"Package": pkg,
	})
}
