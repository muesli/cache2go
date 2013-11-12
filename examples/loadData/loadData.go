package main

import (
	"fmt"
	"github.com/rif/cache2go"
	"strconv"
)

type myStruct struct {
	text     string
	moreData []byte
}

func main() {
	cache := cache2go.Cache("myCache")
	cache.SetDataLoader(func(key interface{}) *cache2go.CacheItem {
		// Apply some clever loading logic here, e.g. read values for
		// this key from database, network or file
		val := myStruct{"This is a test with key " + key.(string), []byte{}}

		// This helper method creates the cached item for us. Yay!
		item := cache2go.CreateCacheItem(key, 0, &val)
		return &item
	})

	// Let's retrieve the item for key "someKey" from the cache
	for i := 0; i < 10; i++ {
		res, err := cache.Value("someKey_" + strconv.Itoa(i))
		if err == nil {
			fmt.Println("Found value in cache:", res.Data().(*myStruct).text)
		} else {
			fmt.Println("Error retrieving value from cache:", err)
		}
	}
}
