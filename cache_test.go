/*
 * Simple caching library with expiration capabilities
 *     Copyright (c) 2013, Christian Muehlhaeuser <muesli@gmail.com>
 *
 *   For license see LICENSE.txt
 */

package cache2go

import (
	"bytes"
	"log"
	"strconv"
	"testing"
	"time"
)

var (
	k = "testkey"
	v = "testvalue"
)

func TestCache(t *testing.T) {
	table := Cache("testCache")
	table.Add(k + "_1", 0*time.Second, v)
	table.Add(k + "_2", 1*time.Second, v)
	p, err := table.Value(k + "_1")
	if err != nil || p == nil || p.Data().(string) != v {
		t.Error("Error retrieving non expiring data from cache", err)
	}
	p, err = table.Value(k + "_2")
	if err != nil || p == nil || p.Data().(string) != v {
		t.Error("Error retrieving data from cache", err)
	}
	if p.AccessCount() != 1 {
		t.Error("Error getting correct access count")
	}
	if p.LifeSpan() != 1*time.Second {
		t.Error("Error getting correct life-span")
	}
	if p.AccessedOn().Unix() == 0 {
		t.Error("Error getting access time")
	}
	if p.CreatedOn().Unix() == 0 {
		t.Error("Error getting creation time")
	}
}

func TestCacheExpire(t *testing.T) {
	table := Cache("testExpire")
	table.Add(k, 250*time.Millisecond, v)
	p, err := table.Value(k)
	if err != nil || p == nil || p.Data().(string) != v {
		t.Error("Error retrieving data from cache", err)
	}
	time.Sleep(500 * time.Millisecond)
	p, err = table.Value(k)
	if err == nil || p != nil {
		t.Error("Error expiring data")
	}
}

func TestCacheKeepAlive(t *testing.T) {
	k2 := k + k
	v2 := v + v
	table := Cache("testKeepAlive")
	table.Add(k, 250*time.Millisecond, v)
	table.Add(k2, 750*time.Millisecond, v2)

	p, err := table.Value(k)
	if err != nil || p == nil || p.Data().(string) != v {
		t.Error("Error retrieving data from cache", err)
	}
	time.Sleep(50 * time.Millisecond)
	p.KeepAlive()

	time.Sleep(450 * time.Millisecond)
	p, err = table.Value(k)
	if err == nil || p != nil {
		t.Error("Error expiring data")
	}
	p, err = table.Value(k2)
	if err != nil || p == nil || p.Data().(string) != v2 {
		t.Error("Error retrieving data from cache", err)
	}
	time.Sleep(1 * time.Second)
	p, err = table.Value(k2)
	if err == nil || p != nil {
		t.Error("Error expiring data")
	}
}

func TestExists(t *testing.T) {
	table := Cache("testExists")
	table.Add(k, 0, v)
	if !table.Exists(k) {
		t.Error("Error verifying existing data in cache")
	}
}

func TestDelete(t *testing.T) {
	table := Cache("testDelete")
	table.Add(k, 0, v)
	p, err := table.Value(k)
	if err != nil || p == nil || p.Data().(string) != v {
		t.Error("Error retrieving data from cache", err)
	}
	table.Delete(k)
	p, err = table.Value(k)
	if err == nil || p != nil {
		t.Error("Error deleting data")
	}

	// test error handling
	_, err = table.Delete(k)
	if err == nil {
		t.Error("Expected error deleting item")
	}
}

func TestFlush(t *testing.T) {
	table := Cache("testFlush")
	table.Add(k, 10*time.Second, v)
	time.Sleep(100 * time.Millisecond)
	table.Flush()

	p, err := table.Value(k)
	if err == nil || p != nil {
		t.Error("Error expiring data")
	}
	if table.Count() != 0 {
		t.Error("Error verifying empty table")
	}
}

func TestFlushNoTimout(t *testing.T) {
	table := Cache("testFlushNoTimeout")
	table.Add(k, 10*time.Second, v)
	table.Flush()

	p, err := table.Value(k)
	if err == nil || p != nil {
		t.Error("Error expiring data")
	}
	if table.Count() != 0 {
		t.Error("Error verifying empty table")
	}
}

func TestCount(t *testing.T) {
	count := 100000
	table := Cache("testCount")
	for i := 0; i < count; i++ {
		key := k + strconv.Itoa(i)
		table.Add(key, 10*time.Second, v)
	}
	for i := 0; i < count; i++ {
		key := k + strconv.Itoa(i)
		p, err := table.Value(key)
		if err != nil || p == nil || p.Data().(string) != v {
			t.Error("Error retrieving data")
		}
	}
	if table.Count() != count {
		t.Error("Data count mismatch")
	}
}

func TestDataLoader(t *testing.T) {
	table := Cache("testDataLoader")
	table.SetDataLoader(func(key interface{}) *CacheItem {
		var item *CacheItem
		if key.(string) != "nil" {
			val := k + key.(string)
			i := CreateCacheItem(key, 500*time.Millisecond, val)
			item = &i
		}

		return item
	})

	p, err := table.Value("nil")
	if err == nil || table.Exists("nil") {
		t.Error("Error validating data loader for nil values")
	}

	for i := 0; i < 10; i++ {
		key := k + strconv.Itoa(i)
		vp := k + key
		p, err = table.Value(key)
		if err != nil || p == nil || p.Data().(string) != vp {
			t.Error("Error validating data loader")
		}
	}
}

func TestAccessCount(t *testing.T) {
	count := 100
	table := Cache("testAccessCount")
	for i := 0; i < count; i++ {
		table.Add(i, 10*time.Second, v)
	}
	for i := 0; i < count; i++ {
		for j := 0; j < i; j++ {
			table.Value(i)
		}
	}

	ma := table.MostAccessed(int64(count))
	for i, key := range ma {
		if key != count - 1 - i {
			t.Error("Most accessed items seem to be sorted incorrectly")
		}
	}

	ma = table.MostAccessed(int64(count - 1))
	if len(ma) != count - 1 {
		t.Error("MostAccessed returns incorrect amount of items")
	}
}

func TestCallbacks(t *testing.T) {
	addedKey := ""
	removedKey := ""
	expired := false

	table := Cache("testCallbacks")
	table.SetAddedItemCallback(func(item *CacheItem) {
		addedKey = item.Key().(string)
	})
	table.SetAboutToDeleteItemCallback(func(item *CacheItem) {
		removedKey = item.Key().(string)
	})

	i := table.Add(k, 500*time.Millisecond, v)
	i.SetAboutToExpireCallback(func(key interface{}) {
		expired = true
	})

	time.Sleep(250 * time.Millisecond)
	if addedKey != k {
		t.Error("AddedItem callback not working")
	}

	time.Sleep(500 * time.Millisecond)
	if removedKey != k {
		t.Error("AboutToDeleteItem callback not working:" + k + "_" + removedKey)
	}
	if !expired {
		t.Error("ExpireCallback callback not working")
	}
}

func TestLogger(t *testing.T) {
	out := new(bytes.Buffer)
	l := log.New(out, "cache2go ", log.Ldate|log.Ltime)

	table := Cache("testLogger")
	table.SetLogger(l)
	table.Add(k, 1*time.Second, v)

	if out.Len() == 0 {
		t.Error("Logger is empty")
	}
}
