package lib

import (
	"sync"
	"time"
)

type CooldownManager struct {
	mu sync.Mutex
	cooldowns map[string]struct{}
	cooldownDuration time.Duration
}

func NewCooldownManager(duration time.Duration) *CooldownManager {
	return &CooldownManager{
		cooldowns: make(map[string]struct{}),
		cooldownDuration: duration,
	}
}

func (c *CooldownManager) Add(userID string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.cooldowns[userID] = struct{}{}

	time.AfterFunc(c.cooldownDuration, func() {
		c.Remove(userID)
	})
}

func (c *CooldownManager) IsOnCooldown(userID string) bool {
	c.mu.Lock()
	defer c.mu.Unlock()

	_, exists := c.cooldowns[userID]
	return exists
}

func (c *CooldownManager) Remove(userID string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	delete(c.cooldowns, userID)
}
