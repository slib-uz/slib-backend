// Package ratelimit auth urinishlarini cheklash uchun atomik hisoblagich
// abstraksiyasini beradi. Implementatsiya (Redis) infratuzilma qatlamida.
package ratelimit

import (
	"context"
	"time"
)

type RateLimiter interface {
	// Hit kalitni atomik oshiradi va oshirilgandan keyingi joriy hisobni
	// qaytaradi. Kalit yangi bo'lsa (birinchi hit), window muddatiga TTL
	// o'rnatiladi.
	Hit(ctx context.Context, key string, window time.Duration) (int64, error)
	// Reset kalitni o'chiradi (muvaffaqiyatli login'dan keyin).
	Reset(ctx context.Context, key string) error
}
