package configs

import (
	"log"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	// Core
	Environment       string `mapstructure:"ENVIRONMENT"`
	HTTPServerAddress string `mapstructure:"HTTP_SERVER_ADDRESS"`
	GRPCServerAddress string `mapstructure:"GRPC_SERVER_ADDRESS"`

	// Database
	DBSource string `mapstructure:"DB_SOURCE_HIAS"`

	// Redis
	RedisURL string `mapstructure:"REDIS_URL"`

	// AWS
	AWSRegion      string `mapstructure:"AWS_REGION"`
	AWSEndpointURL string `mapstructure:"AWS_ENDPOINT_URL"`

	// Cognito
	CognitoClientID     string `mapstructure:"COGNITO_CLIENT_ID"`
	CognitoClientSecret string `mapstructure:"COGNITO_CLIENT_SECRET"`
	CognitoRedirectURI  string `mapstructure:"COGNITO_REDIRECT_URI"`
	CognitoDomain       string `mapstructure:"COGNITO_DOMAIN"`
	CognitoUserPoolID   string `mapstructure:"COGNITO_USER_POOL_ID"`

	// Auth / PASETO
	TokenSymmetricKey    string        `mapstructure:"TOKEN_SYMMETRIC_KEY"`
	AccessTokenDuration  time.Duration `mapstructure:"ACCESS_TOKEN_DURATION"`
	RefreshTokenDuration time.Duration `mapstructure:"REFRESH_TOKEN_DURATION"`
	CookieDomain         string        `mapstructure:"COOKIE_DOMAIN"`

	// CORS
	AllowedOrigins string `mapstructure:"ALLOWED_ORIGINS"`
	DashboardURL   string `mapstructure:"DASHBOARD_URL"`

	// S3
	AWSS3Bucket       string   `mapstructure:"AWS_S3_BUCKET"`
	AWSS3Region       string   `mapstructure:"AWS_S3_REGION"`
	AWSS3CDNDomain    string   `mapstructure:"AWS_S3_CDN_DOMAIN"`
	AWSS3MaxFileSize  int      `mapstructure:"AWS_S3_MAX_FILE_SIZE"`
	AWSS3AllowedTypes []string `mapstructure:"AWS_S3_ALLOWED_TYPES"`

	// SQS Queue URLs
	AWSSQSDocumentProcessingQueueURL string `mapstructure:"AWS_SQS_DOCUMENT_PROCESSING_QUEUE_URL"`
	AWSSQSExtractionResultsQueueURL  string `mapstructure:"AWS_SQS_EXTRACTION_RESULTS_QUEUE_URL"`
	AWSSQSClaimProcessingQueueURL    string `mapstructure:"AWS_SQS_CLAIM_PROCESSING_QUEUE_URL"`
	AWSSQSPaymentEventsQueueURL      string `mapstructure:"AWS_SQS_PAYMENT_EVENTS_QUEUE_URL"`
	AWSSQSNotificationEventsQueueURL string `mapstructure:"AWS_SQS_NOTIFICATION_EVENTS_QUEUE_URL"`

	// Watermill / Consumer Configuration
	ConsumerType               string        `mapstructure:"CONSUMER_TYPE"`
	WatermillWaitTimeSeconds   int32         `mapstructure:"WATERMILL_WAIT_TIME_SECONDS"`
	WatermillMaxMessages       int           `mapstructure:"WATERMILL_MAX_MESSAGES"`
	WatermillVisibilityTimeout time.Duration `mapstructure:"WATERMILL_VISIBILITY_TIMEOUT"`
	QueueManagerType           string        `mapstructure:"QUEUE_MANAGER_TYPE"`

	// M-Pesa
	MpesaConsumerKey    string `mapstructure:"MPESA_CONSUMER_KEY"`
	MpesaConsumerSecret string `mapstructure:"MPESA_CONSUMER_SECRET"`
	MpesaPasskey        string `mapstructure:"MPESA_PASSKEY"`
	MpesaShortcode      string `mapstructure:"MPESA_SHORTCODE"`
	MpesaCallbackURL    string `mapstructure:"MPESA_CALLBACK_URL"`
	MpesaEnvironment    string `mapstructure:"MPESA_ENVIRONMENT"` // sandbox | production

	// IPRS
	IPRSBaseURL string `mapstructure:"IPRS_BASE_URL"`
	IPRSAPIKey  string `mapstructure:"IPRS_API_KEY"`

	// SMART / Slade360
	SMARTBaseURL      string `mapstructure:"SMART_BASE_URL"`
	SMARTAPIKey       string `mapstructure:"SMART_API_KEY"`
	SMARTAPISecret    string `mapstructure:"SMART_API_SECRET"`
	SMARTFacilityCode string `mapstructure:"SMART_FACILITY_CODE"`

	// Bank
	BankBaseURL   string `mapstructure:"BANK_BASE_URL"`
	BankAPIKey    string `mapstructure:"BANK_API_KEY"`
	BankAccountNo string `mapstructure:"BANK_ACCOUNT_NO"`

	// SMS (Africa's Talking)
	SMSAPIKey   string `mapstructure:"SMS_API_KEY"`
	SMSUsername string `mapstructure:"SMS_USERNAME"`
	SMSSenderID string `mapstructure:"SMS_SENDER_ID"`

	// SNS / SES
	SNSTopicArn          string `mapstructure:"SNS_TOPIC_ARN"`
	NotificationTopicARN string `mapstructure:"NOTIFICATION_TOPIC_ARN"`
	SESFromEmail         string `mapstructure:"SES_FROM_EMAIL"`

	// WebSocket Configuration
	WebSocketMaxConnectionsPerUser  int           `mapstructure:"WEBSOCKET_MAX_CONNECTIONS_PER_USER"`
	WebSocketMaxTotalConnections    int64         `mapstructure:"WEBSOCKET_MAX_TOTAL_CONNECTIONS"`
	WebSocketReadTimeout            time.Duration `mapstructure:"WEBSOCKET_READ_TIMEOUT"`
	WebSocketWriteTimeout           time.Duration `mapstructure:"WEBSOCKET_WRITE_TIMEOUT"`
	WebSocketPingInterval           time.Duration `mapstructure:"WEBSOCKET_PING_INTERVAL"`
	WebSocketPongTimeout            time.Duration `mapstructure:"WEBSOCKET_PONG_TIMEOUT"`
	WebSocketClientHeartbeatTimeout time.Duration `mapstructure:"WEBSOCKET_CLIENT_HEARTBEAT_TIMEOUT"`
	WebSocketEnableMetrics          bool          `mapstructure:"WEBSOCKET_ENABLE_METRICS"`
	WebSocketAllowedOrigins         string        `mapstructure:"WEBSOCKET_ALLOWED_ORIGINS"`

	// SSE Configuration
	SSEMaxConnections        int           `mapstructure:"SSE_MAX_CONNECTIONS"`
	SSEMaxConnectionsPerUser int           `mapstructure:"SSE_MAX_CONNECTIONS_PER_USER"`
	SSEReadTimeout           time.Duration `mapstructure:"SSE_READ_TIMEOUT"`
	SSEWriteTimeout          time.Duration `mapstructure:"SSE_WRITE_TIMEOUT"`
	SSEIdleTimeout           time.Duration `mapstructure:"SSE_IDLE_TIMEOUT"`
	SSEEnableMetrics         bool          `mapstructure:"SSE_ENABLE_METRICS"`
	SSEEnableHeartbeat       bool          `mapstructure:"SSE_ENABLE_HEARTBEAT"`
	SSEHeartbeatInterval     time.Duration `mapstructure:"SSE_HEARTBEAT_INTERVAL"`
	SSECleanupInterval       time.Duration `mapstructure:"SSE_CLEANUP_INTERVAL"`

	// Scheduler Configuration
	SchedulerEnabled           bool   `mapstructure:"SCHEDULER_ENABLED"`
	BillingCycleSchedule       string `mapstructure:"BILLING_CYCLE_SCHEDULE"`
	PaymentReminderSchedule    string `mapstructure:"PAYMENT_REMINDER_SCHEDULE"`
	PolicyLapseSchedule        string `mapstructure:"POLICY_LAPSE_SCHEDULE"`
	PreAuthExpirySchedule      string `mapstructure:"PREAUTH_EXPIRY_SCHEDULE"`
	RemittanceCycleSchedule    string `mapstructure:"REMITTANCE_CYCLE_SCHEDULE"`
	PaymentRetrySchedule       string `mapstructure:"PAYMENT_RETRY_SCHEDULE"`
	ReconciliationSchedule     string `mapstructure:"RECONCILIATION_SCHEDULE"`
	NotificationRetrySchedule  string `mapstructure:"NOTIFICATION_RETRY_SCHEDULE"`
	ReportDistributionSchedule string `mapstructure:"REPORT_DISTRIBUTION_SCHEDULE"`
	ReportCleanupSchedule      string `mapstructure:"REPORT_CLEANUP_SCHEDULE"`
}

var AppConfig Config

func GetEnvironment() string {
	return AppConfig.Environment
}

func LoadConfig(path string) (config Config, localConfigLoaded bool, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("app")
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		log.Printf("Warning: Could not read config file: %v", err)
		log.Printf("Continuing with environment variables only")
		localConfigLoaded = false
	} else {
		log.Printf("Using local app.env configuration")
		err = viper.Unmarshal(&config)
		if err != nil {
			return
		}
		localConfigLoaded = true
	}

	return
}
