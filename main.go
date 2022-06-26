package main

import (
	"fmt"
	"io"
	"net"
	"sync"
)

type Backend struct {
	Host string
	Port int
}

var (
	requestCounter int
	mut            sync.Mutex
)

var backends []*Backend = []*Backend{
	{Host: "127.0.0.1", Port: 5001},
	{Host: "127.0.0.1", Port: 5002},
}

func main() {
	lb, err := net.Listen("tcp", "localhost:8080")
	fmt.Println("Welcome to the LoadBalancer")
	if err != nil {
		panic(err.Error())
	}

	defer lb.Close()

	for {
		conn, err := lb.Accept()
		if err != nil {
			panic(err.Error())
		}

		go proxy(conn)
	}
}

func proxy(conn net.Conn) {
	backend := getNextBackend()
	backendConnection, err := net.Dial("tcp", fmt.Sprintf("%s:%d", backend.Host, backend.Port))
	if err != nil {
		backendConnection.Close()
		panic(err.Error())
	}
	go io.Copy(backendConnection, conn)
	go io.Copy(conn, backendConnection)
}

func getNextBackend() *Backend {
	mut.Lock()
	roundRobinIndex := requestCounter % len(backends)
	requestCounter += 1
	mut.Unlock()
	return backends[roundRobinIndex]
}
