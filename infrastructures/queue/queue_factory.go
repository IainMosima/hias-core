package queue

import (
	"context"
	"fmt"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill-aws/sqs"
	"github.com/aws/aws-sdk-go-v2/aws"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
)

type QueueFactory struct {
	config QueueConfig
	logger watermill.LoggerAdapter
}

func NewQueueFactory(config QueueConfig) *QueueFactory {
	return &QueueFactory{
		config: config,
		logger: watermill.NewStdLogger(false, false),
	}
}

func (f *QueueFactory) CreatePublisher() (*sqs.Publisher, error) {
	cfg, err := f.loadAWSConfig()
	if err != nil {
		return nil, err
	}

	publisher, err := sqs.NewPublisher(sqs.PublisherConfig{
		AWSConfig: cfg,
	}, f.logger)
	if err != nil {
		return nil, fmt.Errorf("failed to create SQS publisher: %w", err)
	}

	return publisher, nil
}

func (f *QueueFactory) CreateSubscriber() (*sqs.Subscriber, error) {
	cfg, err := f.loadAWSConfig()
	if err != nil {
		return nil, err
	}

	subscriber, err := sqs.NewSubscriber(sqs.SubscriberConfig{
		AWSConfig: cfg,
	}, f.logger)
	if err != nil {
		return nil, fmt.Errorf("failed to create SQS subscriber: %w", err)
	}

	return subscriber, nil
}

func (f *QueueFactory) loadAWSConfig() (aws.Config, error) {
	var opts []func(*awsconfig.LoadOptions) error
	opts = append(opts, awsconfig.WithRegion(f.config.AWSRegion))

	cfg, err := awsconfig.LoadDefaultConfig(context.TODO(), opts...)
	if err != nil {
		return aws.Config{}, fmt.Errorf("failed to load AWS config: %w", err)
	}

	return cfg, nil
}
