package main

import (
	"fmt"
)

type A struct {
	Id   int
	Name string
}

func main() {
	var aa = map[int]*A{
		1: {1, "xlq"},
		2: {2, "mh"},
	}
	fmt.Println(aa, aa[1])

}
