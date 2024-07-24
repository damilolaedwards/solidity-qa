package api

import (
	"assistant/config"
	"assistant/logging"
	"assistant/types"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/rs/zerolog"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func Start(config *config.ProjectConfig, slitherOutput *types.SlitherOutput) {
	var port string

	if config.Port == 0 {
		port = ":8080" // Default port
	} else {
		port = fmt.Sprint(":", config.Port)
	}

	// Create sub-logger for api module
	logger := logging.NewLogger(zerolog.InfoLevel)
	logger.AddWriter(os.Stdout, logging.UNSTRUCTURED, true)

	// Create a new router
	router := mux.NewRouter()

	// Attach middleware
	attachMiddleware(router)

	// Serve the contracts on a sub-router
	router.HandleFunc("/contracts", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		err := json.NewEncoder(w).Encode(slitherOutput)
		if err != nil {
			logger.Error("Failed to encode contracts: ", err)
			return
		}
	})

	var listener net.Listener
	var err error

	for i := 0; i < 10; i++ {
		listener, err = net.Listen("tcp", port)
		if err == nil {
			break
		}

		logger.Info("Server failed to start on port ", port[1:])
		port = incrementPort(port)
	}

	// Stop further execution if the server failed to start
	if listener == nil {
		logger.Error("Failed to start server: ", err)
		return
	}

	logger.Info("Server started on port ", port[1:])

	// Create a channel to receive interrupt signals
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	// Start the server in a separate goroutine
	serverErrorChan := make(chan error, 1)
	go func() {
		serverErrorChan <- http.Serve(listener, router)
	}()

	// Gracefully shutdown the server if a server error is encountered
	select {
	case <-sigChan:
		logger.Info("Shutting down server...")
		err := listener.Close()
		if err != nil {
			logger.Error("Failed to shut down server: ", err)
			return
		}
	case err := <-serverErrorChan:
		logger.Error("Server error: ", err)
	}
}

func incrementPort(port string) string {
	var portNum int

	_, err := fmt.Sscanf(port, ":%d", &portNum)
	if err != nil {
		panic(err)
	}

	return fmt.Sprintf(":%d", portNum+1)
}
