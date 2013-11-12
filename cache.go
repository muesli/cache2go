/*
 * Simple caching library with expiration capabilities
 *     Copyright (c) 2012, Radu Ioan Fericean
 *                   2013, Christian Muehlhaeuser <muesli@gmail.com>
 *
 *   For license see LICENSE.txt
 */

package cache2go

import (
	"errors"
	"log"
	"sync"
	"time"
)

/* Structure of an item in the cache.
Parameter data contains the user-set value in the cache.
*/
type CacheItem struct {
	sync.RWMutex

	// The item's key
	key interface{}
	// The item's data
	data interface{}
	// How long will the item live in the cache when not being accessed/kept alive
	lifeSpan time.Duration

	// Creation timestamp
	createdOn time.Time
	// Last access timestamp
	accessedOn time.Time
	// How often the item was accessed
	accessCount int64

	// Callback method triggered right before removing the item from the cache
	aboutToExpire func(interface{})
}

// Structure of a table with items in the cache
type CacheTable struct {
	sync.RWMutex

	// The table's name
	name string
	// All cached items
	items map[interface{}]*CacheItem

	// Timer responsible for triggering cleanup
	cleanupTimer *time.Timer
	// Last used timer duration
	cleanupInterval time.Duration

	// The logger used for this table
	logger *log.Logger

	// Callback method triggered when trying to load a non-existing key
	loadData func(interface{}) *CacheItem
	// Callback method triggered when adding a new item to the cache
	addedItem func(*CacheItem)
	// Callback method triggered before deleting an item from the cache
	aboutToDeleteItem func(*CacheItem)
}

var (
	cache = make(map[string]*CacheTable)
	mutex sync.RWMutex
)

/* Returns a newly created CacheItem.
Parameter key is the item's cache-key.
Parameter lifeSpan determines after which time period without an access the item
	will get removed from the cache.
Parameter data is the item's value.
*/
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

// Returns the existing cache table with given name or creates a new one
// if the table does not exist yet.
func Cache(table string) *CacheTable {
	mutex.RLock()
	t, ok := cache[table]
	mutex.RUnlock()

	if !ok {
		t = &CacheTable{
			name:  table,
			items: make(map[interface{}]*CacheItem),
		}

		mutex.Lock()
		cache[table] = t
		mutex.Unlock()
	}

	return t
}

// Returns how many items are currently stored in the cache.
func (table *CacheTable) Count() int {
	table.RLock()
	defer table.RUnlock()

	return len(table.items)
}

// Configures a data-loader callback, which will be called when trying
// to use access a non-existing key.
func (table *CacheTable) SetDataLoader(f func(interface{}) *CacheItem) {
	table.Lock()
	defer table.Unlock()
	table.loadData = f
}

// Configures a callback, which will be called every time a new item
// is added to the cache.
func (table *CacheTable) SetAddedItemCallback(f func(*CacheItem)) {
	table.Lock()
	defer table.Unlock()
	table.addedItem = f
}

// Configures a callback, which will be called every time an item
// is about to be removed from the cache.
func (table *CacheTable) SetAboutToDeleteItemCallback(f func(*CacheItem)) {
	table.Lock()
	defer table.Unlock()
	table.aboutToDeleteItem = f
}

// Sets the logger to be used by this cache table.
func (table *CacheTable) SetLogger(logger *log.Logger) {
	table.Lock()
	defer table.Unlock()
	table.logger = logger
}

