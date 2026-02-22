package configs

import (
	"context"
	"fmt"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
)

const DefaultParameterPrefix = "/hias-core/"

type SSMConfig struct {
	DBSource                         string
	Environment                      string
	HTTPServerAddress                string
	GRPCServerAddress                string
	AWSRegion                        string
	RedisURL                         string
	CognitoClientID                  string
	CognitoClientSecret              string
	CognitoRedirectURI               string
	CognitoDomain                    string
	CognitoUserPoolID                string
	TokenSymmetricKey                string
	AccessTokenDuration              string
	RefreshTokenDuration             string
	CookieDomain                     string
	AllowedOrigins                   string
	DashboardURL                     string
	AWSS3Bucket                      string
	AWSS3Region                      string
	AWSS3CDNDomain                   string
	AWSSQSDocumentProcessingQueueURL string
	AWSSQSExtractionResultsQueueURL  string
	AWSSQSClaimProcessingQueueURL    string
	AWSSQSPaymentEventsQueueURL      string
	AWSSQSNotificationEventsQueueURL string
	ConsumerType                     string
	QueueManagerType                 string
	WatermillWaitTimeSeconds         string
	WatermillMaxMessages             string
	WatermillVisibilityTimeout       string
	MpesaConsumerKey                 string
	MpesaConsumerSecret              string
	MpesaPasskey                     string
	MpesaShortcode                   string
	SchedulerEnabled                 string
}

type SSMManager struct {
	client    *ssm.Client
	cache     map[string]string
	cacheMu   sync.RWMutex
	cacheTTL  time.Duration
	cacheTime time.Time
}

func NewSSMManager(cacheTTL time.Duration) (*SSMManager, error) {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		return nil, fmt.Errorf("failed to load AWS config: %w", err)
	}

	client := ssm.NewFromConfig(cfg)

	return &SSMManager{
		client:   client,
		cache:    make(map[string]string),
		cacheTTL: cacheTTL,
	}, nil
}

func (m *SSMManager) LoadConfig() (*SSMConfig, error) {
	params, err := m.LoadConfigByPrefix(DefaultParameterPrefix)
	if err != nil {
		return nil, err
	}

	ssmConfig := &SSMConfig{}
	for key, value := range params {
		cleanKey := strings.TrimPrefix(key, DefaultParameterPrefix)
		switch cleanKey {
		case "DB_SOURCE":
			ssmConfig.DBSource = value
		case "ENVIRONMENT":
			ssmConfig.Environment = value
		case "HTTP_SERVER_ADDRESS":
			ssmConfig.HTTPServerAddress = value
		case "GRPC_SERVER_ADDRESS":
			ssmConfig.GRPCServerAddress = value
		case "AWS_REGION":
			ssmConfig.AWSRegion = value
		case "REDIS_URL":
			ssmConfig.RedisURL = value
		case "COGNITO_CLIENT_ID":
			ssmConfig.CognitoClientID = value
		case "COGNITO_CLIENT_SECRET":
			ssmConfig.CognitoClientSecret = value
		case "COGNITO_REDIRECT_URI":
			ssmConfig.CognitoRedirectURI = value
		case "COGNITO_DOMAIN":
			ssmConfig.CognitoDomain = value
		case "COGNITO_USER_POOL_ID":
			ssmConfig.CognitoUserPoolID = value
		case "TOKEN_SYMMETRIC_KEY":
			ssmConfig.TokenSymmetricKey = value
		case "ACCESS_TOKEN_DURATION":
			ssmConfig.AccessTokenDuration = value
		case "REFRESH_TOKEN_DURATION":
			ssmConfig.RefreshTokenDuration = value
		case "COOKIE_DOMAIN":
			ssmConfig.CookieDomain = value
		case "ALLOWED_ORIGINS":
			ssmConfig.AllowedOrigins = value
		case "DASHBOARD_URL":
			ssmConfig.DashboardURL = value
		case "AWS_S3_BUCKET":
			ssmConfig.AWSS3Bucket = value
		case "AWS_S3_REGION":
			ssmConfig.AWSS3Region = value
		case "AWS_S3_CDN_DOMAIN":
			ssmConfig.AWSS3CDNDomain = value
		case "AWS_SQS_DOCUMENT_PROCESSING_QUEUE_URL":
			ssmConfig.AWSSQSDocumentProcessingQueueURL = value
		case "AWS_SQS_EXTRACTION_RESULTS_QUEUE_URL":
			ssmConfig.AWSSQSExtractionResultsQueueURL = value
		case "AWS_SQS_CLAIM_PROCESSING_QUEUE_URL":
			ssmConfig.AWSSQSClaimProcessingQueueURL = value
		case "AWS_SQS_PAYMENT_EVENTS_QUEUE_URL":
			ssmConfig.AWSSQSPaymentEventsQueueURL = value
		case "AWS_SQS_NOTIFICATION_EVENTS_QUEUE_URL":
			ssmConfig.AWSSQSNotificationEventsQueueURL = value
		case "CONSUMER_TYPE":
			ssmConfig.ConsumerType = value
		case "QUEUE_MANAGER_TYPE":
			ssmConfig.QueueManagerType = value
		case "WATERMILL_WAIT_TIME_SECONDS":
			ssmConfig.WatermillWaitTimeSeconds = value
		case "WATERMILL_MAX_MESSAGES":
			ssmConfig.WatermillMaxMessages = value
		case "WATERMILL_VISIBILITY_TIMEOUT":
			ssmConfig.WatermillVisibilityTimeout = value
		case "MPESA_CONSUMER_KEY":
			ssmConfig.MpesaConsumerKey = value
		case "MPESA_CONSUMER_SECRET":
			ssmConfig.MpesaConsumerSecret = value
		case "MPESA_PASSKEY":
			ssmConfig.MpesaPasskey = value
		case "MPESA_SHORTCODE":
			ssmConfig.MpesaShortcode = value
		case "SCHEDULER_ENABLED":
			ssmConfig.SchedulerEnabled = value
		}
	}

	return ssmConfig, nil
}

func (m *SSMManager) LoadConfigByPrefix(prefix string) (map[string]string, error) {
	m.cacheMu.RLock()
	if time.Since(m.cacheTime) < m.cacheTTL && len(m.cache) > 0 {
		result := make(map[string]string, len(m.cache))
		for k, v := range m.cache {
			result[k] = v
		}
		m.cacheMu.RUnlock()
		return result, nil
	}
	m.cacheMu.RUnlock()

	m.cacheMu.Lock()
	defer m.cacheMu.Unlock()

	params := make(map[string]string)
	var nextToken *string

	for {
		input := &ssm.GetParametersByPathInput{
			Path:           &prefix,
			Recursive:      boolPtr(true),
			WithDecryption: boolPtr(true),
			NextToken:      nextToken,
		}

		output, err := m.client.GetParametersByPath(context.TODO(), input)
		if err != nil {
			log.Printf("Warning: Failed to load SSM parameters with prefix %s: %v", prefix, err)
			return params, nil
		}

		for _, param := range output.Parameters {
			if param.Name != nil && param.Value != nil {
				params[*param.Name] = *param.Value
			}
		}

		if output.NextToken == nil {
			break
		}
		nextToken = output.NextToken
	}

	m.cache = params
	m.cacheTime = time.Now()

	return params, nil
}

func boolPtr(b bool) *bool {
	return &b
}
