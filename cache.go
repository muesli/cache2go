// Simple caching library with expiration capabilities
package cache2go

import (
	"errors"
	"sync"
	"time"
)

// Structure of an entry in the cache
// data contains the user-set value in the cache
type CacheEntry struct {
	sync.Mutex
	key            string
	keepAlive      bool
	expireDuration time.Duration
	expiringSince  time.Time
	data           interface{}

	// Callback method triggered right before removing the item from the cache
	aboutToExpire func(string)
}

// Structure of a table with items in the cache
type CacheTable struct {
	sync.RWMutex
	name        string
	items       map[string]*CacheEntry
	expTimer    *time.Timer
	expDuration time.Duration
}

var (
	cache = make(map[string]*CacheTable)
	mutex sync.RWMutex
)

// Mark entry to be kept for another expireDuration period
func (entry *CacheEntry) KeepAlive() {
	entry.Lock()
	defer entry.Unlock()
	entry.expiringSince = time.Now()
}

// Returns this entry's expiration duration
func (entry *CacheEntry) ExpireDuration() time.Duration {
	entry.Lock()
	defer entry.Unlock()
	return entry.expireDuration
}

// Returns since when this entry is expiring
func (entry *CacheEntry) ExpiringSince() time.Time {
	entry.Lock()
	defer entry.Unlock()
	return entry.expiringSince
}

// Returns the value of this cached item
func (entry *CacheEntry) Data() interface{} {
	entry.Lock()
	defer entry.Unlock()
	return entry.data
}

// Returns the existing cache table with given name or creates a new one
// if the table does not exist yet
func Cache(table string) *CacheTable {
	mutex.RLock()
	t, ok := cache[table]
	mutex.RUnlock()

	if !ok {
		t = &CacheTable{
			name:  table,
			items: make(map[string]*CacheEntry),
		}
		mutex.Lock()
		cache[table] = t
		mutex.Unlock()
	}

	return t
}

// Returns how many items are currently stored in the cache
func (table *CacheTable) CacheCount() int {
	table.RLock()
	defer table.RUnlock()

	return len(table.items)
}

// Expiration check loop, triggered by a self-adjusting timer
func (table *CacheTable) expirationCheck() {
	table.Lock()
	if table.expTimer != nil {
		table.expTimer.Stop()
	}

	// Take a copy of cache so we can iterate over it without blocking the mutex
	cc := table.items
	table.Unlock()

	// To be more accurate with timers, we would need to update 'now' on every
	// loop iteration. Not sure it's really efficient though.
	now := time.Now()
	smallestDuration := 0 * time.Second
	for key, c := range cc {
		if c.expireDuration == 0 {
			continue
		}
		if now.Sub(c.expiringSince) >= c.expireDuration {
			table.Lock()
			if c.aboutToExpire != nil {
				c.aboutToExpire(key)
			}
			delete(table.items, key)
			table.Unlock()
		} else {
			if smallestDuration == 0 || c.ExpireDuration() < smallestDuration {
				smallestDuration = c.ExpireDuration() - now.Sub(c.ExpiringSince())
			}
		}
	}

	table.Lock()
	table.expDuration = smallestDuration
	if smallestDuration > 0 {
		table.expTimer = time.AfterFunc(smallestDuration, func() {
			go table.expirationCheck()
		})
	}
	table.Unlock()
}

// Adds a key/value pair to the cache
// The last parameter abouToExpireFunc can be nil. Otherwise abouToExpireFunc
// will be called (with this item's key as its only parameter), right before
// removing this item from the cache
func (table *CacheTable) Cache(key string, expire time.Duration, data interface{}, aboutToExpireFunc func(string)) {
	entry := CacheEntry{}
	entry.keepAlive = true
	entry.key = key
	entry.expireDuration = expire
	entry.expiringSince = time.Now()
	entry.aboutToExpire = aboutToExpireFunc
	entry.data = data

	table.Lock()
	table.items[key] = &entry
	expDur := table.expDuration
	table.Unlock()

	// If we haven't set up any expiration check timer or found a more imminent item
	if expire > 0 && (expDur == 0 || expire < expDur) {
		table.expirationCheck()
	}
}

// Get an entry from the cache and mark it to be kept alive
func (table *CacheTable) Value(key string) (*CacheEntry, error) {
	table.RLock()
	defer table.RUnlock()
	if r, ok := table.items[key]; ok {
		r.KeepAlive()
		return r, nil
	}
	return nil, errors.New("Key not found in cache")
}

// Delete all items from cache
func (table *CacheTable) Flush() {
	table.Lock()
	defer table.Unlock()

	table.items = make(map[string]*CacheEntry)
	table.expDuration = 0
	if table.expTimer != nil {
		table.expTimer.Stop()
	}

	mutex.Lock()
	defer mutex.Unlock()
	delete(cache, table.name)
}
