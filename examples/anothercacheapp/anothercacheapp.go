package main

import (
	"fmt"
	"github.com/ck89119/cache2go"
	"time"
)

// Keys & values in cache2go can be off arbitrary types, e.g. a struct.
type myStruct struct {
	text     string
	moreData []byte
}

func main() {
	// Accessing a new cache table for the first time will create it.
	cache := cache2go.Cache("myCache")

	// We will put a new item in the cache. It will expire after
	// not being accessed via Value(key) for more than 5 seconds.
	val1 := myStruct{"This is a test!", []byte{1}}
	cache.Add("1", 4*time.Second, &val1)
	val2 := myStruct{"This is a test!", []byte{2}}
	cache.Add("2", 6*time.Second, &val2)

	time.Sleep(3 * time.Second)

	res, err := cache.Value("1")
	if err == nil {
		fmt.Println("Found value in cache 1:", res.Data().(*myStruct).moreData)
	} else {
		fmt.Println("Error retrieving value from cache:", err)

	}

	time.Sleep(3500 * time.Millisecond)

	res, err = cache.Value("1")
	if err == nil {
		fmt.Println("Found value in cache 1:", res.Data().(*myStruct).moreData)
	} else {
		fmt.Println("Error retrieving value from cache:", err)
	}

	res, err = cache.Value("2")
	if err == nil {
		fmt.Println("Found value in cache 2:", res.Data().(*myStruct).moreData)
	} else {
		fmt.Println("Error retrieving value from cache:", err)
	}

	// Remove the item from the cache.
	cache.Delete("1")
	cache.Delete("2")

	// And wipe the entire cache table.
	cache.Flush()
}
