package main

import (
	"log"

	"github.com/gin-gonic/gin"
)

type Applicationlist []string

func CheckIfAppNameExist(c *gin.Context) {
	rdb := ConnectToRedis()

	res, err := rdb.HGetAll(ctx, "proxy_server_app").Result()
	if err != nil {
		log.Fatalf("Could not retrieve hash: %v", err)
		return
	}

	// if strings.HasPrefix(c.GetHeader("Authorization"), "Bearer") && os.Getenv("AUTH_TOKEN") == strings.Split(c.GetHeader("Authorization"), " ")[1] {
	app := c.Param("app")
	if _, exists := res[app]; exists {
		// value := strings.Split(res[app], ":")
		c.Next()
	} else {
		log.Println("No match")
		c.JSON(404, gin.H{"error": "App not found"})
		c.Abort()
	}
	// } else {
	// 	c.JSON(404, gin.H{"error": "no proper token"})
	// 	c.Abort()
	// }
}
