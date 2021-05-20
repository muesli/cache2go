package cache2go

import (
	"testing"
	"time"
)

func TestCacheTable_LPush_LPop_RPush_RPop_ListLength(t *testing.T) {
	table := Cache("newlist")
	var err error
	var j interface{}
	for i := 1; i <= 10; i++ {
		err = table.LPush("newlist", 0*time.Second, i)
		if err != nil {
			t.Error("lpush error", err)
		}
	}
	if length, err := table.ListLength("newlist"); length != 10 || err != nil {
		t.Error("length error", err)
	}

	for i := 10; i >= 1; i-- {
		j, err = table.LPop("newlist")
		jInt, ok := j.(int)
		if err != nil || i != jInt || !ok {
			t.Error("lpop error", err)
		}
	}

	if length, err := table.ListLength("newlist"); length != 0 || err != nil {
		t.Error("length error", err)
	}

	for i := 1; i <= 10; i++ {
		err = table.RPush("newlist", 0*time.Second, i)
		if err != nil {
			t.Error("lpush error", err)
		}
	}

	for i := 10; i >= 1; i-- {
		j, err = table.RPop("newlist")
		jInt, ok := j.(int)
		if err != nil || i != jInt || !ok {
			t.Error("lpop error", err)
		}
	}

}
