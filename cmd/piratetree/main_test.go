package main

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/giuseppemaniscalco/piratetree/internal/handler/booking"
	adapter "github.com/giuseppemaniscalco/piratetree/internal/handler/booking/adapter/windingtree"
	"github.com/giuseppemaniscalco/piratetree/internal/handler/booking/response"
	provider "github.com/giuseppemaniscalco/piratetree/internal/provider/windingtree"
)

type mockHttpClient struct {
	response *http.Response
	error    error
}

func init() {
	log.SetOutput(ioutil.Discard)
}

func TestMain(m *testing.M) {
	code := m.Run()
	os.Exit(code)
}

func fixtureBytes(t *testing.T, test string) []byte {
	currentDir, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	filePath := filepath.Join(currentDir, test)
	content, err := ioutil.ReadFile(filePath)
	if err != nil {
		t.Fatal(err)
	}

	return content
}

func fixtureReadCloser(t *testing.T, test string) io.ReadCloser {
	content := fixtureBytes(t, test)
	c := bytes.NewReader(content)

	return ioutil.NopCloser(c)
}

func (m *mockHttpClient) Do(*http.Request) (*http.Response, error) {
	return m.response, m.error
}

func TestBooking(t *testing.T) {
	dataProvider := map[string]struct {
		mockProviderResponse provider.HttpClient
		actualRequest        io.Reader
		expectedStatus       int
		expectedResponse     []byte
	}{
		"200": {
			mockProviderResponse: &mockHttpClient{
				response: &http.Response{
					StatusCode: http.StatusOK,
					Body:       fixtureReadCloser(t, "test/data/integration/booking/actualProviderResponse.json"),
				},
			},
			actualRequest:    fixtureReadCloser(t, "test/data/integration/booking/actualRequest.json"),
			expectedStatus:   http.StatusOK,
			expectedResponse: fixtureBytes(t, "test/data/integration/booking/expectedResponse.json"),
		},
	}
	for name, tt := range dataProvider {
		t.Run(name, func(t *testing.T) {
			//
			// Request/Response
			//
			actualRequest, err := http.NewRequest(http.MethodPost, "/booking", tt.actualRequest)
			if err != nil {
				t.Fatal(err)
			}
			actualResponse := httptest.NewRecorder()
			//
			// Dependency
			//
			wtProvider := provider.NewWindingTree(tt.mockProviderResponse, "")
			wtAdapter := adapter.NewWindingTree(wtProvider)
			bookingHandler := booking.NewHandler(wtAdapter)
			mux := http.NewServeMux()
			//
			// Router
			//
			mux.Handle("/booking", bookingHandler)
			//
			// Server
			//
			mux.ServeHTTP(actualResponse, actualRequest)
			//
			// Asserting
			//
			if actualResponse.Code != tt.expectedStatus {
				t.Errorf("status code want %d got %d", actualResponse.Code, tt.expectedStatus)
			}
			actualResponseOjb := new(response.Response)
			if err := json.Unmarshal(actualResponse.Body.Bytes(), actualResponseOjb); err != nil {
				t.Fatal(err)
			}
			expectedResponseObj := new(response.Response)
			if err := json.Unmarshal(tt.expectedResponse, expectedResponseObj); err != nil {
				t.Fatal(err)
			}
			if !reflect.DeepEqual(actualResponseOjb, expectedResponseObj) {
				t.Errorf("response want %s got %s", actualResponse.Body.String(), tt.expectedResponse)
			}
		})
	}
}
