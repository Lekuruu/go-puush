package app

import (
	"sync"
	"time"
)

// CooldownManager manages cooldowns in memory.
type CooldownManager struct {
	data  map[string]time.Time
	mutex sync.RWMutex
	ttl   time.Duration
}

// NewCooldownManager creates a CooldownManager with a cooldown duration
func NewCooldownManager(ttl time.Duration) *CooldownManager {
	return &CooldownManager{
		data: make(map[string]time.Time),
		ttl:  ttl,
	}
}

// NewCooldownManagerWithCleanup creates a CooldownManager and starts automatic cleanup
func NewCooldownManagerWithCleanup(ttl time.Duration, cleanupInterval time.Duration) *CooldownManager {
	cm := NewCooldownManager(ttl)
	cm.StartCleanup(cleanupInterval)
	return cm
}

// Allow checks if an action is allowed, based on cooldown
func (cm *CooldownManager) Allow(key string) bool {
	cm.mutex.Lock()
	defer cm.mutex.Unlock()

	now := time.Now()
	last, exists := cm.data[key]

	if !exists || now.Sub(last) >= cm.ttl {
		cm.data[key] = now
		return true
	}
	return false
}

// Cleanup removes expired entries to avoid unbounded memory growth
func (cm *CooldownManager) Cleanup() {
	cm.mutex.Lock()
	defer cm.mutex.Unlock()

	now := time.Now()
	for key, last := range cm.data {
		if now.Sub(last) >= cm.ttl*2 {
			delete(cm.data, key)
		}
	}
}

// StartCleanup starts a background goroutine that periodically cleans up expired entries
func (cm *CooldownManager) StartCleanup(interval time.Duration) {
	go func() {
		ticker := time.NewTicker(interval)
		defer ticker.Stop()
		for range ticker.C {
			cm.Cleanup()
		}
	}()
}
