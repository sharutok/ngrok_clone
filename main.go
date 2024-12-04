package main

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
)

func ReverseProxyHandler(w http.ResponseWriter, r *http.Request) {

	rdb := ConnectToRedis()
	res, err := rdb.HGetAll(ctx, "reverse_proxy_server").Result()
	if err != nil {
		log.Fatalf("Could not retrieve hash: %v", err)
		return
	}

	for prefix, target := range res {
		if strings.HasPrefix(r.URL.Path, prefix) {
			// Parse the target URL
			targetURL, err := url.Parse(target)
			if err != nil {
				http.Error(w, "Invalid target URL", http.StatusInternalServerError)
				log.Printf("Error parsing target URL: %v\n", err)
				return
			}

			// Create a reverse proxy for the matched backend
			proxy := httputil.NewSingleHostReverseProxy(targetURL)

			// Rewrite the path for the backend (optional)
			r.URL.Path = strings.TrimPrefix(r.URL.Path, prefix)

			// Forward the request
			proxy.ServeHTTP(w, r)
			return
		}
	}

	// If no route matches, return a 404
	http.Error(w, "Route not found", http.StatusNotFound)
}

func main() {
	// Create a reverse proxy handler
	reverseProxyHandler := http.HandlerFunc(ReverseProxyHandler)

	// Wrap the handler with the logging middleware
	loggedHandler := (reverseProxyHandler)

	// Start the server
	port := "8888"
	log.Printf("Reverse proxy server running on port %s\n", port)

	if err := http.ListenAndServe(":"+port, loggedHandler); err != nil {
		log.Fatalf("Failed to start server: %v\n", err)
	}
}
