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
