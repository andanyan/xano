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
	aa := new(A)
	a(aa)
	fmt.Printf("%+v\n", aa)
}

func a(v interface{}) {
	s := b()
	fmt.Printf("%+v\n", s)

	fv := reflect.Indirect(reflect.ValueOf(v))
	fmt.Println(fv.CanAddr(), fv.CanInterface(), fv.CanSet())

	fs := reflect.ValueOf(s)

	fmt.Println("sss", fs.IsValid())
	if !fs.IsValid() {
		return
	}

	fv.Set(fs)

}

func b() interface{} {
	c := A{
		Id:   1,
		Name: "xlq",
	}
	//var c *A
	return c
}
