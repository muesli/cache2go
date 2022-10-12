package cache2go

import (
	"testing"
	"time"
)

func TestCacheTable_HAdd(t *testing.T) {
	cache := Cache("htable")
	var hs1 int8 = 100
	cache.HAdd("hstudents", 0*time.Second, "s1", hs1)
	hv, err := cache.HValue("hstudents", "s1")
	if err != nil || hv == nil || hv.(int8) != hs1 {
		t.Error("Error retrieving non expiring data from cache", err)
	}

	var ht1 int8 = 99
	cache.HAdd("hteacher", 1*time.Second, "t1", ht1)
	time.Sleep(2 * time.Second)
	htv, err := cache.HValue("hteacher", "t1")
	if err == nil || htv != nil {
		t.Error("Error retrieving non expiring data from cache", err)
	}
}

func TestCacheTable_HDelete(t *testing.T) {
	var hv interface{}
	var err error
	var hs1 int8 = 100

	cache := Cache("htable")
	cache.HAdd("hstudents", 0*time.Second, "s1", hs1)
	hv, err = cache.HValue("hstudents", "s1")
	if err != nil || hv == nil || hv.(int8) != hs1 {
		t.Error("Error retrieving non expiring data from cache", err)
	}

	err = cache.HDelete("hstudents", "s1")
	if err != nil {
		t.Error("Error delete hash key fail", err)
	}

	hv, err = cache.HValue("hstudents", "s1")
	if err == nil || hv != nil {
		t.Error("Error retrieving non expiring data from cache", err)
	}
}
