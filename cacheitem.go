/*
 * Simple caching library with expiration capabilities
 *     Copyright (c) 2013, Christian Muehlhaeuser <muesli@gmail.com>
 *
 *   For license see LICENSE.txt
 */

package cache2go

import (
	"sync"
	"time"
)

// Structure of an item in the cache.
// Parameter data contains the user-set value in the cache.
type CacheItem struct {
	sync.RWMutex

	// The item's key.
	key interface{}
	// The item's data.
	data interface{}
	// How long will the item live in the cache when not being accessed/kept alive.
	lifeSpan time.Duration

	// Creation timestamp.
	createdOn time.Time
	// Last access timestamp.
	accessedOn time.Time
	// How often the item was accessed.
	accessCount int64

	// Callback method triggered right before removing the item from the cache
	aboutToExpire func(interface{})
}

// Returns a newly created CacheItem.
// Parameter key is the item's cache-key.
// Parameter lifeSpan determines after which time period without an access the item
// will get removed from the cache.
// Parameter data is the item's value.
func CreateCacheItem(key interface{}, lifeSpan time.Duration, data interface{}) CacheItem {
	t := time.Now()
	return CacheItem{
		key:           key,
		lifeSpan:      lifeSpan,
		createdOn:     t,
		accessedOn:    t,
		accessCount:   0,
		aboutToExpire: nil,
		data:          data,
	}
}

// Mark item to be kept for another expireDuration period.
func (item *CacheItem) KeepAlive() {
	item.Lock()
	defer item.Unlock()
	item.accessedOn = time.Now()
	item.accessCount++
}

// Returns this item's expiration duration.
func (item *CacheItem) LifeSpan() time.Duration {
	// immutable
	return item.lifeSpan
}

// Returns when this item was last accessed.
func (item *CacheItem) AccessedOn() time.Time {
	item.RLock()
	defer item.RUnlock()
	return item.accessedOn
}

// Returns when this item was added to the cache.
func (item *CacheItem) CreatedOn() time.Time {
	// immutable
	return item.createdOn
}

// Returns how often this item has been accessed.
func (item *CacheItem) AccessCount() int64 {
	item.RLock()
	defer item.RUnlock()
	return item.accessCount
}

// Returns the key of this cached item.
func (item *CacheItem) Key() interface{} {
	// immutable
	return item.key
}

// Returns the value of this cached item.
func (item *CacheItem) Data() interface{} {
	// immutable
	return item.data
}

// Configures a callback, which will be called right before the item
// is about to be removed from the cache.
func (item *CacheItem) SetAboutToExpireCallback(f func(interface{})) {
	item.Lock()
	defer item.Unlock()
	item.aboutToExpire = f
}
