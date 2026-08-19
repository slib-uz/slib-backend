package service_test

import (
	"context"
	"errors"
	"strconv"
	"testing"
	"time"

	"slib.uz/src/core/application/service"
)

func TestRevokeStoresEntryWithTTLFromExp(t *testing.T) {
	c := newFakeCache()
	s := service.NewTokenRevocationService(c)

	exp := time.Now().Add(10 * time.Minute)
	if err := s.Revoke(context.Background(), "jti-1", exp); err != nil {
		t.Fatalf("Revoke xato qaytardi: %v", err)
	}

	ttl, ok := c.ttls["revoked:jti:jti-1"]
	if !ok {
		t.Fatal("denylist yozuvi yaratilmadi")
	}
	// TTL exp dan hisoblanadi, qat'iy raqamdan emas.
	if ttl > 10*time.Minute || ttl < 9*time.Minute {
		t.Fatalf("TTL exp dan hisoblanmadi: %v", ttl)
	}
}

func TestRevokeIgnoresAlreadyExpiredToken(t *testing.T) {
	c := newFakeCache()
	s := service.NewTokenRevocationService(c)

	err := s.Revoke(context.Background(), "jti-old", time.Now().Add(-time.Minute))
	if err != nil {
		t.Fatalf("Revoke xato qaytardi: %v", err)
	}
	if len(c.data) != 0 {
		t.Fatal("muddati o'tgan token uchun yozuv yaratildi")
	}
}

func TestRevokeReturnsErrorWhenCacheWriteFails(t *testing.T) {
	c := newFakeCache()
	c.setErr = errors.New("redis down")
	s := service.NewTokenRevocationService(c)

	err := s.Revoke(context.Background(), "jti-1", time.Now().Add(time.Minute))
	if err == nil {
		t.Fatal("yozish nosozligi yashirildi — logout jimgina muvaffaqiyatsiz bo'lardi")
	}
}

func TestRevokedAtReturnsNilForUnknownToken(t *testing.T) {
	s := service.NewTokenRevocationService(newFakeCache())

	at, err := s.RevokedAt(context.Background(), "jti-yoq")
	if err != nil {
		t.Fatalf("RevokedAt xato qaytardi: %v", err)
	}
	if at != nil {
		t.Fatal("bekor qilinmagan token bekor qilingan deb qaytdi")
	}
}

func TestRevokedAtReturnsRevocationTime(t *testing.T) {
	c := newFakeCache()
	s := service.NewTokenRevocationService(c)

	if err := s.Revoke(context.Background(), "jti-1", time.Now().Add(time.Minute)); err != nil {
		t.Fatalf("Revoke xato qaytardi: %v", err)
	}

	at, err := s.RevokedAt(context.Background(), "jti-1")
	if err != nil {
		t.Fatalf("RevokedAt xato qaytardi: %v", err)
	}
	if at == nil {
		t.Fatal("bekor qilingan token topilmadi")
	}
	if time.Since(at.At) > time.Minute {
		t.Fatalf("bekor qilingan vaqt noto'g'ri: %v", at.At)
	}
}

// Bekor qilishning ikki turi bir-biridan ajralishi shart: logout yakuniy,
// rotatsiya esa grace oynasiga bo'ysunadi. Ajratilmasa, logout qilingan refresh
// token grace oynasi ichida yangi juftlik olib bera olardi.
func TestRevokeMarksRevocationAsFinal(t *testing.T) {
	s := service.NewTokenRevocationService(newFakeCache())

	if err := s.Revoke(context.Background(), "jti-1", time.Now().Add(time.Minute)); err != nil {
		t.Fatalf("Revoke xato qaytardi: %v", err)
	}

	rev, err := s.RevokedAt(context.Background(), "jti-1")
	if err != nil {
		t.Fatalf("RevokedAt xato qaytardi: %v", err)
	}
	if rev == nil || !rev.Final {
		t.Fatalf("logout bekor qilishi yakuniy deb belgilanmadi: %+v", rev)
	}
}

