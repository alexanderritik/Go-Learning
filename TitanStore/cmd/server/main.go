package main

import (
	"bufio"
	"fmt"
	"net"
	"strings"
	"sync"

	"github.com/alexanderritik/TitanStore/internal/store"
)

var clients sync.Map
var cache *store.Cache

func main() {

	cache = store.NewCache(4)

	listener, err := net.Listen("tcp", ":9000")
	if err != nil {
		fmt.Printf("Error")
	}

	for {
		conn, err := listener.Accept()
		if err != nil {
			continue
		}
		clients.Store(conn, true)
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()
	defer clients.Delete(conn)

	reader := bufio.NewReader(conn)

	for {
		message, err := reader.ReadString('\n')
		if err != nil {
			return
		}
		splitMessage := strings.Split(strings.TrimSpace(message), " ")
		switch splitMessage[0] {
		case "SET":
			fmt.Printf("SET ")
			cache.Set(splitMessage[1], splitMessage[2])
			conn.Write([]byte("OK\n"))
		case "GET":
			fmt.Printf("GET ")
			res, bool := cache.Get(splitMessage[1])
			if bool {
				result := fmt.Sprintf("%v", res)
				conn.Write([]byte(result))
			} else {
				conn.Write([]byte("Not Found\n"))
			}

		}
	}

}
