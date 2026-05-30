package ses

import (
	"context"
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/sesv2"
	"github.com/knadh/listmonk/models"
)

// Options holds all configuration for the SES messenger.
type Options struct {
	Region          string
	AccessKeyID     string
	SecretAccessKey string
	RoleARN         string  // optional, for IAM role assumption
	ConfigSet       string  // optional SES configuration set name
	MaxMsgRetries   int
	PricePerMessage float64 // default 0.0001
}

// SESMessenger implements manager.Messenger via the SES v2 SendEmail API.
type SESMessenger struct {
	name   string
	opt    Options
	client *sesv2.Client
	quota  *QuotaManager
	log    *log.Logger
}

// New creates a new SESMessenger with an authenticated sesv2.Client.
func New(name string, opt Options, lo *log.Logger) (*SESMessenger, error) {
	if opt.Region == "" {
		return nil, fmt.Errorf("ses: region is required")
	}

	if opt.PricePerMessage == 0 {
		opt.PricePerMessage = 0.0001
	}

	cfg, err := config.LoadDefaultConfig(context.Background(),
		config.WithRegion(opt.Region),
		config.WithCredentialsProvider(
			credentials.NewStaticCredentialsProvider(opt.AccessKeyID, opt.SecretAccessKey, ""),
		),
	)
	if err != nil {
		return nil, fmt.Errorf("ses: failed to load AWS config: %w", err)
	}

	client := sesv2.NewFromConfig(cfg)
	quota := newQuotaManager(client, lo)
	quota.Start()

	return &SESMessenger{
		name:   name,
		opt:    opt,
		client: client,
		quota:  quota,
		log:    lo,
	}, nil
}

// Name returns the messenger's name.
func (s *SESMessenger) Name() string {
	return s.name
}

// Push sends a message via SES. Stub – full implementation in Task 5.
func (s *SESMessenger) Push(_ models.Message) error {
	return nil
}

// Flush is a no-op for SES.
func (s *SESMessenger) Flush() error {
	return nil
}

// Close stops the QuotaManager goroutine.
func (s *SESMessenger) Close() error {
	s.quota.Stop()
	return nil
}
