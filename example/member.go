package main

import (
	"xano"
)

func main() {
	xano.WithConfig("./config/member.toml")

	xano.Run()
}
