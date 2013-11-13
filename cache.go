/*
 * Simple caching library with expiration capabilities
 *     Copyright (c) 2012, Radu Ioan Fericean
 *                   2013, Christian Muehlhaeuser <muesli@gmail.com>
 *
 *   For license see LICENSE.txt
 */

package cache2go

import (
	"sync"
)

var (
	cache = make(map[string]*CacheTable)
	mutex sync.RWMutex
)

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
