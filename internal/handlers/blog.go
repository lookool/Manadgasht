package handlers

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

// BlogList همه‌ی مطالب منتشرشده‌ی وبلاگ را نشان می‌دهد.
func (h *Handlers) BlogList(w http.ResponseWriter, r *http.Request) {
	posts, err := h.Blog.List()
	if err != nil {
		logError("بارگذاری لیست بلاگ", err)
		http.Error(w, "خطای داخلی سرور", http.StatusInternalServerError)
		return
	}

	render(w, r, "blog.html", map[string]interface{}{
		"Title": "بلاگ",
		"Posts": posts,
	})
}

// BlogDetail یک پست خاص را بر اساس Slug نشان می‌دهد.
func (h *Handlers) BlogDetail(w http.ResponseWriter, r *http.Request) {
	slug := chi.URLParam(r, "slug")

	post, err := h.Blog.FindBySlug(slug)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	render(w, r, "blog_detail.html", map[string]interface{}{
		"Title": post.Title,
		"Post":  post,
	})
}
