package handlers

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"

	"manadgasht/internal/models"
	"manadgasht/internal/service"
)

func (h *Handlers) AdminPackagesList(w http.ResponseWriter, r *http.Request) {
	packages, err := h.Packages.ListAll()
	if err != nil {
		logError("بارگذاری لیست ادمین تور/هتل", err)
		http.Error(w, "خطای داخلی سرور", http.StatusInternalServerError)
		return
	}

	renderAdmin(w, r, "admin_packages_list.html", map[string]interface{}{
		"Title":    "مدیریت تور/هتل",
		"Packages": packages,
	})
}

func (h *Handlers) AdminPackageNewForm(w http.ResponseWriter, r *http.Request) {
	renderAdmin(w, r, "admin_package_form.html", map[string]interface{}{
		"Title":    "تور/هتل جدید",
		"IsEdit":   false,
		"FormPath": "/admin/packages",
	})
}

func (h *Handlers) AdminPackageCreate(w http.ResponseWriter, r *http.Request) {
	input, err := parsePackageForm(r)
	if err != nil {
		http.Error(w, "درخواست نامعتبر است", http.StatusBadRequest)
		return
	}

	_, err = h.Packages.CreatePackage(input)
	if err != nil {
		var verr *service.ValidationError
		if errors.As(err, &verr) {
			renderAdmin(w, r, "admin_package_form.html", map[string]interface{}{
				"Title":    "تور/هتل جدید",
				"IsEdit":   false,
				"FormPath": "/admin/packages",
				"Errors":   verr.Fields,
				"Form":     input,
			})
			return
		}
		logError("ساخت تور/هتل", err)
		http.Error(w, "خطای داخلی سرور", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/admin/packages", http.StatusSeeOther)
}

func (h *Handlers) AdminPackageEditForm(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseUint(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	pkg, err := h.Packages.FindByID(uint(id))
	if err != nil {
		http.NotFound(w, r)
		return
	}

	renderAdmin(w, r, "admin_package_form.html", map[string]interface{}{
		"Title":    "ویرایش تور/هتل",
		"IsEdit":   true,
		"FormPath": "/admin/packages/" + chi.URLParam(r, "id") + "/update",
		"Package":  pkg,
	})
}

func (h *Handlers) AdminPackageUpdate(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseUint(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	input, err := parsePackageForm(r)
	if err != nil {
		http.Error(w, "درخواست نامعتبر است", http.StatusBadRequest)
		return
	}

	_, err = h.Packages.UpdatePackage(uint(id), input)
	if err != nil {
		var verr *service.ValidationError
		if errors.As(err, &verr) {
			renderAdmin(w, r, "admin_package_form.html", map[string]interface{}{
				"Title":    "ویرایش تور/هتل",
				"IsEdit":   true,
				"FormPath": "/admin/packages/" + chi.URLParam(r, "id") + "/update",
				"Errors":   verr.Fields,
				"Package":  input,
			})
			return
		}
		logError("به‌روزرسانی تور/هتل", err)
		http.Error(w, "خطای داخلی سرور", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/admin/packages", http.StatusSeeOther)
}

func (h *Handlers) AdminPackageDelete(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseUint(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	if err := h.Packages.DeletePackage(uint(id)); err != nil {
		logError("حذف تور/هتل", err)
	}

	http.Redirect(w, r, "/admin/packages", http.StatusSeeOther)
}

func parsePackageForm(r *http.Request) (service.PackageAdminInput, error) {
	coverImage, err := saveUploadedFile(r, "cover_image")
	if err != nil {
		return service.PackageAdminInput{}, err
	}

	if err := r.ParseMultipartForm(maxUploadSizeMB << 20); err != nil {
		return service.PackageAdminInput{}, err
	}

	price, _ := strconv.ParseFloat(r.FormValue("price"), 64)
	duration, _ := strconv.Atoi(r.FormValue("duration_days"))

	return service.PackageAdminInput{
		Title:           r.FormValue("title"),
		Slug:            r.FormValue("slug"),
		Type:            models.PackageType(r.FormValue("type")),
		Description:     r.FormValue("description"),
		Location:        r.FormValue("location"),
		DurationDays:    duration,
		Price:           price,
		IsInternational: r.FormValue("is_international") == "on",
		IsActive:        r.FormValue("is_active") == "on",
		CoverImage:      coverImage,
		Gallery:         r.FormValue("gallery"),
	}, nil
}