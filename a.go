package main

import (
	"fmt"
	"reflect"
)

type A struct {
	Id   int
	Name string
}

func main() {
	a := A{
		1, "xlq",
	}

	at := reflect.TypeOf(a)

	na := reflect.New(at)
	ttt(na.Interface())
	fmt.Println(na)
}

func ttt(a interface{}) {
	b := a.(*A)
	b.Id = 2
	b.Name = "mh"
}
