package ses

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/sesv2"
	"golang.org/x/time/rate"
)

// QuotaSnapshot holds a point-in-time view of the SES sending quota.
type QuotaSnapshot struct {
	MaxSendRate     float64
	Max24HourSend   float64
	SentLast24Hours float64
}

// QuotaManager polls SES GetAccount every 60 s and enforces the send-rate limit.
type QuotaManager struct {
	client  *sesv2.Client
	limiter *rate.Limiter
	max24h  float64
	sent24h float64
	mu      sync.RWMutex
	stopCh  chan struct{}
	log     *log.Logger
}

func newQuotaManager(client *sesv2.Client, lo *log.Logger) *QuotaManager {
	return &QuotaManager{
		client:  client,
		limiter: rate.NewLimiter(rate.Inf, 1),
		stopCh:  make(chan struct{}),
		log:     lo,
	}
}

// Start launches the background polling goroutine.
func (q *QuotaManager) Start() {
	go func() {
		q.refresh()
		ticker := time.NewTicker(60 * time.Second)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				q.refresh()
			case <-q.stopCh:
				return
			}
		}
	}()
}
