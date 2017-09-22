package main

import (
	"fmt"
	"github.com/danman113/simpleServer/server"
)

func main() {
	fmt.Printf("Hai\n")
	s := server.NewServer()
	s.AddStaticPage("/", "html/index.html")
	s.AddStaticFileserver("/static/", "static")
	s.Start(8080)
}
