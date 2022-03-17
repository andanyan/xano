package main

import (
	"xano"
)

func main() {
	xano.WithConfig("./config/member1.toml")

	xano.Run()
}
