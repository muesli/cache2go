/**
cacheItem struct
	data map[interface{}]interface{}

So data is hash, cacheItem is hashItem
*/
package cache2go

import "time"

func (table *CacheTable) HAdd(key interface{}, lifeSpan time.Duration, hkey interface{}, hvalue interface{}) *CacheItem {
	table.Lock()
	item, ok := table.items[key]
	if !ok {
		hash := make(map[interface{}]interface{})
		hash[hkey] = hvalue
		hashItem := NewCacheItem(key, lifeSpan, hash)
		table.addInternal(hashItem)
		return hashItem
	}
	table.Unlock()

	item.Lock()
	if hash, ok := item.Data().(map[interface{}]interface{}); ok {
		hash[hkey] = hvalue
	}
	item.Unlock()

	return item

}

func (table *CacheTable) HValue(key interface{}, hkey interface{}) (interface{}, error) {
	table.RLock()
	hashItem, ok := table.items[key]
	table.RUnlock()
	if !ok {
		return nil, ErrKeyNotFound
	}

	hvalue, hok := hashItem.Data().(map[interface{}]interface{})[hkey]
	if !hok {
		return nil, ErrKeyNotFound
	}

	hashItem.KeepAlive()
	return hvalue, nil
}

func (table *CacheTable) HDelete(key interface{}, hkey interface{}) error {
	table.RLock()
	defer table.RUnlock()

	hashItem, ok := table.items[key]
	if !ok {
		return ErrKeyNotFound
	}

	hashItem.Lock()
	defer hashItem.Unlock()

	hash, ok := hashItem.Data().(map[interface{}]interface{})
	if !ok {
		return ErrKeyTypeNotHash
	}

	delete(hash, hkey)
	return nil
}
