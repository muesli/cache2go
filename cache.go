/*
 * Simple caching library with expiration capabilities
 *     Copyright (c) 2012, Radu Ioan Fericean
 *                   2013-2017, Christian Muehlhaeuser <muesli@gmail.com>
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

// Cache returns the existing cache table with given name or creates a new one
// if the table does not exist yet.
func Cache(table string) *CacheTable {
	mutex.RLock()
	t, ok := cache[table]
	mutex.RUnlock()

	if !ok {
		mutex.Lock()
		t, ok = cache[table]
		// Double check whether the table exists or not.
		if !ok {
			t = &CacheTable{
				name:  table,
				items: make(map[interface{}]*CacheItem),
			}
			cache[table] = t
		}
		mutex.Unlock()
	}

	return t
}

// RemoveCache flushes and removes the table if it exists.
func RemoveCache(table string) {
	mutex.RLock()
	defer mutex.RUnlock()

	// since we remove it, its cache is unnecessary ao we flush it.
	cacheTable := cache[table]
	if cacheTable != nil {
		cacheTable.Flush()
	}

	delete(cache, table)
}

// AllTables returns name list of all tables.
func AllTables() []string {
	mutex.RLock()
	defer mutex.RUnlock()

	var tables []string

	for k := range cache {
		tables = append(tables, k)
	}

	return tables
}
