package api

import (
	"assistant/config"
	"assistant/logging"
	"assistant/logging/colors"
	"assistant/types"
	"context"
	"embed"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"

	"golang.ngrok.com/ngrok"
	ngrokConfig "golang.ngrok.com/ngrok/config"

	"github.com/fsnotify/fsnotify"
	"github.com/gorilla/mux"
	"github.com/rs/zerolog"
)

//go:embed assets/public
var staticAssets embed.FS

type API struct {
	targetContracts string
	contracts       []types.Contract
	logger          *logging.Logger
	projectName     string
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

	// Set project name
	api.projectName = projectConfig.Name

	// Create sub-logger for api module
	logger := logging.NewLogger(zerolog.InfoLevel)
	logger.AddWriter(os.Stdout, logging.UNSTRUCTURED, true)

	// Create a new router
	router := mux.NewRouter()

	// Serve static content
	router.PathPrefix("/static/").HandlerFunc(serveStaticFilesHandler)

	// Attach routes
	err := api.attachRoutes(router)
	if err != nil {
		logger.Error("Failed to start server: ", err)
		return
	}

	// Attach middleware
	api.attachMiddleware(router)

	var listener net.Listener

	// Set up local server
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
		logger.Error("Failed to start local server: ", err)
		return
	}

	var remoteUrl string

	// Host server using ngrok if specified in config
	if projectConfig.Host {
		logger.Info("Hosting server using ngrok")
		remoteUrl, err = startNgrok(projectConfig.Port)
		if err != nil {
			logger.Error(err)
			return
		}
	}

	// Log server details
	logger.Info("Server started successfully\n")
	logger.Info("Local: http://localhost" + port)

	if remoteUrl != "" {
		logger.Info("Remote: ", remoteUrl)
	} else {
		logger.Info("Remote: Use ", colors.Green, "--host", colors.White, " or the ", colors.Green, "host", colors.White, " config option to host the server (requires the NGROK_AUTHTOKEN env variable)")
	}

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
	}
}

// startNgrok starts a ngrok server listening to the server at port and returns the URL
func startNgrok(port int) (string, error) {
	// Obtain ngrok auth token
	ngrokAuthtoken := os.Getenv("NGROK_AUTHTOKEN")
	if ngrokAuthtoken == "" {
		return "", fmt.Errorf("NGROK_AUTHTOKEN environment variable is not set")
	}

	localServerUrl, err := url.Parse("http://localhost:" + strconv.Itoa(port))
	if err != nil {
		return "", fmt.Errorf("unable to parse local server URL: %v", err)
	}

	// Create ngrok listener
	ctx := context.Background()
	listener, err := ngrok.ListenAndForward(ctx,
		localServerUrl,
		ngrokConfig.HTTPEndpoint(ngrokConfig.WithRequestHeader("ngrok-skip-browser-warning", "true")),
		ngrok.WithAuthtokenFromEnv(),
	)
	if err != nil {
		return "", fmt.Errorf("failed to create ngrok listener: %v", err)
	}

	return listener.URL(), nil
}

func (api *API) attachRoutes(router *mux.Router) error {
	return attachFrontendRoutes(router, api.contracts, api.targetContracts, api.projectName)
}

func (api *API) attachMiddleware(router *mux.Router) {
	// Handle cancelled requests
	router.Use(func(next http.Handler) http.Handler {
		return http.TimeoutHandler(next, 300*time.Second, "Request timed out/was cancelled")
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
