package log

import (
	"fmt"
	"os"
)

func Println(v ...interface{}) {
	fmt.Println(v...)
}

func Fatal(v ...interface{}) {
	fmt.Println(v...)
	os.Exit(1)
}
