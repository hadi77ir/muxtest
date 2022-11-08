package main

import (
	"fmt"
	"os"
)

// commands
import (
	"muxtest/client"
	"muxtest/server"
)

func main() {
	if len(os.Args) == 1 {
		usage()
		return
	}
	args := os.Args[2:]
	switch os.Args[1] {
	case "server":
		if len(args) != 1 {
			usage()
			return
		}
		server.Run(args)
	case "client":
		if len(args) != 2 {
			usage()
			return
		}
		client.Run(args)
	}
}
func usage() {
	fmt.Println("muxtest: simple hello world tester for multiplexer")
	fmt.Println("usage: muxtest server HOST:PORT")
	fmt.Println("usage: muxtest client HOST:PORT COUNT")
}
