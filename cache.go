//Simple caching library with expiration capabilities
package cache2go

import (
	"errors"
	"sync"
	"time"
)

type expiringCacheEntry interface {
	XCache(key string, expire time.Duration, value expiringCacheEntry)
	timer() *time.Timer
	KeepAlive()
}

// Structure that must be embeded in the objectst that must be cached with expiration.
// If the expiration is not needed this can be ignored
type XEntry struct {
	sync.Mutex
	key            string
	keepAlive      bool
	expireDuration time.Duration
	t              *time.Timer
}

var (
	xcache = make(map[string]expiringCacheEntry)
	xMux   sync.RWMutex
	cache  = make(map[string]interface{})
	mux    sync.RWMutex
)

// The main function to cache with expiration
func (xe *XEntry) XCache(key string, expire time.Duration, value expiringCacheEntry) {
	xe.keepAlive = true
	xe.key = key
	xe.expireDuration = expire
	xMux.Lock()
	xcache[key] = value
	xMux.Unlock()
	go xe.expire()
}

// The internal mechanism for expiartion
func (xe *XEntry) expire() {
	for xe.keepAlive {
		xe.Lock()
		xe.keepAlive = false
		xe.Unlock()
		xe.t = time.NewTimer(xe.expireDuration)
		<-xe.t.C
		if !xe.keepAlive {
			xMux.Lock()
			delete(xcache, xe.key)
			xMux.Unlock()
		}
	}
}

// Getter for the timer
func (xe *XEntry) timer() *time.Timer {
	return xe.t
}

// Mark entry to be kept another expirationDuration period
func (xe *XEntry) KeepAlive() {
	xe.Lock()
	defer xe.Unlock()
	xe.keepAlive = true
}

// Get an entry from the expiration cache and mark it for keeping alive
func GetXCached(key string) (ece expiringCacheEntry, err error) {
	xMux.RLock()
	defer xMux.RUnlock()
	if r, ok := xcache[key]; ok {
		r.KeepAlive()
		return r, nil
	}
	return nil, errors.New("not found")
}

// The function to be used to cache a key/value pair when expiration is not needed
func Cache(key string, value interface{}) {
	mux.Lock()
	defer mux.Unlock()
	cache[key] = value
}

// The function to extract a value for a key that never expire
func GetCached(key string) (v interface{}, err error) {
	mux.RLock()
	defer mux.RUnlock()
	if r, ok := cache[key]; ok {
		return r, nil
	}
	return nil, errors.New("not found")
}

// Delete all keys from expiraton cache
func XFlush() {
	xMux.Lock()
	defer xMux.Unlock()
	for _, v := range xcache {
		if v.timer() != nil {
			v.timer().Stop()
		}
	}
	xcache = make(map[string]expiringCacheEntry)
}

// Delete all keys from cache
func Flush() {
	mux.Lock()
	defer mux.Unlock()
	cache = make(map[string]interface{})
}
