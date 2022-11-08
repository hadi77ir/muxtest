package server

import (
	"bytes"
	"fmt"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

var addr string
var wg *sync.WaitGroup

func Run(args []string) {
	addr = args[0]
	wg = &sync.WaitGroup{}
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		panic(err)
	}

	finChan := make(chan struct{}, 1)
	// Handle SIGINT and SIGTERM.
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	ch := make(chan struct{}, 1)

	wg.Add(1)
	go server(listener, ch)
	go func() {
		wg.Wait()
		close(finChan)
	}()
	select {
	case <-sigChan:
		_ = listener.Close()
		close(ch)
	case <-finChan:
	}
}

func server(listener net.Listener, done chan struct{}) {
	defer wg.Done()
	for {
		conn, err := listener.Accept()
		if err != nil {
			return
		}
		fmt.Println("received connection")
		wg.Add(1)
		go serveConn(conn)
		select {
		case <-done:
			return
		default:
		}
	}
}

const startingMessage = "hello world"
const endingMessage = "hello back"

func serveConn(conn net.Conn) {
	defer wg.Done()
	fmt.Println("handling connection")
	buf := []byte(startingMessage)
	n, err := conn.Read(buf)
	if err == nil && n == len(buf) && bytes.Equal([]byte(startingMessage), buf) {
		_, _ = conn.Write([]byte(endingMessage))
	}
	_ = conn.Close()
	fmt.Println("ended connection")
}
