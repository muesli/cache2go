/**
cacheItem struct
	data *container/list.List

So data is *list.List,cacheItem is listItem
*/
package cache2go

import (
	"container/list"
	"time"
)

func (table *CacheTable) LPush(key interface{}, lifeSpan time.Duration, value interface{}) error {
	return table.push(true, key, lifeSpan, value)
}

func (table *CacheTable) LPop(key interface{}) (interface{}, error) {
	return table.pop(true, key)
}

func (table *CacheTable) RPush(key interface{}, lifeSpan time.Duration, value interface{}) error {
	return table.push(false, key, lifeSpan, value)
}

func (table *CacheTable) RPop(key interface{}) (interface{}, error) {
	return table.pop(false, key)
}

func (table *CacheTable) ListLength(key interface{}) (int, error) {
	table.RLock()
	defer table.RUnlock()

	listItem, ok := table.items[key]
	if !ok {
		return 0, ErrKeyNotFound
	}

	listItem.RLock()
	defer listItem.RUnlock()

	list, ok := listItem.Data().(*list.List)
	if !ok {
		return 0, ErrKeyTypeNotList
	}
	return list.Len(), nil
}

// push value into list
// fromLeft true:lpush Or lPop false:rpush,rPop
func (table *CacheTable) push(fromLeft bool, key interface{}, lifeSpan time.Duration, value interface{}) error {
	table.Lock()
	listItem, ok := table.items[key]
	if !ok {
		newList := list.New()
		if fromLeft {
			newList.PushFront(value)
		} else {
			newList.PushBack(value)
		}
		newListItem := NewCacheItem(key, lifeSpan, newList)
		table.addInternal(newListItem)
		return nil
	}
	table.Unlock()

	listItem.Lock()
	defer listItem.Unlock()
	listObj, ok := listItem.Data().(*list.List)
	if !ok {
		return ErrKeyTypeNotList
	}
	if fromLeft {
		listObj.PushFront(value)
	} else {
		listObj.PushBack(value)
	}
	return nil
}

func (table *CacheTable) pop(fromLeft bool, key interface{}) (interface{}, error) {
	table.RLock()
	listItem, ok := table.items[key]
	table.RUnlock()
	if !ok {
		return nil, ErrKeyNotFound
	}

	listItem.RLock()
	listObj, ok := listItem.Data().(*list.List)
	listItem.RUnlock()
	if !ok {
		return nil, ErrKeyTypeNotList
	}

	var popElement *list.Element
	if fromLeft {
		popElement = listObj.Front()
	} else {
		popElement = listObj.Back()
	}

	listObj.Remove(popElement)
	return popElement.Value, nil
}
