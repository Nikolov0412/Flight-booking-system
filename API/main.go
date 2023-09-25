package main

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	r.Run()
	db, err := initDynamoDB()
	if err != nil {
		log.Fatalf("Error initializing DynamoDB client: %v", err)
	}
	fmt.Println(db.Endpoint)
}
