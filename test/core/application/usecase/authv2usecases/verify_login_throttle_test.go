package authv2usecases_test

import (
	"context"
	"errors"
	"testing"

	"slib.uz/src/core/application/response"
	"slib.uz/src/core/application/usecase/authv2usecases"
)

// Throttle bloklaganda OTP Check chaqirilmaydi va 429 qaytadi.
func TestVerifyBlockedReturnsError(t *testing.T) {
	// blockingThrottle() send_otp_throttle_test.go dagi yordamchi.
	// Boshqa bog'liqliklar nil — blok yo'lida chaqirilmaydi.
	uc := authv2usecases.NewVerifyAndLoginUseCase(nil, nil, nil, nil, nil, blockingThrottle())

	_, _, err := uc.Execute(context.Background(), "sess", "123456", "", "")

	if !errors.Is(err, response.TooManyAttemptsError) {
		t.Fatalf("TooManyAttemptsError kutilgandi, %v keldi", err)
	}
}
