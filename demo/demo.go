package main

import (
	"fmt"
	"sync"
	"time"

	"github.com/chenjinya/go-idalloc"
)

func main() {
	var ida idalloc.Pool
	idcu := make(chan uint64)
	idca := make(chan uint64)

	ida.Run("user", idcu)
	ida.Run("article", idca)
	ida.BootAutoIncre(0)
	ida.Debug(true)
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		wg.Add(1)
		for {
			time.Sleep(time.Second)
			id := <-idcu
			fmt.Println(id)
		}

	}()
	go func() {
		wg.Add(1)
		for {
			time.Sleep(time.Second)
			id := <-idcu
			fmt.Println(id)
		}
	}()

	wg.Wait()
}
