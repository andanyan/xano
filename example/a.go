package main

import (
	"fmt"
	"time"
)

func main() {
	m := make(map[int]int)

	go func() {
		for i := 0; i < 1000; i++ {
			fmt.Println("aaaaaaa")
			m[i] = i
		}
	}()

	go func() {
		for i := 0; i < 1000; i++ {
			fmt.Println("bbbb", m[i])
		}
	}()

	time.Sleep(20 * time.Second)
}
