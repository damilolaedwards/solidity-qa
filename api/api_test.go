package api

import (
	"assistant/types"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
)

func TestInitializeAPI(t *testing.T) {
	contracts := []types.Contract{{Name: "TestContract"}}
	api := InitializeAPI("testContractCodes", contracts)

	if api == nil {
		t.Fatal("Expected non-nil API instance")
	}

	if api.contractCodes != "testContractCodes" {
		t.Errorf("Expected targetContracts to be 'testContractCodes', got %s", api.contractCodes)
	}

	if len(api.contracts) != 1 || api.contracts[0].Name != "TestContract" {
		t.Errorf("Expected contracts to contain one TestContract, got %v", api.contracts)
	}

	if api.logger == nil {
		t.Error("Expected non-nil logger")
	}
}

func TestAttachRoutes(t *testing.T) {
	api := &API{
		contractCodes: "testContracts",
		contracts:     []types.Contract{{Name: "TestContract"}},
		projectName:   "TestProject",
	}

	router := mux.NewRouter()
	err := api.attachRoutes(router)

	if err != nil {
		t.Fatalf("attachRoutes returned an error: %v", err)
	}

	// Test that at least one route was attached by making a request
	// and checking if it's handled (not resulting in a 404)
	req, _ := http.NewRequest("GET", "/", nil)
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	if rr.Code == http.StatusNotFound {
		t.Error("Expected routes to be attached, but got a 404 response")
	}
}

func TestAttachMiddleware(t *testing.T) {
	api := &API{}
	router := mux.NewRouter()

	api.attachMiddleware(router)

	// Create a test server using the router
	server := httptest.NewServer(router)
	defer server.Close()

	// Send a request to the server
	resp, err := http.Get(server.URL)
	if err != nil {
		t.Fatalf("Failed to send request: %v", err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Println("Error closing response body: ", err)
		}
	}(resp.Body)
}

func TestServeStaticFilesHandler(t *testing.T) {
	// Create a request to /static/test.txt
	req, err := http.NewRequest("GET", "/static/test.txt", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Create a ResponseRecorder to record the response
	rr := httptest.NewRecorder()

	// Call the handler function
	handler := http.HandlerFunc(serveStaticFilesHandler)
	handler.ServeHTTP(rr, req)

	// Check the status code
	if status := rr.Code; status != http.StatusNotFound {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusNotFound)
	}

	// Send request for a file that exists
	req, err = http.NewRequest("GET", "/static/css/style.css", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Create a ResponseRecorder to record the response
	rr = httptest.NewRecorder()

	// Call the handler function
	handler = http.HandlerFunc(serveStaticFilesHandler)
	handler.ServeHTTP(rr, req)

	// Check the status code
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	// Check the content type
	if contentType := rr.Header().Get("Content-Type"); contentType != "text/plain; charset=utf-8" {
		t.Errorf("handler returned wrong content type: got %v want %v", contentType, "text/plain; charset=utf-8")
	}
}

func TestIncrementPort(t *testing.T) {
	testCases := []struct {
		input    string
		expected string
	}{
		{":8080", ":8081"},
		{":9000", ":9001"},
		{":3000", ":3001"},
	}

	for _, tc := range testCases {
		result := incrementPort(tc.input)
		if result != tc.expected {
			t.Errorf("incrementPort(%s) = %s; want %s", tc.input, result, tc.expected)
		}
	}
}
