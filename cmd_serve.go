package main

import (
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

type ServeCmd struct {
	BindAddr string `env:"BIND_ADDR" help:"Port to listen on" default:"0.0.0.0:1984"`
	Input    string `env:"INPUT_FILE" help:"Input file to read base configuration from" default:"./input.yaml"`
	Rules    string `env:"RULES_FILE" help:"Rules file" default:"./rules.yaml"`
}

func (s *ServeCmd) Run(cli *CLI) error {
	mux := http.NewServeMux()
	mux.HandleFunc("/", inputTransformHandler(s))

	server := &http.Server{Addr: ":9200", Handler: mux}

	// Handle shutdown
	errChan := make(chan error, 1)
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		ln, err := net.Listen("tcp", s.BindAddr)
		if err != nil {
			errChan <- fmt.Errorf("failed to bind to %s: %w", s.BindAddr, err)
			return
		}

		fmt.Printf("✅ Server is up and listening on http://%s\n", s.BindAddr)

		if err := server.Serve(ln); err != nil {
			if err != http.ErrServerClosed {
				errChan <- fmt.Errorf("server error: %w", err)
				return
			} else {
				errChan <- nil
				return
			}
		}
	}()

	select {
	case <-stop:
		return server.Close()

	case err := <-errChan:
		if err != nil {
			return err
		}
		return nil
	}
}
