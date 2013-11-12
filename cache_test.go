package cache2go

import (
	"strconv"
	"testing"
	"time"
)

func TestCache(t *testing.T) {
	table := Cache("table")
	a := "testy test"
	table.Cache("test", 1*time.Second, a, nil)
	b, err := table.Value("test")
	if err != nil || b == nil || b.Data().(string) != a {
		t.Error("Error retrieving data from cache", err)
	}
}

func TestCacheExpire(t *testing.T) {
	table := Cache("table")
	a := "testy test"
	table.Cache("test", 1*time.Second, a, nil)
	b, err := table.Value("test")
	if err != nil || b == nil || b.Data().(string) != a {
		t.Error("Error retrieving data from cache", err)
	}
	time.Sleep(1500 * time.Millisecond)
	b, err = table.Value("test")
	if err == nil || b != nil {
		t.Error("Error expiring data")
	}
}

func TestCacheNonExpiring(t *testing.T) {
	table := Cache("table")
	a := "testy test"
	table.Cache("test", 0, a, nil)
	time.Sleep(500 * time.Millisecond)
	b, err := table.Value("test")
	if err != nil || b == nil || b.Data().(string) != a {
		t.Error("Error retrieving data from cache", err)
	}
}

func TestCacheKeepAlive(t *testing.T) {
	table := Cache("table")
	a := "testy test"
	table.Cache("test", 500*time.Millisecond, a, nil)
	a = "testest test"
	table.Cache("test2", 1250*time.Millisecond, a, nil)
	b, err := table.Value("test")
	if err != nil || b == nil || b.Data().(string) != "testy test" {
		t.Error("Error retrieving data from cache", err)
	}
	time.Sleep(200 * time.Millisecond)
	b.KeepAlive()
	time.Sleep(750 * time.Millisecond)
	b, err = table.Value("test")
	if err == nil || b != nil {
		t.Error("Error expiring data")
	}
	b, err = table.Value("test2")
	if err != nil || b == nil || b.Data().(string) != "testest test" {
		t.Error("Error retrieving data from cache", err)
	}
	time.Sleep(1500 * time.Millisecond)
	b, err = table.Value("test2")
	if err == nil || b != nil {
		t.Error("Error expiring data")
	}
}

func TestFlush(t *testing.T) {
	table := Cache("table")
	a := "testy test"
	table.Cache("test", 10*time.Second, a, nil)
	time.Sleep(1000 * time.Millisecond)
	table.Flush()
	b, err := table.Value("test")
	if err == nil || b != nil {
		t.Error("Error expiring data")
	}
}

func TestFlushNoTimout(t *testing.T) {
	table := Cache("table")
	a := "testy test"
	table.Cache("test", 10*time.Second, a, nil)
	table.Flush()
	b, err := table.Value("test")
	if err == nil || b != nil {
		t.Error("Error expiring data")
	}
}

func TestMassive(t *testing.T) {
	table := Cache("table")
	val := "testy test"
	count := 100000
	for i := 0; i < count; i++ {
		key := "test_" + strconv.Itoa(i)
		table.Cache(key, 2*time.Second, val, nil)
	}
	for i := 0; i < count; i++ {
		key := "test_" + strconv.Itoa(i)
		d, err := table.Value(key)
		if err != nil || d == nil || d.Data().(string) != val {
			t.Error("Error retrieving data")
		}
	}
	if table.CacheCount() != count {
		t.Error("Data count mismatch")
	}
}

func TestDataLoader(t *testing.T) {
	table := Cache("dataLoaderTest")
	table.SetDataLoader(func(key interface{}) *CacheEntry{
		val := "test_" + key.(string)
		entry := CreateCacheEntry(key, 500*time.Millisecond, val, nil)
		return &entry
	})

	for i := 0; i < 10; i++ {
		key := "test_" + strconv.Itoa(i)
		vp := "test_" + key
		p, err := table.Value(key)
		if err != nil || p == nil || p.Data().(string) != vp {
			t.Error("Error validating data loader")
		}
	}
}

func TestCallbacks(t *testing.T) {
	table := Cache("callbackTest")

	addedKey := ""
	removedKey := ""

	table.SetAddedItemCallback(func(item *CacheEntry) {
		addedKey = item.Key().(string)
	})
	table.SetAboutToDeleteItemCallback(func(item *CacheEntry) {
		removedKey = item.Key().(string)
	})

	k := "testkey"
	v := "testvalue"
	table.Cache(k, 500*time.Millisecond, v, nil)

	time.Sleep(250 * time.Millisecond)
	if addedKey != k {
		t.Error("AddedItem callback not working")
	}

	time.Sleep(500 * time.Millisecond)
	if removedKey != k {
		t.Error("AboutToDeleteItem callback not working:" + k + "_" + removedKey)
	}
}
