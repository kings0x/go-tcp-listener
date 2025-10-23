package main

import "sync"

func main() {
	addr := "127.0.0.1:8080"
	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		StartServer(addr)
	}()

	go func() {
		defer wg.Done()
		StartClient(addr)
	}()

	wg.Wait()
}
