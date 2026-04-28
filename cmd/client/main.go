package main

import (
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
)

func main() {
	var host string
	var port int
	flag.StringVar(&host, "host", "localhost", "server host")
	flag.IntVar(&port, "port", 8080, "server port")
	flag.Parse()

	args := flag.Args()
	if len(args) == 0 {
		fmt.Fprintf(os.Stderr, "usage: cli [flags] <command> [subcommand]\n")
		fmt.Fprintf(os.Stderr, "commands: ping, ping reset, pong\n")
		os.Exit(1)
	}

	baseURL := fmt.Sprintf("http://%s:%d", host, port)

	fmt.Printf(os.Stdout, "Server: %s\n", baseURL)
	fmt.Printf(os.Stdout, "-----------\n")
	fmt.Printf(os.Stdout, "Test Client\n")
	fmt.Printf(os.Stdout, "-----------\n")

	var path string
	switch args[0] {
	case "ping":
		if len(args) > 1 && args[1] == "reset" {
			path = "/ping/reset"
		} else {
			path = "/ping"
		}
	case "pong":
		path = "/pong"
	default:
		fmt.Fprintf(os.Stderr, "unknown command: %s\n", args[0])
		fmt.Fprintf(os.Stderr, "commands: ping, ping reset, pong\n")
		os.Exit(1)
	}

	resp, err := http.Get(baseURL + path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "request failed: %v\n", err)
		os.Exit(1)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to read response: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("status: %s\n", resp.Status)
	fmt.Printf("body: %s\n", body)
}
