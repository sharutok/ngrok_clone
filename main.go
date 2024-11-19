package main

import (
	"fmt"
	"log"
	"net/http/httputil"
	"net/url"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	_ "github.com/joho/godotenv/autoload"
)

// createProxy creates a reverse proxy to the given target URL
func createProxy(target string) *httputil.ReverseProxy {
	url, err := url.Parse(target)
	if err != nil {
		log.Fatalf("Failed to parse target URL %s: %v", target, err)
	}
	return httputil.NewSingleHostReverseProxy(url)
}

func reverseProxyHandler(c *gin.Context) {
	path := c.Request.URL.Path
	rdb := ConnectToRedis()

	res, err := rdb.HGetAll(ctx, "reverse_proxy_server").Result()
	if err != nil {
		log.Fatalf("Could not retrieve hash: %v", err)
	}

	app := strings.Split(path, "/")[1]

	for app_name, link := range res {
		s := fmt.Sprintf("/%s", app_name)
		if app == app_name {
			c.Request.URL.Path = strings.TrimPrefix(path, s)
			createProxy(link).ServeHTTP(c.Writer, c.Request)
			return
		}
	}
	log.Println("No match")
	c.JSON(404, gin.H{"error": "App not found"})
	c.Abort()
}

func main() {
	godotenv.Load(".env")
	router := gin.Default()
	router.NoRoute(reverseProxyHandler)
	router.POST("/add-app", ADD_APP)
	router.POST("/delete-app", DELETE_APP)
	router.GET("/health-check", HEALTH_CHECK)

	port := ":8081"
	log.Printf("Reverse proxy server is running on port %s", port)
	if err := router.Run(port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
