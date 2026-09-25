package handlers

import "net/http"

// About صفحه‌ی درباره‌ی ما را نشان می‌دهد.
func (h *Handlers) About(w http.ResponseWriter, r *http.Request) {
	render(w, "about.html", map[string]interface{}{
		"Title": "درباره ما",
	})
}

// Contact صفحه‌ی تماس با ما را نشان می‌دهد.
func (h *Handlers) Contact(w http.ResponseWriter, r *http.Request) {
	render(w, "contact.html", map[string]interface{}{
		"Title": "تماس با ما",
	})
}
