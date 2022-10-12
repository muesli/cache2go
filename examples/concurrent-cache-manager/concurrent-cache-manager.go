package main

import (
	"fmt"
	"time"

	"github.com/muesli/cache2go"
)

func main() {
	store := cache2go.Cache("test")
	store.AddAboutToDeleteItemCallback(func(v *cache2go.CacheItem) {
		fmt.Println("delete:", v)
	})

	key := "test1"
	store.Add(key, time.Second*1800, "asdfasdfasd")
	fmt.Println("fitst delete:")
	store.Delete(key)

	store.Add(key, time.Second*1800, "asdfasdfasd")

	go func() {
		time.Sleep(time.Second * 2)
		fmt.Println("goroutine delete:")
		store.Foreach(func(key interface{}, v *cache2go.CacheItem) {
			fmt.Println("start")
			store.Delete(key)
			fmt.Println("end")
		})
	}()

	time.Sleep(time.Second * 10)
}
