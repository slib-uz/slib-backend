package config

import (
	"os"

	"github.com/caarlos0/env/v10"
	"github.com/joho/godotenv"
)

type Config struct {

	// Database
	DBUser     string `env:"DB_USER,required"`
	DBPassword string `env:"DB_PASSWORD,required"`
	DBName     string `env:"DB_NAME,required"`
	DBHost     string `env:"DB_HOST,required"`
	DBPort     string `env:"DB_PORT,required"`
	DBSSLMode  string `env:"DB_SSL_MODE,required"`

	// Redis
	RedisAddress string `env:"REDIS_ADDRESS,required"`

	// CORS
	CorsAllowOrigin []string `env:"CORS_ALLOW_ORIGIN,required" envSeparator:","`

	// Auth
	JwtPrivateKeyPath            string `env:"JWT_PRIVATE_KEY_PATH,required"`
	JwtPublicKeyPath             string `env:"JWT_PUBLIC_KEY_PATH,required"`
	JwtAccessTokenExpireMinutes  int    `env:"JWT_ACCESS_TOKEN_EXPIRE_MINUTES,required"`
	JwtRefreshTokenExpireMinutes int    `env:"JWT_REFRESH_TOKEN_EXPIRE_MINUTES,required"`

	// Sessiya. Bu uchtasi required EMAS: env/.env git'da kuzatilmaydi,
	// required qilinsa mavjud muhitlar ishga tushmay qoladi.
	RefreshRotationStrict       bool `env:"REFRESH_ROTATION_STRICT" envDefault:"false"`
	RefreshRotationGraceSeconds int  `env:"REFRESH_ROTATION_GRACE_SECONDS" envDefault:"60"`

	// Production default true — xavfsiz tomonga og'ish.
	// O'rnatilmagan bo'lsa sandbox login o'chiq bo'ladi.
	Production bool `env:"PRODUCTION" envDefault:"true"`

	// ScienceID
	ScienceIDClientID     string `env:"SCIENCE_ID_CLIENT_ID,required"`
	ScienceIDClientSecret string `env:"SCIENCE_ID_CLIENT_SECRET,required"`
	ScienceIDBaseURL      string `env:"SCIENCE_ID_BASE_URL,required"`
	OAuthCodeVerifier     string `env:"OAUTH_CODE_VERIFIER,required"`

	// Kafka
	KafkaBrokers  string `env:"KAFKA_BROKERS,required"`
	KafkaTopic    string `env:"KAFKA_TOPIC,required"`
	KafkaUsername string `env:"KAFKA_USERNAME,required"`
	KafkaPassword string `env:"KAFKA_PASSWORD,required"`

	// Minio
	MinioEndpoint       string `env:"MINIO_ENDPOINT,required"`
	MinioAccessKey      string `env:"MINIO_ACCESS_KEY,required"`
	MinioSecretKey      string `env:"MINIO_SECRET_KEY,required"`
	MinioUseSSL         bool   `env:"MINIO_USE_HTTPS,required"`
	MinioBucketName     string `env:"MAIN_BUCKET_NAME,required"`
	PublicBucketName    string `env:"PUBLIC_BUCKET_NAME,required"`
	UploadedFileMaxSize int    `env:"UPLOADED_FILE_MAX_SIZE,required"`

	//// Payme
	//PaymeID  string `env:"PAYME_ID,required"`
	//PaymeKey string `env:"PAYME_KEY,required"`
	//PaymeURL string `env:"PAYME_URL,required"`

	// Literacy
	LiteracyBaseURL string `env:"LITERACY_BASE_URL,required"`
	LiteracyAPIKey  string `env:"LITERACY_API_KEY,required"`

	// ETaqriz
	ETaqrizBaseURL string `env:"ETAQRIZ_BASE_URL,required"`
	ETaqrizToken   string `env:"ETAQRIZ_TOKEN,required"`
	ETaqrizSecret  string `env:"ETAQRIZ_SECRET,required"`

	// AntiPlag
	AntiPlagBaseURL      string `env:"ANTIPLAG_BASE_URL,required"`
	AntiPlagClientID     string `env:"ANTIPLAG_CLIENT_ID,required"`
	AntiPlagClientSecret string `env:"ANTIPLAG_CLIENT_SECRET,required"`

	// Limits
	ViewsCountLifetimeMinute int `env:"VIEWS_COUNT_LIFETIME_MINUTES,required"`
	ReviewDeadlineDays       int `env:"REVIEW_DEADLINE_DAYS,required"`
	ResubmitDeadlineDays     int `env:"RESUBMIT_DEADLINE_DAYS,required"`

	// Telegram
	TelegramBotToken    string `env:"TELEGRAM_BOT_TOKEN,required"`
	TelegramAdminChatID string `env:"TELEGRAM_ADMIN_CHAT_ID,required"`

	// ROI
	ROIClientID     string `env:"ROI_CLIENT_ID,required"`
	ROIClientSecret string `env:"ROI_CLIENT_SECRET,required"`
	ROIBaseURL      string `env:"ROI_BASE_URL,required"`
	ROIFrontendURL  string `env:"ROI_FRONTEND_URL,required"`

	// Frontend
	FrontendURL string `env:"FRONTEND_URL,required"`

	// ORCID
	ORCIDClientID     string `env:"ORCID_CLIENT_ID,required"`
	ORCIDClientSecret string `env:"ORCID_CLIENT_SECRET,required"`
	ORCIDBaseURL      string `env:"ORCID_BASE_URL,required"`

	// PDF
	Html2PdfBaseURL string `env:"HTML2PDF_BASE_URL,required"`

	// OTP
	OtpTtlMinutes int `env:"OTP_TTL_MINUTES,required"`

	// Brute-force himoyasi (CWE-307): bu uchtasi required EMAS — env/.env
	// git'da kuzatilmaydi, required qilinsa mavjud deploy ishga tushmay
	// qoladi.
	OtpVerifyMaxAttempts int `env:"OTP_VERIFY_MAX_ATTEMPTS" envDefault:"5"`
	OtpSendPerMinute     int `env:"OTP_SEND_PER_MINUTE" envDefault:"1"`
	OtpSendPerHour       int `env:"OTP_SEND_PER_HOUR" envDefault:"5"`

	// SMS
	SmsEtcBaseURL  string `env:"SMS_ETC_BASE_URL,required"`
	SmsEtcUsername string `env:"SMS_ETC_USERNAME,required"`
	SmsEtcPassword string `env:"SMS_ETC_PASSWORD,required"`
	SmsAlphaName   string `env:"SMS_ETC_ALPHA_NAME,required"`

	// Basic Auth Credentials
	ClientBasicAuthUsername string `env:"CLIENT_USERNAME,required"`
	ClientBasicAuthSecret   string `env:"CLIENT_SECRET,required"`

	// Developer Auth
	DeveloperUsername     string `env:"DEVELOPER_USERNAME,required"`
	DeveloperPasswordHash string `env:"DEVELOPER_PASSWORD_HASH,required"`

	// Slib Billing
	SlibBillingBaseURL string `env:"SLIB_BILLING_BASE_URL,required"`

	// CrossRef
	CrossRefDepositURL         string `env:"CROSSREF_DEPOSIT_URL,required"`
	CrossRefDepositorEmail     string `env:"CROSSREF_DEPOSITOR_EMAIL,required"`
	CrossRefArticleResourceURL string `env:"CROSSREF_ARTICLE_RESOURCE_URL" envDefault:""`
	CrossRefSenderEmail        string `env:"CROSSREF_SENDER_EMAIL" envDefault:"admin@crossref.org"`
	CrossRefLoginURL           string `env:"CROSSREF_LOGIN_URL" envDefault:"https://doi.crossref.org/servlet/login"`

	// UzSci
	UzSciBaseURL string `env:"UZSCI_BASE_URL,required"`

	// DeepSeek
	DeepSeekAPIKey  string `env:"DEEPSEEK_API_KEY,required"`
	DeepSeekBaseURL string `env:"DEEPSEEK_BASE_URL,required"`
	DeepSeekModel   string `env:"DEEPSEEK_MODEL" envDefault:"deepseek-v4-pro"`

	// OpenRouter
	OpenRouterAPIKey  string `env:"OPENROUTER_API_KEY,required"`
	OpenRouterBaseURL string `env:"OPENROUTER_BASE_URL,required"`
}

// @inject
func NewConfig() *Config {
	return loadConfig()
}

func loadConfig() *Config {
	loadEnv()

	cfg := &Config{}
	if err := env.Parse(cfg); err != nil {
		panic(err)
	}

	return cfg
}

func loadEnv() {
	if os.Getenv("CONTAINER_MODE") != "1" {
		_ = godotenv.Load("env/.env.local")
		_ = godotenv.Load("env/.env")
	}
}
