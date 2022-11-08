package client

import (
	"bytes"
	"fmt"
	"net"
	"strconv"
	"sync"
	"sync/atomic"
)

var addr string
var result *atomic.Int32
var wg *sync.WaitGroup

func Run(args []string) {
	addr = args[0]
	times, _ := strconv.Atoi(args[1])
	wg = &sync.WaitGroup{}
	result = &atomic.Int32{}

	for i := 0; i < times; i++ {
		fmt.Println("sending connection number", i)
		wg.Add(1)
		go dialTarget(addr)
	}
	wg.Wait()
	fmt.Println("result is", result.Load())
}

const startingMessage = "hello world"
const endingMessage = "hello back"

func dialTarget(addr string) {
	defer wg.Done()
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return
	}

	_, err = conn.Write([]byte(startingMessage))
	if err != nil {
		return
	}

	buf := []byte(endingMessage)
	n, err := conn.Read(buf)
	if err == nil && n == len(buf) && bytes.Equal([]byte(endingMessage), buf) {
		result.Add(1)
	}
	_ = conn.Close()
}
