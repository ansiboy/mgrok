package main

import (
	"mgrok/server"
	_ "net/http/pprof"
)

func main() {
	server.Main()
}
