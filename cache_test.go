package cache2go

import (
	"testing"
	"time"
)

type myStruct struct {
	XEntry
	data string
}

func TestCache(t *testing.T) {
	a := &myStruct{data: "mama are mere"}
	a.XCache("mama", 1*time.Second, a, nil)
	b, err := GetXCached("mama")
	if err != nil || b == nil || b != a {
		t.Error("Error retrieving data from cache", err)
	}
}

func TestCacheExpire(t *testing.T) {
	a := &myStruct{data: "mama are mere"}
	a.XCache("mama", 1*time.Second, a, nil)
	b, err := GetXCached("mama")
	if err != nil || b == nil || b.(*myStruct).data != "mama are mere" {
		t.Error("Error retrieving data from cache", err)
	}
	time.Sleep(1001 * time.Millisecond)
	b, err = GetXCached("mama")
	if err == nil || b != nil {
		t.Error("Error expiring data")
	}
}

func TestCacheKeepAlive(t *testing.T) {
	a := &myStruct{data: "mama are mere"}
	a.XCache("mama", 1*time.Second, a, nil)
	a = &myStruct{data: "mama are mere2"}
	a.XCache("mama2", 2*time.Second, a, nil)
	b, err := GetXCached("mama")
	if err != nil || b == nil || b.(*myStruct).data != "mama are mere" {
		t.Error("Error retrieving data from cache", err)
	}
	time.Sleep(500 * time.Millisecond)
	b.KeepAlive()
	time.Sleep(1001 * time.Millisecond)
	b, err = GetXCached("mama")
	if err == nil || b != nil {
		t.Error("Error expiring data")
	}
	b, err = GetXCached("mama2")
	if err != nil || b == nil || b.(*myStruct).data != "mama are mere2" {
		t.Error("Error retrieving data from cache", err)
	}
	time.Sleep(2001 * time.Millisecond)
	b, err = GetXCached("mama2")
	if err == nil || b != nil {
		t.Error("Error expiring data")
	}
}

func TestFlush(t *testing.T) {
	a := &myStruct{data: "mama are mere"}
	a.XCache("mama", 10*time.Second, a, nil)
	XFlush()
	b, err := GetXCached("mama")
	if err == nil || b != nil {
		t.Error("Error expiring data")
	}
}

func TestFlushNoTimout(t *testing.T) {
	a := &myStruct{data: "mama are mere"}
	a.XCache("mama", 10*time.Second, a, nil)
	XFlush()
	b, err := GetXCached("mama")
	if err == nil || b != nil {
		t.Error("Error expiring data")
	}
}
