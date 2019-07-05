package main

import (
	"log"
	"net/http"
	"os"

	"github.com/giuseppemaniscalco/piratetree/internal/handler/booking"
	adapter "github.com/giuseppemaniscalco/piratetree/internal/handler/booking/adapter/windingtree"
	provider "github.com/giuseppemaniscalco/piratetree/internal/provider/windingtree"
)

func main() {
	//
	// Environments
	//
	port := os.Getenv("PORT")
	if len(port) == 0 {
		log.Fatal("$PORT must be set")
	}
	windingTreeUrl := os.Getenv("WINDINGTREE_URL")
	if len(windingTreeUrl) == 0 {
		log.Fatal("$WINDINGTREE_URL must be set")
	}
	//
	// Dependencies
	//
	httpClient := http.DefaultClient
	wtProvider := provider.NewWindingTree(httpClient, windingTreeUrl)
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
	log.Println("HTTP server listen on port: " + port)
	log.Fatal(http.ListenAndServe(":"+port, mux))
}
