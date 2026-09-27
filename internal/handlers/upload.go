package handlers

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const (
	uploadDir       = "web/static/uploads"
	maxUploadSizeMB = 10
)

// saveUploadedFile اگر کاربر فایلی برای فیلد fieldName انتخاب کرده
// باشد، آن را داخل web/static/uploads ذخیره می‌کند و مسیر قابل‌استفاده
// در تگ <img> را برمی‌گرداند. اگر فایلی انتخاب نشده باشه، رشته‌ی خالی
// و nil برمی‌گردونه (خطا نیست، چون این فیلد همیشه اختیاریه).
func saveUploadedFile(r *http.Request, fieldName string) (string, error) {
	if err := r.ParseMultipartForm(maxUploadSizeMB << 20); err != nil {
		return "", nil
	}

	file, header, err := r.FormFile(fieldName)
	if err != nil {
		return "", nil
	}
	defer file.Close()

	if err := os.MkdirAll(uploadDir, 0o755); err != nil {
		return "", err
	}

	ext := strings.ToLower(filepath.Ext(header.Filename))
	if ext == "" {
		ext = ".jpg"
	}
	filename := fmt.Sprintf("%d%s", time.Now().UnixNano(), ext)
	destPath := filepath.Join(uploadDir, filename)

	dest, err := os.Create(destPath)
	if err != nil {
		return "", err
	}
	defer dest.Close()

	if _, err := io.Copy(dest, file); err != nil {
		return "", err
	}

	return "/static/uploads/" + filename, nil
}