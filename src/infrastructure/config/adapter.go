package config

import (
	"slib.uz/src/core/domain/ports/conf"
)

type ConfigAdapterImpl struct {
	env *Config
}

// @inject
func NewConfigAdapter(env *Config) conf.ConfigAdapter {
	return &ConfigAdapterImpl{env: env}
}

func (this *ConfigAdapterImpl) OtpTTLMinutes() int {
	return this.env.OtpTtlMinutes
}

func (this *ConfigAdapterImpl) OtpVerifyMaxAttempts() int {
	return this.env.OtpVerifyMaxAttempts
}

func (this *ConfigAdapterImpl) OtpSendPerMinute() int {
	return this.env.OtpSendPerMinute
}

func (this *ConfigAdapterImpl) OtpSendPerHour() int {
	return this.env.OtpSendPerHour
}

func (this *ConfigAdapterImpl) GetReviewDeadlineDays() int {
	return this.env.ReviewDeadlineDays
}

func (this *ConfigAdapterImpl) GetFrontendURL() string {
	return this.env.FrontendURL
}

func (this *ConfigAdapterImpl) GetROIFrontendURL() string {
	return this.env.ROIFrontendURL
}

func (this *ConfigAdapterImpl) GetJwtAccessTokenExpireMinutes() int {
	return this.env.JwtAccessTokenExpireMinutes
}

func (this *ConfigAdapterImpl) GetJwtRefreshTokenExpireMinutes() int {
	return this.env.JwtRefreshTokenExpireMinutes
}

func (this *ConfigAdapterImpl) GetCrossRefSenderEmail() string {
	return this.env.CrossRefSenderEmail
}

func (this *ConfigAdapterImpl) GetClientBasicAuthCredentials() (string, string) {
	return this.env.ClientBasicAuthUsername, this.env.ClientBasicAuthSecret
}

func (this *ConfigAdapterImpl) IsRefreshRotationStrict() bool {
	return this.env.RefreshRotationStrict
}

func (this *ConfigAdapterImpl) GetRefreshRotationGraceSeconds() int {
	return this.env.RefreshRotationGraceSeconds
}

func (this *ConfigAdapterImpl) IsProduction() bool {
	return this.env.Production
}
