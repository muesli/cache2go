package cache2go

import (
	"strconv"
	"testing"
	"time"
	_ "fmt"
)

type myStruct struct {
	XEntry
	data string
}

func TestCache(t *testing.T) {
	a := &myStruct{data: "testy test"}
	a.XCache("test", 1*time.Second, a, nil)
	b, err := GetXCached("test")
	if err != nil || b == nil || b != a {
		t.Error("Error retrieving data from cache", err)
	}
}

func TestCacheExpire(t *testing.T) {
	a := &myStruct{data: "testy test"}
	a.XCache("test", 1*time.Second, a, nil)
	b, err := GetXCached("test")
	if err != nil || b == nil || b.(*myStruct).data != "testy test" {
		t.Error("Error retrieving data from cache", err)
	}
	time.Sleep(1500 * time.Millisecond)
	b, err = GetXCached("test")
	if err == nil || b != nil {
		t.Error("Error expiring data")
	}
}

func TestCacheKeepAlive(t *testing.T) {
	a := &myStruct{data: "testy test"}
	a.XCache("test", 500*time.Millisecond, a, nil)
	a = &myStruct{data: "testest test"}
	a.XCache("test2", 1250*time.Millisecond, a, nil)
	b, err := GetXCached("test")
	if err != nil || b == nil || b.(*myStruct).data != "testy test" {
		t.Error("Error retrieving data from cache", err)
	}
	time.Sleep(200 * time.Millisecond)
	b.KeepAlive()
	time.Sleep(750 * time.Millisecond)
	b, err = GetXCached("test")
	if err == nil || b != nil {
		t.Error("Error expiring data")
	}
	b, err = GetXCached("test2")
	if err != nil || b == nil || b.(*myStruct).data != "testest test" {
		t.Error("Error retrieving data from cache", err)
	}
	time.Sleep(1500 * time.Millisecond)
	b, err = GetXCached("test2")
	if err == nil || b != nil {
		t.Error("Error expiring data")
	}
}

func TestFlush(t *testing.T) {
	a := &myStruct{data: "testy test"}
	a.XCache("test", 10*time.Second, a, nil)
	XFlush()
	b, err := GetXCached("test")
	if err == nil || b != nil {
		t.Error("Error expiring data")
	}
}

func TestFlushNoTimout(t *testing.T) {
	a := &myStruct{data: "testy test"}
	a.XCache("test", 10*time.Second, a, nil)
	XFlush()
	b, err := GetXCached("test")
	if err == nil || b != nil {
		t.Error("Error expiring data")
	}
}

func TestMassive(t *testing.T) {
	return

	for i := 0; i < 1000000; i++ {
		a := &myStruct{data: "testy test"}
		key := "test_" + strconv.Itoa(i)
		a.XCache(key, 1*time.Second, a, nil)
	}

//	fmt.Println("In Cache:", XCacheCount())
	time.Sleep(500 * time.Millisecond)
//	fmt.Println("In Cache:", XCacheCount())
	time.Sleep(1500 * time.Millisecond)
//	fmt.Println("In Cache:", XCacheCount())
}
