package handlers

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"

	"manadgasht/internal/service"
)

func (h *Handlers) AdminBlogList(w http.ResponseWriter, r *http.Request) {
	posts, err := h.Blog.ListAll()
	if err != nil {
		logError("بارگذاری لیست ادمین بلاگ", err)
		http.Error(w, "خطای داخلی سرور", http.StatusInternalServerError)
		return
	}

	renderAdmin(w, r, "admin_blog_list.html", map[string]interface{}{
		"Title": "مدیریت بلاگ",
		"Posts": posts,
	})
}

func (h *Handlers) AdminBlogNewForm(w http.ResponseWriter, r *http.Request) {
	renderAdmin(w, r, "admin_blog_form.html", map[string]interface{}{
		"Title":    "پست جدید",
		"IsEdit":   false,
		"FormPath": "/admin/blog",
	})
}

func (h *Handlers) AdminBlogCreate(w http.ResponseWriter, r *http.Request) {
	input, err := parseBlogForm(r)
	if err != nil {
		http.Error(w, "درخواست نامعتبر است", http.StatusBadRequest)
		return
	}

	_, err = h.Blog.CreatePost(input)
	if err != nil {
		var verr *service.ValidationError
		if errors.As(err, &verr) {
			renderAdmin(w, r, "admin_blog_form.html", map[string]interface{}{
				"Title":    "پست جدید",
				"IsEdit":   false,
				"FormPath": "/admin/blog",
				"Errors":   verr.Fields,
				"Form":     input,
			})
			return
		}
		logError("ساخت پست بلاگ", err)
		http.Error(w, "خطای داخلی سرور", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/admin/blog", http.StatusSeeOther)
}

func (h *Handlers) AdminBlogEditForm(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseUint(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	post, err := h.Blog.FindByID(uint(id))
	if err != nil {
		http.NotFound(w, r)
		return
	}

	renderAdmin(w, r, "admin_blog_form.html", map[string]interface{}{
		"Title":    "ویرایش پست",
		"IsEdit":   true,
		"FormPath": "/admin/blog/" + chi.URLParam(r, "id") + "/update",
		"Post":     post,
	})
}

func (h *Handlers) AdminBlogUpdate(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseUint(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	input, err := parseBlogForm(r)
	if err != nil {
		http.Error(w, "درخواست نامعتبر است", http.StatusBadRequest)
		return
	}

	_, err = h.Blog.UpdatePost(uint(id), input)
	if err != nil {
		var verr *service.ValidationError
		if errors.As(err, &verr) {
			renderAdmin(w, r, "admin_blog_form.html", map[string]interface{}{
				"Title":    "ویرایش پست",
				"IsEdit":   true,
				"FormPath": "/admin/blog/" + chi.URLParam(r, "id") + "/update",
				"Errors":   verr.Fields,
				"Post":     input,
			})
			return
		}
		logError("به‌روزرسانی پست بلاگ", err)
		http.Error(w, "خطای داخلی سرور", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/admin/blog", http.StatusSeeOther)
}

func (h *Handlers) AdminBlogDelete(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseUint(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	if err := h.Blog.DeletePost(uint(id)); err != nil {
		logError("حذف پست بلاگ", err)
	}

	http.Redirect(w, r, "/admin/blog", http.StatusSeeOther)
}

func parseBlogForm(r *http.Request) (service.BlogAdminInput, error) {
	coverImage, err := saveUploadedFile(r, "cover_image")
	if err != nil {
		return service.BlogAdminInput{}, err
	}

	if err := r.ParseMultipartForm(maxUploadSizeMB << 20); err != nil {
		return service.BlogAdminInput{}, err
	}

	return service.BlogAdminInput{
		Title:       r.FormValue("title"),
		Slug:        r.FormValue("slug"),
		Summary:     r.FormValue("summary"),
		Content:     r.FormValue("content"),
		IsPublished: r.FormValue("is_published") == "on",
		CoverImage:  coverImage,
	}, nil
}