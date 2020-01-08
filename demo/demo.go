package main

import (
	"fmt"
	"sync"
	"time"

	"github.com/chenjinya/go-idalloc"
)

func main() {
	idcu := make(chan uint64)
	idca := make(chan uint64)
	pl := idalloc.Pool{}
	pl.Run("user", idcu).Run("article", idca)
	idalloc.BootAutoIncre(0)
	idalloc.Debug(true)
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
