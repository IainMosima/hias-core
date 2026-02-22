package queue

type QueueConfig struct {
	AWSRegion                    string
	AWSEndpointURL               string
	DocumentProcessingQueueURL   string
	ExtractionResultsQueueURL    string
	ClaimProcessingQueueURL      string
	PaymentEventsQueueURL        string
	NotificationEventsQueueURL   string
	WaitTimeSeconds              int32
	MaxMessages                  int
	VisibilityTimeout            int
}

const (
	TopicDocumentProcessing = "document-processing"
	TopicExtractionResults  = "extraction-results"
	TopicClaimProcessing    = "claim-processing"
	TopicPaymentEvents      = "payment-events"
	TopicNotificationEvents = "notification-events"
)

func (c *QueueConfig) GetQueueURL(topic string) string {
	switch topic {
	case TopicDocumentProcessing:
		return c.DocumentProcessingQueueURL
	case TopicExtractionResults:
		return c.ExtractionResultsQueueURL
	case TopicClaimProcessing:
		return c.ClaimProcessingQueueURL
	case TopicPaymentEvents:
		return c.PaymentEventsQueueURL
	case TopicNotificationEvents:
		return c.NotificationEventsQueueURL
	default:
		return ""
	}
}
