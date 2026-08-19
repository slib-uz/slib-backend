package service

import (
	"context"
	"strconv"
	"strings"
	"time"

	"github.com/labstack/gommon/log"
	"slib.uz/src/core/domain/ports/cache"
)

// RevokedKeyPrefix — denylist kalitlari uchun Redis prefiksi.
//
// Eksport qilingan: bu format spetsifikatsiyada hujjatlashtirilgan shartnoma
// (CWE-613 dizayn hujjati), va test/core/application/service paketidagi
// testlar soat siljishi (clock skew) kabi holatlarni simulyatsiya qilish
// uchun fake cache'ga to'g'ridan-to'g'ri shu prefiks bilan yozadi — Revoke
// orqali kelajakdagi vaqt belgisini hosil qilib bo'lmaydi.
const RevokedKeyPrefix = "revoked:jti:"

// finalMarker — qiymatga qo'shiladigan "bu bekor qilish yakuniy" belgisi.
// Qiymat formati: "<unix>" (rotatsiya) yoki "<unix>:final" (logout).
const finalMarker = "final"

// Revocation — denylist yozuvi: token qachon va qanday bekor qilingani.
//
// Bekor qilishning ikki turi bor va ular bir xil emas:
//
//   - Rotatsiya (Final == false): refresh aylantirilganda eski token bekor
//     qilinadi. Bu yerda grace oynasi kerak — parallel yuborilgan so'rovlar va
//     tarmoq qayta urinishlari o'sha tokenni bir necha soniya ichida ikkinchi
//     marta ko'rsatishi normal holat.
//   - Logout (Final == true): foydalanuvchi ataylab chiqdi. Bu yerda grace
//     bo'lishi mumkin emas — aks holda o'g'irlangan refresh token bilan
//     grace oynasi ichida yangi juftlik olib, chiqishni bekor qilib bo'lardi
//     (60 soniyalik oyna 7 kunlik yangi sessiyaga aylanardi).
type Revocation struct {
	At    time.Time
	Final bool
}

// TokenRevocationService bekor qilingan tokenlarning jti ro'yxatini boshqaradi.
//
// Kalit qiymati sifatida bekor qilingan vaqt (unix soniya) va bekor qilish
// turi saqlanadi. Vaqt refresh rotatsiyasidagi "grace" oynasini qo'shimcha
// kalitsiz beradi, tur esa grace qaysi yozuvlarga tegishli emasligini aytadi.
type TokenRevocationService struct {
	cache cache.CacheProvider
}

// @inject
func NewTokenRevocationService(c cache.CacheProvider) *TokenRevocationService {
	return &TokenRevocationService{cache: c}
}

// Revoke tokenni yakuniy ravishda denylist'ga yozadi: grace oynasi qo'llanilmaydi.
// Logout shu metoddan foydalanadi.
//
// TTL tokenning qolgan umriga teng: token o'z-o'zidan o'lgach yozuvni saqlashning
// ma'nosi yo'q. Muddati o'tgan token uchun hech narsa qilmaydi.
func (this *TokenRevocationService) Revoke(ctx context.Context, jti string, exp time.Time) error {
	return this.revoke(ctx, jti, exp, true)
}

// RevokeWithGrace tokenni grace oynasiga bo'ysunadigan qilib bekor qiladi.
// Faqat refresh rotatsiyasi uchun: qayta ishlatilgan token grace ichida hali ham
// qabul qilinadi. Yakuniy bekor qilish kerak bo'lsa Revoke ishlatiladi.
func (this *TokenRevocationService) RevokeWithGrace(ctx context.Context, jti string, exp time.Time) error {
	return this.revoke(ctx, jti, exp, false)
}

func (this *TokenRevocationService) revoke(ctx context.Context, jti string, exp time.Time, final bool) error {
	if jti == "" {
		return nil
	}

	ttl := time.Until(exp)
	if ttl <= 0 {
		return nil
	}

	value := strconv.FormatInt(time.Now().Unix(), 10)
	if final {
		value += ":" + finalMarker
	}

	return this.cache.Set(ctx, RevokedKeyPrefix+jti, value, ttl)
}

// RevokedAt tokenning bekor qilingan vaqtini va turini qaytaradi.
// (nil, nil) — token bekor qilinmagan.
//
// Redis yetib bo'lmasa ham (nil, nil) qaytaradi va ERROR loglaydi: bu ongli
// fail-open qarori — Redis yiqilishi butun saytni to'xtatmasligi kerak.
// Xavf oynasi access token TTL si bilan chegaralanadi.
func (this *TokenRevocationService) RevokedAt(ctx context.Context, jti string) (*Revocation, error) {
	if jti == "" {
		return nil, nil
	}

	raw, err := this.cache.GetByKey(ctx, RevokedKeyPrefix+jti)
	if err != nil {
		log.Error("TokenRevocationService.RevokedAt: denylist o'qib bo'lmadi, fail-open. jti=", jti, " err=", err.Error())
		return nil, nil
	}

	if raw == "" {
		return nil, nil
	}

	// Belgisiz qiymat — rotatsiya. Migratsiya yo'q: bosqichma-bosqich deploy
	// paytida eski instance yozgan yozuvlar ham shu tarmoqdan o'tadi.
	timestamp, marker, hasMarker := strings.Cut(raw, ":")
	final := hasMarker && marker == finalMarker

	seconds, err := strconv.ParseInt(timestamp, 10, 64)
	if err != nil || (hasMarker && !final) {
		// Buzuq qiymat — xavfsiz tomonga og'amiz: token yakuniy bekor qilingan
		// deb hisoblanadi, ya'ni grace ham yordam bermaydi.
		log.Error("TokenRevocationService.RevokedAt: yaroqsiz qiymat, token bekor qilingan deb hisoblanadi. jti=", jti, " value=", raw)
		return &Revocation{At: time.Unix(0, 0), Final: true}, nil
	}

	return &Revocation{At: time.Unix(seconds, 0), Final: final}, nil
}
