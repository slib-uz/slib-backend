package security_test

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"encoding/pem"
	"os"
	"path/filepath"
	"testing"
	"time"

	"slib.uz/src/core/domain/entity"
	"slib.uz/src/core/domain/ports/security"
	"slib.uz/src/infrastructure/config"
	infrasecurity "slib.uz/src/infrastructure/security"
)

// newTestService haqiqiy ECDSA kalit jufti yaratadi, uni PEM fayllarga yozadi
// va servisni haqiqiy konstruktor orqali quradi.
//
// Test endi alohida paketda (security_test), shuning uchun JwtTokenService
// ning eksport qilinmagan maydonlariga to'g'ridan-to'g'ri kira olmaydi. Struct
// literalidan foydalanish o'rniga NewJwtTokenService orqali quriladi — bu
// haqiqiy ishlab chiqarish yo'lini (PEM fayldan o'qish) ham sinaydi.
func newTestService(t *testing.T) security.TokenService {
	t.Helper()

	key, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		t.Fatalf("test kaliti yaratilmadi: %v", err)
	}

	dir := t.TempDir()
	privPath := filepath.Join(dir, "private.pem")
	pubPath := filepath.Join(dir, "public.pem")

	privBytes, err := x509.MarshalECPrivateKey(key)
	if err != nil {
		t.Fatalf("private key marshal qilinmadi: %v", err)
	}
	privPEM := pem.EncodeToMemory(&pem.Block{Type: "EC PRIVATE KEY", Bytes: privBytes})
	if err := os.WriteFile(privPath, privPEM, 0o600); err != nil {
		t.Fatalf("private key yozilmadi: %v", err)
	}

	pubBytes, err := x509.MarshalPKIXPublicKey(&key.PublicKey)
	if err != nil {
		t.Fatalf("public key marshal qilinmadi: %v", err)
	}
	pubPEM := pem.EncodeToMemory(&pem.Block{Type: "PUBLIC KEY", Bytes: pubBytes})
	if err := os.WriteFile(pubPath, pubPEM, 0o600); err != nil {
		t.Fatalf("public key yozilmadi: %v", err)
	}

	cfg := &config.Config{
		JwtPrivateKeyPath: privPath,
		JwtPublicKeyPath:  pubPath,
	}

	return infrasecurity.NewJwtTokenService(cfg)
}

func newTestClaim() *entity.TokenEntity {
	return entity.NewTokenEntity(
		time.Now().Add(time.Minute),
		"42",
		map[string]any{"type": "ACCESS"},
	)
}

func TestDecodePreservesJTI(t *testing.T) {
	s := newTestService(t)

	signed := s.Encode(newTestClaim())
	if signed == "" {
		t.Fatal("Encode bo'sh token qaytardi")
	}

	decoded, err := s.Decode(signed)
	if err != nil {
		t.Fatalf("Decode xato qaytardi: %v", err)
	}
	if decoded.ID == "" {
		t.Fatal("Decode jti ni yo'qotdi: TokenEntity.ID bo'sh")
	}
}

func TestEachTokenGetsDistinctJTI(t *testing.T) {
	s := newTestService(t)

	first, err := s.Decode(s.Encode(newTestClaim()))
	if err != nil {
		t.Fatalf("birinchi Decode xato: %v", err)
	}
	second, err := s.Decode(s.Encode(newTestClaim()))
	if err != nil {
		t.Fatalf("ikkinchi Decode xato: %v", err)
	}

	if first.ID == second.ID {
		t.Fatalf("ikki xil token bir xil jti oldi: %s", first.ID)
	}
}
