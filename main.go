package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	socketPathPtr := flag.String("socket", "/tmp/auth.sock", "path to the UDS that serves the creds")
	flag.Parse()
	socketPath := *socketPathPtr

	// 1. Read JSON credentials from stdin
	fmt.Fprintln(os.Stderr, "Reading credentials from stdin... (Ctrl+D to finish)")
	creds, err := io.ReadAll(os.Stdin)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading stdin: %v\n", err)
		os.Exit(1)
	}

	if len(creds) == 0 {
		fmt.Fprintln(os.Stderr, "Error: No credentials provided.")
		os.Exit(1)
	}

	// 2. Clean up existing socket file if it exists
	if _, err := os.Stat(socketPath); err == nil {
		os.Remove(socketPath)
	}

	// 3. Create the Unix Domain Socket listener
	listener, err := net.Listen("unix", socketPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating socket: %v\n", err)
		os.Exit(1)
	}
	defer os.Remove(socketPath)

	// Set permissions so the socket is accessible (adjust as needed)
	os.Chmod(socketPath, 0600)

	// 4. Define the HTTP Handler
	mux := http.NewServeMux()
	mux.HandleFunc("/token", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(creds)
	})

	server := &http.Server{Handler: mux}

	// 5. Handle graceful shutdown (Ctrl+C)
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	go func() {
		fmt.Printf("Serving credentials on %s\n", socketPath)
		if err := server.Serve(listener); err != nil && err != http.ErrServerClosed {
			fmt.Fprintf(os.Stderr, "Server error: %v\n", err)
		}
	}()

	<-stop
	fmt.Println("\nShutting down...")
	server.Shutdown(context.Background())
}
