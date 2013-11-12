package cache2go

import (
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
	table.Cache(k, 1*time.Second, v)
	p, err := table.Value(k)
	if err != nil || p == nil || p.Data().(string) != v {
		t.Error("Error retrieving data from cache", err)
	}
}

func TestCacheExpire(t *testing.T) {
	table := Cache("testExpire")
	table.Cache(k, 250*time.Millisecond, v)
	p, err := table.Value(k)
	if err != nil || p == nil || p.Data().(string) != v {
		t.Error("Error retrieving data from cache", err)
	}
	time.Sleep(500*time.Millisecond)
	p, err = table.Value(k)
	if err == nil || p != nil {
		t.Error("Error expiring data")
	}
}

func TestCacheNonExpiring(t *testing.T) {
	table := Cache("testNonExpiring")
	table.Cache(k, 0, v)
	time.Sleep(500*time.Millisecond)
	p, err := table.Value(k)
	if err != nil || p == nil || p.Data().(string) != v {
		t.Error("Error retrieving data from cache", err)
	}
}

func TestCacheKeepAlive(t *testing.T) {
	k2 := k + k
	v2 := v + v
	table := Cache("testKeepAlive")
	table.Cache(k, 250*time.Millisecond, v)
	table.Cache(k2, 750*time.Millisecond, v2)

	p, err := table.Value(k)
	if err != nil || p == nil || p.Data().(string) != v {
		t.Error("Error retrieving data from cache", err)
	}
	time.Sleep(50*time.Millisecond)
	p.KeepAlive()

	time.Sleep(450*time.Millisecond)
	p, err = table.Value(k)
	if err == nil || p != nil {
		t.Error("Error expiring data")
	}
	p, err = table.Value(k2)
	if err != nil || p == nil || p.Data().(string) != v2 {
		t.Error("Error retrieving data from cache", err)
	}
	time.Sleep(1*time.Second)
	p, err = table.Value(k2)
	if err == nil || p != nil {
		t.Error("Error expiring data")
	}
}

func TestDelete(t *testing.T) {
	table := Cache("testDelete")
	table.Cache(k, 0, v)
	p, err := table.Value(k)
	if err != nil || p == nil || p.Data().(string) != v {
		t.Error("Error retrieving data from cache", err)
	}
	table.Delete(k)
	p, err = table.Value(k)
	if err == nil || p != nil {
		t.Error("Error deleting data")
	}
}

func TestFlush(t *testing.T) {
	table := Cache("testFlush")
	table.Cache(k, 10*time.Second, v)
	time.Sleep(100*time.Millisecond)
	table.Flush()
	p, err := table.Value(k)
	if err == nil || p != nil {
		t.Error("Error expiring data")
	}
}

func TestFlushNoTimout(t *testing.T) {
	table := Cache("testFlushNoTimeout")
	table.Cache(k, 10*time.Second, v)
	table.Flush()
	p, err := table.Value(k)
	if err == nil || p != nil {
		t.Error("Error expiring data")
	}
}

func TestMassive(t *testing.T) {
	count := 100000
	table := Cache("testMassive")
	for i := 0; i < count; i++ {
		key := k + strconv.Itoa(i)
		table.Cache(key, 2*time.Second, v)
	}
	for i := 0; i < count; i++ {
		key := k + strconv.Itoa(i)
		p, err := table.Value(key)
		if err != nil || p == nil || p.Data().(string) != v {
			t.Error("Error retrieving data")
		}
	}
	if table.CacheCount() != count {
		t.Error("Data count mismatch")
	}
}

func TestDataLoader(t *testing.T) {
	table := Cache("testDataLoader")
	table.SetDataLoader(func(key interface{}) *CacheEntry{
		val := k + key.(string)
		entry := CreateCacheEntry(key, 500*time.Millisecond, val)
		return &entry
	})

	for i := 0; i < 10; i++ {
		key := k + strconv.Itoa(i)
		vp := k + key
		p, err := table.Value(key)
		if err != nil || p == nil || p.Data().(string) != vp {
			t.Error("Error validating data loader")
		}
	}
}

func TestCallbacks(t *testing.T) {
	addedKey := ""
	removedKey := ""

	table := Cache("testCallbacks")
	table.SetAddedItemCallback(func(item *CacheEntry) {
		addedKey = item.Key().(string)
	})
	table.SetAboutToDeleteItemCallback(func(item *CacheEntry) {
		removedKey = item.Key().(string)
	})

	table.Cache(k, 500*time.Millisecond, v)

	time.Sleep(250*time.Millisecond)
	if addedKey != k {
		t.Error("AddedItem callback not working")
	}

	time.Sleep(500*time.Millisecond)
	if removedKey != k {
		t.Error("AboutToDeleteItem callback not working:" + k + "_" + removedKey)
	}
}
