package main

import (
	"fmt"

	"github.com/kiselev-nikolay/direct-to-me/pkg/server"
)

func main() {
	server.Serve("127.0.0.1", "8080", func(method, url string) {
		fmt.Println("method:", method)
		fmt.Println("url:   ", url)
	})
}
