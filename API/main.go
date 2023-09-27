package main

import (
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/gin-gonic/gin"
)

var svc *dynamodb.DynamoDB

func init() {
	var err error
	svc, err = initDynamoDB()
	if err != nil {
		panic(err)
	}
}
func main() {
	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	r.Run()

}
