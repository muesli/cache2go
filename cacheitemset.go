package cache2go

import (
	"sync"
	"time"
)

type Set struct {
	sync.RWMutex
	len     int
	members map[interface{}]int
}

func NewSet() *Set {
	set := &Set{
		len:     0,
		members: make(map[interface{}]int),
	}
	return set
}

func (set *Set) Len() int {
	set.RLock()
	defer set.RUnlock()
	return set.len
}

func (set *Set) SetAdd(member interface{}) error {
	set.Lock()
	defer set.Unlock()
	set.members[member] = 0
	set.len++
	return nil
}

func (set *Set) SetHas(member interface{}) bool {
	set.RLock()
	defer set.RUnlock()

	_, ok := set.members[member]
	return ok
}

func (set *Set) SetRemove(member interface{}) error {
	set.Lock()
	defer set.Unlock()

	if _, ok := set.members[member]; ok {
		delete(set.members, member)
		set.len--
	}
	return nil
}

// add member to table
func (table *CacheTable) SAdd(key interface{}, lifeSpan time.Duration, member interface{}) *CacheItem {
	table.Lock()
	item, ok := table.items[key]
	if !ok {
		set := NewSet()
		set.SetAdd(member)
		hashItem := NewCacheItem(key, lifeSpan, set)
		table.addInternal(hashItem)
		return hashItem
	}
	table.Unlock()

	item.Lock()
	if set, ok := item.Data().(*Set); ok {
		set.SetAdd(member)
	}
	item.Unlock()

	return item
}

// is member
func (table *CacheTable) SIsMember(key interface{}, member interface{}) bool {
	table.RLock()
	defer table.RUnlock()

	setItem, ok := table.items[key]
	if !ok {
		return false
	}

	set, ok := setItem.Data().(*Set)
	if !ok {
		return false
	}

	return set.SetHas(member)
}

// delete member
func (table *CacheTable) SDelete(key interface{}, member interface{}) error {
	table.Lock()
	defer table.Unlock()

	setItem, ok := table.items[key]
	if !ok {
		return ErrKeyNotFound
	}

	set, ok := setItem.Data().(*Set)
	if !ok {
		return ErrKeyTypeNotSet
	}

	return set.SetRemove(member)
}
