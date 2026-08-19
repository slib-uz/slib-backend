package entity

import "time"

type TokenEntity struct {
	// ID — JWT "jti" da'vosi. Encode paytida generatsiya qilinadi,
	// Decode paytida to'ldiriladi. Bekor qilish shunga tayanadi.
	ID      string
	Exp     time.Time
	Subject string
	Payload map[string]any
}

func NewTokenEntity(exp time.Time, subject string, payload map[string]any) *TokenEntity {
	return &TokenEntity{Exp: exp, Subject: subject, Payload: payload}
}
