package handlers

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"time"
)

// funcMap توابع کمکی قابل استفاده داخل فایل‌های html/template.
var funcMap = template.FuncMap{
	// formatPrice عدد قیمت را با جداکننده‌ی هزارگان چاپ می‌کند،
	// مثلاً ۱۲۰۰۰۰۰۰ → 12,000,000
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

// هر صفحه یک فایل "content" جدا دارد؛ چون همه از همین نام "content"
// استفاده می‌کنند، نمی‌شود همه را با هم در یک template.Template پارس کرد
// (باعث تداخل/پنیک می‌شود). به همین دلیل برای هر صفحه، یک نسخه‌ی جدا از
// base+header+footer کلون می‌کنیم و فقط همان یک فایل content را رویش
// اضافه می‌کنیم؛ نتیجه در pages کش می‌شود تا هر بار از نو پارس نشود.
var pages = map[string]*template.Template{}

func init() {
	pageFiles := []string{
		"home.html",
		"packages.html",
		"package_detail.html",
		"blog.html",
		"blog_detail.html",
		"about.html",
		"contact.html",
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
}

// render یک صفحه‌ی خاص را به همراه لایوت اصلی (base) رندر می‌کند.
// data باید یک map باشد که داده‌های مورد نیاز صفحه را در بر دارد؛
// CurrentYear همیشه خودکار اضافه می‌شود.
func render(w http.ResponseWriter, name string, data map[string]interface{}) {
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

	// Region همیشه باید مقدار داشته باشه (حتی خالی)، چون تمپلیت هدر
	// برای مشخص‌کردن آیتم فعال منو، همیشه با eq مقایسه‌ش می‌کنه.
	if _, ok := data["Region"]; !ok {
		data["Region"] = ""
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	if err := tmpl.ExecuteTemplate(w, "base", data); err != nil {
		log.Println("خطا در اجرای تمپلیت:", err)
		http.Error(w, "خطای داخلی سرور", http.StatusInternalServerError)
	}
}
