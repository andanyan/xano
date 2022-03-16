package main

import (
	"xano"
)

func main() {
	xano.WithConfig("./config/master.toml")

	xano.Run()
}
