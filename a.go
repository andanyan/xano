package main

import (
	"fmt"
	"time"
)

func main() {
	c := make(chan struct{})

	go func() {
		time.Sleep(2 * time.Second)
		c <- struct{}{}
	}()

	select {
	case <-c:
		fmt.Println("call successfully!!!")
		return
	case <-time.After(time.Duration(3 * time.Second)):
		fmt.Println("timeout!!!")
		return
	}

}