func TestRevokeWithGraceMarksRevocationAsNonFinal(t *testing.T) {
	s := service.NewTokenRevocationService(newFakeCache())

	if err := s.RevokeWithGrace(context.Background(), "jti-1", time.Now().Add(time.Minute)); err != nil {
		t.Fatalf("RevokeWithGrace xato qaytardi: %v", err)
	}

	rev, err := s.RevokedAt(context.Background(), "jti-1")
	if err != nil {
		t.Fatalf("RevokedAt xato qaytardi: %v", err)
	}
	if rev == nil {
		t.Fatal("bekor qilingan token topilmadi")
	}
	if rev.Final {
		t.Fatal("rotatsiya bekor qilishi yakuniy deb belgilandi — grace oynasi yo'qoladi")
	}
}

// Migratsiya yo'q: belgi qo'shilishidan oldin yozilgan qiymatlar ham o'qilishi
// kerak. Ular rotatsiya deb qabul qilinadi — ya'ni avvalgi xatti-harakat saqlanadi.
func TestRevokedAtTreatsUnmarkedLegacyValueAsRotation(t *testing.T) {
	c := newFakeCache()
	c.data[service.RevokedKeyPrefix+"jti-1"] = strconv.FormatInt(time.Now().Unix(), 10)
	s := service.NewTokenRevocationService(c)

	rev, err := s.RevokedAt(context.Background(), "jti-1")
	if err != nil {
		t.Fatalf("RevokedAt xato qaytardi: %v", err)
	}
	if rev == nil {
		t.Fatal("belgisiz eski qiymat o'qilmadi")
	}
	if rev.Final {
		t.Fatal("belgisiz eski qiymat yakuniy deb o'qildi")
	}
	if time.Since(rev.At) > time.Minute {
		t.Fatalf("bekor qilingan vaqt noto'g'ri: %v", rev.At)
	}
}

func TestRevokedAtTreatsCorruptValueAsFinalRevocation(t *testing.T) {
	c := newFakeCache()
	c.data[service.RevokedKeyPrefix+"jti-1"] = "axlat:qiymat"
	s := service.NewTokenRevocationService(c)

	rev, err := s.RevokedAt(context.Background(), "jti-1")
	if err != nil {
		t.Fatalf("RevokedAt xato qaytardi: %v", err)
	}
	if rev == nil || !rev.Final {
		t.Fatalf("buzuq qiymat yakuniy bekor qilish deb hisoblanmadi: %+v", rev)
	}
}

// jti bo'sh bo'lishi mumkin: Task 1 dan oldin berilgan tokenlarda "jti" da'vosi
// yo'q, Decode ularni bo'sh ID bilan qaytaradi. Bo'sh jti uchun yozuv yaratilsa,
// kalit yalang'och prefiksga ("revoked:jti:") aylanadi va barcha eski tokenlar
// bitta kalitni bo'lishadi — bitta logout hammasini o'chirgan bo'lardi.
func TestRevokeIgnoresEmptyJti(t *testing.T) {
	c := newFakeCache()
	s := service.NewTokenRevocationService(c)

	if err := s.Revoke(context.Background(), "", time.Now().Add(time.Minute)); err != nil {
		t.Fatalf("Revoke xato qaytardi: %v", err)
	}
	if len(c.data) != 0 {
		t.Fatalf("bo'sh jti uchun yozuv yaratildi: %v", c.data)
	}
}

// Yalang'och prefiks kaliti mavjud bo'lsa ham, bo'sh jti hech qachon
// "bekor qilingan" deb hisoblanmasligi kerak.
func TestRevokedAtIgnoresEmptyJtiEvenIfPrefixKeyExists(t *testing.T) {
	c := newFakeCache()
	c.data["revoked:jti:"] = "1"
	s := service.NewTokenRevocationService(c)

	at, err := s.RevokedAt(context.Background(), "")
	if err != nil {
		t.Fatalf("RevokedAt xato qaytardi: %v", err)
	}
	if at != nil {
		t.Fatal("bo'sh jti bekor qilingan deb qaytdi — eski tokenlar ommaviy o'chirilardi")
	}
}

func TestRevokedAtFailsOpenWhenCacheUnavailable(t *testing.T) {
	c := newFakeCache()
	c.getErr = errors.New("redis down")
	s := service.NewTokenRevocationService(c)

	at, err := s.RevokedAt(context.Background(), "jti-1")
	if err != nil {
		t.Fatalf("fail-open buzildi: xato qaytdi: %v", err)
	}
	if at != nil {
		t.Fatal("fail-open buzildi: token bekor qilingan deb qaytdi")
	}
}