// Expiration check loop, triggered by a self-adjusting timer.
func (table *CacheTable) expirationCheck() {
	table.Lock()
	if table.cleanupInterval > 0 {
		table.log("Expiration check triggered after", table.cleanupInterval, "for table", table.name)
	} else {
		table.log("Expiration check installed for table", table.name)
	}
	if table.cleanupTimer != nil {
		table.cleanupTimer.Stop()
	}

	// Take a copy of cache so we can iterate over it without blocking the mutex.
	cc := table.items
	table.Unlock()

	// To be more accurate with timers, we would need to update 'now' on every
	// loop iteration. Not sure it's really efficient though.
	now := time.Now()
	smallestDuration := 0 * time.Second
	for key, c := range cc {
		c.RLock()
		lifeSpan := c.lifeSpan
		accessedOn := c.accessedOn
		c.RUnlock()

		if lifeSpan == 0 {
			continue
		}
		if now.Sub(accessedOn) >= lifeSpan {
			table.Delete(key)
		} else {
			if smallestDuration == 0 || lifeSpan < smallestDuration {
				smallestDuration = lifeSpan - now.Sub(accessedOn)
			}
		}
	}

	table.Lock()
	table.cleanupInterval = smallestDuration
	if smallestDuration > 0 {
		table.cleanupTimer = time.AfterFunc(smallestDuration, func() {
			go table.expirationCheck()
		})
	}
	table.Unlock()
}

/* Adds a key/value pair to the cache.
Parameter key is the item's cache-key.
Parameter lifeSpan determines after which time period without an access the item
	will get removed from the cache.
Parameter data is the item's value.
*/
func (table *CacheTable) Cache(key interface{}, lifeSpan time.Duration, data interface{}) *CacheItem {
	item := CreateCacheItem(key, lifeSpan, data)

	table.Lock()
	table.log("Adding item with key", key, "and lifespan of", lifeSpan, "to table", table.name)
	table.items[key] = &item
	expDur := table.cleanupInterval
	table.Unlock()

	// Trigger callback after adding an item to cache
	if table.addedItem != nil {
		table.addedItem(&item)
	}

	// If we haven't set up any expiration check timer or found a more imminent item
	if lifeSpan > 0 && (expDur == 0 || lifeSpan < expDur) {
		table.expirationCheck()
	}

	return &item
}

// Delete an item from the cache.
func (table *CacheTable) Delete(key interface{}) (*CacheItem, error) {
	table.RLock()
	r, ok := table.items[key]

	if !ok {
		table.RUnlock()
		return nil, errors.New("Key not found in cache")
	}

	// Trigger callbacks before deleting an item from cache.
	if table.aboutToDeleteItem != nil {
		table.aboutToDeleteItem(r)
	}
	table.RUnlock()
	r.RLock()
	defer r.RUnlock()
	if r.aboutToExpire != nil {
		r.aboutToExpire(key)
	}

	table.Lock()
	defer table.Unlock()

	table.log("Deleting item with key", key, "created on", r.createdOn, "and hit", r.accessCount, "times from table", table.name)
	delete(table.items, key)
	return r, nil
}

// Test whether an item exists in the cache. Unlike the Value method
// Exists neither tries to fetch data via the loadData callback nor
// does it keep the item alive in the cache.
func (table *CacheTable) Exists(key interface{}) bool {
	table.RLock()
	defer table.RUnlock()
	_, ok := table.items[key]

	return ok
}

// Get an item from the cache and mark it to be kept alive.
func (table *CacheTable) Value(key interface{}) (*CacheItem, error) {
	table.RLock()
	r, ok := table.items[key]
	table.RUnlock()

	if ok {
		r.KeepAlive()
		return r, nil
	}

	if table.loadData != nil {
		item := table.loadData(key)
		table.Cache(key, item.lifeSpan, item.data)
		if item != nil {
			return item, nil
		}

		return nil, errors.New("Key not found and could not be loaded into cache")
	}

	return nil, errors.New("Key not found in cache")
}

// Delete all items from cache.
func (table *CacheTable) Flush() {
	table.Lock()
	defer table.Unlock()

	table.log("Flushing table", table.name)

	table.items = make(map[interface{}]*CacheItem)
	table.cleanupInterval = 0
	if table.cleanupTimer != nil {
		table.cleanupTimer.Stop()
	}
}

// Internal logging method for convenience.
func (table *CacheTable) log(v ...interface{}) {
	if table.logger == nil {
		return
	}

	table.logger.Println(v)
}
