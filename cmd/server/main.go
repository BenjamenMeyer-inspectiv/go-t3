package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/BenjamenMeyer-inspectiv/go-t3/internal/api"
)

func main() {
	var port int
	flag.IntVar(&port, "port", 8080, "port to listen on")
	flag.Parse()

	addr := fmt.Sprintf(":%d", port)
	fmt.Printf("go-t3 server starting on %s\n", addr)
	router := api.NewRouter()
	log.Fatal(http.ListenAndServe(addr, router))
}
