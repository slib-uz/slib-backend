package service

import (
	"context"
	"crypto/rand"
	"errors"
	"fmt"
	"math/big"
	"time"

	"slib.uz/src/core/application/response"
	"slib.uz/src/core/domain/entity"
	"slib.uz/src/core/domain/entity/enum"
	"slib.uz/src/core/domain/ports/conf"
	"slib.uz/src/core/domain/ports/repository"
	"slib.uz/src/core/utils"
)

type OTPService struct {
	repository repository.OTPCodeRepository
	conf       conf.ConfigAdapter
}

// @inject
func NewOTPService(repository repository.OTPCodeRepository, conf conf.ConfigAdapter) *OTPService {
	return &OTPService{repository: repository, conf: conf}
}

func (this *OTPService) Make(ctx context.Context, phoneNumber string, purpose enum.OTPPurpose) (*entity.OTPCodeEntity, error) {
	ttl := time.Duration(this.conf.OtpTTLMinutes()) * time.Minute
	expiresAt := time.Now().Add(ttl)

	otp := entity.NewOTPCodeEntity(0, phoneNumber, this.generateRandomCode(), utils.RandomHex(32), purpose, expiresAt, nil, time.Now())

	id, err := this.repository.Create(ctx, otp)
	if err != nil {
		return nil, err
	}

	otp.ID = id
	return otp, nil
}

func (this *OTPService) Check(ctx context.Context, code string, sessionID string) (*entity.OTPCodeEntity, error) {
	var invalidCodeError = func(message string) error {
		return response.NewOptionalResponse(200, response.CodeInvalidOrReusedCode, nil, message)
	}
	var expiredCodeError = response.NewOptionalResponse(200, response.CodeExpired, nil, "")

	otp, err := this.repository.GetByCodeAndSessionID(ctx, sessionID, code)

	if errors.Is(err, response.NotFoundError) {
		return nil, invalidCodeError("invalid code (not found)")
	} else if err != nil {
		return nil, err
	}

	if otp.ExpiresAt.Before(time.Now()) {
		return nil, expiredCodeError
	}

	if otp.UsedAt != nil {
		return nil, invalidCodeError("invalid code (already used)")
	}

	if err := this.repository.MarkAsUsed(ctx, otp.ID); err != nil {
		return nil, err
	}

	return otp, nil
}

func (*OTPService) generateRandomCode() string {
	n, _ := rand.Int(rand.Reader, big.NewInt(1000000))
	return fmt.Sprintf("%06d", n.Int64())
}

// func (*OTPService) generateSessionID() string {
// 	return utils.RandomHex(32)
// }
