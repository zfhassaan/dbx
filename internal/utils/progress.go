package utils

import (
	"fmt"
	"time"
)

// ProgressIndicator provides simple progress feedback for long-running operations
type ProgressIndicator struct {
	message   string
	startTime time.Time
	lastUpdate time.Time
	interval  time.Duration
}

// NewProgressIndicator creates a new progress indicator
func NewProgressIndicator(message string) *ProgressIndicator {
	return &ProgressIndicator{
		message:    message,
		startTime:  time.Now(),
		lastUpdate: time.Now(),
		interval:   5 * time.Second, // Update every 5 seconds
	}
}

// Update prints progress message if enough time has passed
func (p *ProgressIndicator) Update() {
	now := time.Now()
	if now.Sub(p.lastUpdate) >= p.interval {
		elapsed := now.Sub(p.startTime).Round(time.Second)
		fmt.Printf("⏳ %s... (elapsed: %s)\n", p.message, elapsed)
		p.lastUpdate = now
	}
}

// Finish prints final progress message
func (p *ProgressIndicator) Finish(success bool) {
	elapsed := time.Since(p.startTime).Round(time.Second)
	if success {
		fmt.Printf("✅ %s completed in %s\n", p.message, elapsed)
	} else {
		fmt.Printf("❌ %s failed after %s\n", p.message, elapsed)
	}
}

