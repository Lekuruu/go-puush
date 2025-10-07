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

// Cleanup is a routine to clean up old entries to avoid unbounded memory growth
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
