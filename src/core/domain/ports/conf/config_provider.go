package conf

type ConfigAdapter interface {
	GetReviewDeadlineDays() int
	GetFrontendURL() string
	GetROIFrontendURL() string
	OtpTTLMinutes() int
	OtpVerifyMaxAttempts() int
	OtpSendPerMinute() int
	OtpSendPerHour() int
	GetJwtAccessTokenExpireMinutes() int
	GetJwtRefreshTokenExpireMinutes() int
	GetCrossRefSenderEmail() string
	GetClientBasicAuthCredentials() (string, string)
	IsRefreshRotationStrict() bool
	GetRefreshRotationGraceSeconds() int
	IsProduction() bool
}
