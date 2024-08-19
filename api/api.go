package api

import (
	"assistant/config"
	"assistant/logging"
	"assistant/types"
	"embed"
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/gorilla/mux"
	"github.com/rs/zerolog"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"
)

//go:embed assets/public
var staticAssets embed.FS

type API struct {
	targetContracts string
	contracts       []types.Contract
	logger          *logging.Logger
}

func InitializeAPI(contractCodes string, contracts []types.Contract) *API {
	// Create sub-logger for api module
	logger := logging.NewLogger(zerolog.InfoLevel)
	logger.AddWriter(os.Stdout, logging.UNSTRUCTURED, true)

	return &API{
		targetContracts: contractCodes,
		contracts:       contracts,
		logger:          logger,
	}
}

func (api *API) Start(projectConfig *config.ProjectConfig) {
	var port string

	if projectConfig.Port == 0 {
		port = ":8080" // Default port
	} else {
		port = fmt.Sprint(":", projectConfig.Port)
	}

	// Create sub-logger for api module
	logger := logging.NewLogger(zerolog.InfoLevel)
	logger.AddWriter(os.Stdout, logging.UNSTRUCTURED, true)

	// Create a new router
	router := mux.NewRouter()

	// Attach other middleware
	api.attachMiddleware(router)

	// Serve static content
	router.PathPrefix("/static/").HandlerFunc(serveStaticFilesHandler)

	// Attach routes
	api.attachRoutes(router)

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

	// Create new watcher
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()

	//err = watcher.Add("assistant.json")
	//if err != nil {
	//	log.Fatal(err)
	//}

	select {
	// Shutdown the server upon keyboard interrupt
	case <-sigChan:
		logger.Info("Shutting down server...")
		err := listener.Close()
		if err != nil {
			logger.Error("Failed to shut down server: ", err)
			return
		}
	// Gracefully shutdown the server if a server error is encountered
	case err := <-serverErrorChan:
		logger.Error("Server error: ", err)
		// Restart server if config file is modified
		//case event, ok := <-watcher.Events:
		//	if !ok {
		//		return
		//	}
		//	if event.Op&fsnotify.Write == fsnotify.Write {
		//		fmt.Println(fmt.Sprintf("%s modifed. Restarting server...", event.Name))
		//		err := listener.Close()
		//		if err != nil {
		//			logger.Error("Failed to shut down server: ", err)
		//			return
		//		}
		//
		//		workingDirectory, err := os.Getwd()
		//		if err != nil {
		//			logger.Error("Failed to obtain working directory", err)
		//			return
		//		}
		//
		//		configPath := filepath.Join(workingDirectory, "assistant.json")
		//		projectConfig, err = config.ReadProjectConfigFromFile(configPath)
		//		if err != nil {
		//			logger.Error("Failed to read project config: ", err)
		//			return
		//		}
		//		api.Start(projectConfig)
		//	}
	}
}

func (api *API) attachRoutes(router *mux.Router) {
	attachFrontendRoutes(router, api.contracts, api.targetContracts)
}

func (api *API) attachMiddleware(router *mux.Router) {
	// Handle cancelled requests
	router.Use(func(next http.Handler) http.Handler {
		return http.TimeoutHandler(next, 30*time.Second, "Request timed out")
	})
}

func serveStaticFilesHandler(w http.ResponseWriter, r *http.Request) {
	// Remove "/static/" prefix from the request path
	filePath := strings.TrimPrefix(r.URL.Path, "/static/")
	serveStaticFile(w, r, "assets/public/"+filePath)
}

func serveStaticFile(w http.ResponseWriter, r *http.Request, filePath string) {
	file, err := staticAssets.Open(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			fmt.Println("Not found: " + filePath)
			http.NotFound(w, r)
		} else {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
		return
	}
	defer file.Close()

	info, err := file.Stat()
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	if info.IsDir() {
		fmt.Println("Is dir: " + filePath)
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}

	content, err := io.ReadAll(file)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	contentType := http.DetectContentType(content)
	w.Header().Set("Content-Type", contentType)
	w.Header().Set("Content-Length", fmt.Sprintf("%d", len(content)))
	w.Header().Set("Last-Modified", info.ModTime().UTC().Format(http.TimeFormat))

	if r.Method != "HEAD" {
		_, err = w.Write(content)
		if err != nil {
			log.Printf("Error writing response: %v", err)
		}
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
