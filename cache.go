/*
Rating system designed to be used in VoIP Carriers World
Copyright (C) 2012  Radu Ioan Fericean

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program.  If not, see <http://www.gnu.org/licenses/>
*/

package cache

import (
	"errors"
	"sync"
	"time"
)

type ExpiringCacheEntry interface {
	XCache(key string, expire time.Duration, value ExpiringCacheEntry)
	Expire()
	KeepAlive()
}

type XEntry struct {
	sync.Mutex
	key       string
	keepAlive bool
	expire    time.Duration
}

var (
	xcache = make(map[string]ExpiringCacheEntry)
	cache  = make(map[string]interface{})
)

func (xe *XEntry) XCache(key string, expire time.Duration, value ExpiringCacheEntry) {
	xe.keepAlive = true
	xe.key = key
	xe.expire = expire
	xcache[key] = value
	go xe.Expire()
}

func (xe *XEntry) Expire() {
	for xe.keepAlive {
		xe.Lock()
		xe.keepAlive = false
		xe.Unlock()
		t := time.NewTimer(xe.expire)
		<-t.C
		xe.Lock()
		if !xe.keepAlive {
			delete(xcache, xe.key)
		}
		xe.Unlock()
	}
}

func (xe *XEntry) KeepAlive() {
	xe.Lock()
	defer xe.Unlock()
	xe.keepAlive = true
}

func GetXCached(key string) (ece ExpiringCacheEntry, err error) {
	if r, ok := xcache[key]; ok {
		r.KeepAlive()
		return r, nil
	}
	return nil, errors.New("not found")
}

func Cache(key string, value interface{}) {
	cache[key] = value
}

func GetCached(key string) (v interface{}, err error) {
	if r, ok := cache[key]; ok {
		return r, nil
	}
	return nil, errors.New("not found")
}
