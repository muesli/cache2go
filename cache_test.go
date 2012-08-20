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
	a.XCache("mama", 1*time.Second, a)
	b, err := GetXCached("mama")
	if err != nil || b == nil || b != a {
		t.Error("Error retriving data from cache", err)
	}
}

func TestCacheExpire(t *testing.T) {
	a := &myStruct{data: "mama are mere"}
	a.XCache("mama", 1*time.Second, a)
	b, err := GetXCached("mama")
	if err != nil || b == nil || b.(*myStruct).data != "mama are mere" {
		t.Error("Error retriving data from cache", err)
	}
	time.Sleep(1001 * time.Millisecond)
	b, err = GetXCached("mama")
	if err == nil || b != nil {
		t.Error("Error expiring data")
	}
}

func TestCacheKeepAlive(t *testing.T) {
	a := &myStruct{data: "mama are mere"}
	a.XCache("mama", 1*time.Second, a)
	b, err := GetXCached("mama")
	if err != nil || b == nil || b.(*myStruct).data != "mama are mere" {
		t.Error("Error retriving data from cache", err)
	}
	time.Sleep(500 * time.Millisecond)
	b.KeepAlive()
	time.Sleep(501 * time.Millisecond)
	if err != nil {
		t.Error("Error keeping cached data alive", err)
	}
	time.Sleep(1000 * time.Millisecond)
	b, err = GetXCached("mama")
	if err == nil || b != nil {
		t.Error("Error expiring data")
	}
}

func TestFlush(t *testing.T) {
	a := &myStruct{data: "mama are mere"}
	a.XCache("mama", 10*time.Second, a)
	time.Sleep(1000 * time.Millisecond)
	XFlush()
	b, err := GetXCached("mama")
	if err == nil || b != nil {
		t.Error("Error expiring data")
	}
}

func TestFlushNoTimout(t *testing.T) {
	a := &myStruct{data: "mama are mere"}
	a.XCache("mama", 10*time.Second, a)
	XFlush()
	b, err := GetXCached("mama")
	if err == nil || b != nil {
		t.Error("Error expiring data")
	}
}
