package main

import (
	"math/rand"
	"time"
)

// GenerateKey - Tasodifiy 6 xonali kalit yaratish funksiyasi
func GenerateKey() string {
	const charset = "abcdefghijklmnopqrstuvwxyz0123456789"
	// Har safar har xil natija olish uchun yangi seed
	seed := rand.New(rand.NewSource(time.Now().UnixNano()))

	b := make([]byte, 6)
	for i := range b {
		b[i] = charset[seed.Intn(len(charset))]
	}
	return string(b)
}

// IsValidURL - Kiritilgan URL xavfsiz ekanligini tekshirish (Oddiy misol)
func IsValidURL(url string) bool {
	// Bu yerda qo'shimcha tekshiruvlar qo'shish mumkin
	return len(url) > 10
}
