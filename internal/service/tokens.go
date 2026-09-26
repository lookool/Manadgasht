package service

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"strings"
)

// generateToken یک رشته‌ی تصادفیِ امن (برای توکن نشست/کوکی) می‌سازد.
func generateToken(nBytes int) (string, error) {
	b := make([]byte, nBytes)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}

const referralAlphabet = "ABCDEFGHJKLMNPQRSTUVWXYZ23456789" // بدون حروف/ارقام شبیه به هم (O/0, I/1)

// generateReferralCode یک کد معرفِ کوتاه و خوانا برای همکار می‌سازد،
// مثلاً «MG-7K2N4X».
func generateReferralCode() (string, error) {
	b := make([]byte, 6)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	var sb strings.Builder
	for _, v := range b {
		sb.WriteByte(referralAlphabet[int(v)%len(referralAlphabet)])
	}
	return fmt.Sprintf("MG-%s", sb.String()), nil
}
