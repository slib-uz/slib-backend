package enum

type OTPPurpose string

const (
	OTPPurposeLogin         OTPPurpose = "OTP_LOGIN"
	OTPPurposeRegister      OTPPurpose = "OTP_REGISTER"
	OTPPurposeResetPassword OTPPurpose = "OTP_RESET_PASSWORD"
)
