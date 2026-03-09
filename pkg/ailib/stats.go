package ailib

import (
	"time"
)

// Version information.
const (
	VersionMajor = 0
	VersionMinor = 5
	VersionPatch = 0
	VersionStr   = "0.5.0"
)

// CallStats tracks call statistics.
type CallStats struct {
	StartTime      time.Time
	EndTime        time.Time
	Duration       time.Duration
	RequestTokens  int
	ResponseTokens int
	TotalTokens    int
	Model          string
	Provider       string
	Success        bool
	StatusCode     int
	Error          error
}

// ClientMetrics tracks client-level metrics.
type ClientMetrics struct {
	TotalRequests    int64
	SuccessfulCalls  int64
	FailedCalls      int64
	TotalTokens      int64
	TotalLatency     time.Duration
	AverageLatency   time.Duration
	LastCallTime     time.Time
}

// Record records a call.
func (m *ClientMetrics) Record(stats *CallStats) {
	m.TotalRequests++
	m.LastCallTime = stats.StartTime
	m.TotalLatency += stats.Duration

	if stats.Success {
		m.SuccessfulCalls++
		m.TotalTokens += int64(stats.TotalTokens)
	} else {
		m.FailedCalls++
	}

	if m.TotalRequests > 0 {
		m.AverageLatency = time.Duration(int64(m.TotalLatency) / m.TotalRequests)
	}
}
