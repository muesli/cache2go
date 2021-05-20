package cache2go

import (
	"testing"
	"time"
)

func TestCacheTable_SAdd_SIsMember_SDelete(t *testing.T) {
	table := Cache("newset")
	for i := 1; i <= 10; i++ {
		table.SAdd("newset", 0*time.Second, i)
	}
	for i := 1; i <= 10; i++ {
		if has := table.SIsMember("newset", i); !has {
			t.Error("SAdd  i is", i, "err is", has)
		}
	}
	for i := 1; i <= 10; i++ {
		table.SDelete("newset", i)
	}
	for i := 1; i <= 10; i++ {
		if has := table.SIsMember("newset", i); has {
			t.Error("SDelete  i is", i, "err is", has)
		}
	}
}
