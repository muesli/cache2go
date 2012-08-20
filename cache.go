//Simple caching library with expiration capabilities
package cache

import (
	"errors"
	"sync"
	"time"
)

type expiringCacheEntry interface {
	XCache(key string, expire time.Duration, value expiringCacheEntry)
	KeepAlive()
}

// Structure that must be embeded in the objectst that must be cached with expiration.
// If the expiration is not needed this can be ignored
type XEntry struct {
	sync.Mutex
	key            string
	keepAlive      bool
	expireDuration time.Duration
}

var (
	xcache = make(map[string]expiringCacheEntry)
	cache  = make(map[string]interface{})
)

// The main function to cache with expiration
func (xe *XEntry) XCache(key string, expire time.Duration, value expiringCacheEntry) {
	xe.keepAlive = true
	xe.key = key
	xe.expireDuration = expire
	xcache[key] = value
	go xe.expire()
}

// The internal mechanism for expiartion
func (xe *XEntry) expire() {
	for xe.keepAlive {
		xe.Lock()
		xe.keepAlive = false
		xe.Unlock()
		t := time.NewTimer(xe.expireDuration)
		<-t.C
		xe.Lock()
		if !xe.keepAlive {
			delete(xcache, xe.key)
		}
		xe.Unlock()
	}
}

// Mark entry to be kept another expirationDuration period
func (xe *XEntry) KeepAlive() {
	xe.Lock()
	defer xe.Unlock()
	xe.keepAlive = true
}

// Delete all keys from expiraton cache
func (xe *XEntry) Flush() {
	xcache = make(map[string]expiringCacheEntry)
}

// Get an entry from the expiration cache and mark it for keeping alive
func GetXCached(key string) (ece expiringCacheEntry, err error) {
	if r, ok := xcache[key]; ok {
		r.KeepAlive()
		return r, nil
	}
	return nil, errors.New("not found")
}

// The function to be used to cache a key/value pair when expiration is not needed
func Cache(key string, value interface{}) {
	cache[key] = value
}

// The function to extract a value for a key that never expire
func GetCached(key string) (v interface{}, err error) {
	if r, ok := cache[key]; ok {
		return r, nil
	}
	return nil, errors.New("not found")
}

// Delete all keys from cache
func Flush() {
	cache = make(map[string]interface{})
}
