package handlers

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"time"

	"manadgasht/internal/middleware"
)

// funcMap توابع کمکی قابل استفاده داخل فایل‌های html/template.
var funcMap = template.FuncMap{
	"formatPrice": func(price float64) string {
		s := fmt.Sprintf("%.0f", price)
		n := len(s)
		if n <= 3 {
			return s
		}
		var out []byte
		for i, c := range []byte(s) {
			if i > 0 && (n-i)%3 == 0 {
				out = append(out, ',')
			}
			out = append(out, c)
		}
		return string(out)
	},
}

var pages = map[string]*template.Template{}
var adminPages = map[string]*template.Template{}

func init() {
	pageFiles := []string{
		"home.html",
		"packages.html",
		"package_detail.html",
		"blog.html",
		"blog_detail.html",
		"about.html",
		"contact.html",
		"user_register.html",
		"user_login.html",
		"account.html",
		"agent_register.html",
		"agent_login.html",
		"agent_dashboard.html",
	}

	base := template.Must(template.New("base.html").Funcs(funcMap).ParseFiles(
		"web/templates/base.html",
		"web/templates/partials/header.html",
		"web/templates/partials/footer.html",
	))

	for _, page := range pageFiles {
		clone, err := base.Clone()
		if err != nil {
			log.Fatalf("خطا در کلون تمپلیت برای %s: %v", page, err)
		}
		clone = template.Must(clone.ParseFiles("web/templates/" + page))
		pages[page] = clone
	}

	// پنل ادمین لایوت جدای خودش رو داره (بدون هدر/فوتر عمومی سایت،
	// با یک سایدبار)، پس یک base دوم و یک map جدا براش می‌سازیم.
	adminPageFiles := []string{
		"admin_login.html",
		"admin_dashboard.html",
		"admin_packages_list.html",
		"admin_package_form.html",
		"admin_blog_list.html",
		"admin_blog_form.html",
	}

	adminBase := template.Must(template.New("admin_base.html").Funcs(funcMap).ParseFiles(
		"web/templates/admin_base.html",
		"web/templates/partials/admin_sidebar.html",
	))

	for _, page := range adminPageFiles {
		clone, err := adminBase.Clone()
		if err != nil {
			log.Fatalf("خطا در کلون تمپلیت ادمین برای %s: %v", page, err)
		}
		clone = template.Must(clone.ParseFiles("web/templates/" + page))
		adminPages[page] = clone
	}
}

// render یک صفحه‌ی خاص را به همراه لایوت اصلی (base) رندر می‌کند.
func render(w http.ResponseWriter, r *http.Request, name string, data map[string]interface{}) {
	tmpl, ok := pages[name]
	if !ok {
		log.Println("تمپلیت پیدا نشد:", name)
		http.Error(w, "خطای داخلی سرور", http.StatusInternalServerError)
		return
	}

	if data == nil {
		data = map[string]interface{}{}
	}
	data["CurrentYear"] = time.Now().Year()

	if _, ok := data["Region"]; !ok {
		data["Region"] = ""
	}
	if _, ok := data["CurrentUser"]; !ok {
		data["CurrentUser"] = middleware.CurrentUser(r)
	}
	if _, ok := data["CurrentAgent"]; !ok {
		data["CurrentAgent"] = middleware.CurrentAgent(r)
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	if err := tmpl.ExecuteTemplate(w, "base", data); err != nil {
		log.Println("خطا در اجرای تمپلیت:", err)
		http.Error(w, "خطای داخلی سرور", http.StatusInternalServerError)
	}
}

// renderAdmin مشابه render ولی برای صفحات پنل ادمین (لایوت جدا، بدون
// هدر/فوتر عمومی سایت).
func renderAdmin(w http.ResponseWriter, r *http.Request, name string, data map[string]interface{}) {
	tmpl, ok := adminPages[name]
	if !ok {
		log.Println("تمپلیت ادمین پیدا نشد:", name)
		http.Error(w, "خطای داخلی سرور", http.StatusInternalServerError)
		return
	}

	if data == nil {
		data = map[string]interface{}{}
	}
	data["CurrentYear"] = time.Now().Year()
	data["CurrentAdmin"] = middleware.CurrentAdmin(r)

	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	if err := tmpl.ExecuteTemplate(w, "admin_base", data); err != nil {
		log.Println("خطا در اجرای تمپلیت ادمین:", err)
		http.Error(w, "خطای داخلی سرور", http.StatusInternalServerError)
	}
}