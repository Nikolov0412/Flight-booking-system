package main

import (
	"net/http"

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
	// Define a route for creating airlines
	r.POST("/airlines", func(c *gin.Context) {
		var airline Airline

		// Bind the request body to the Airline struct
		if err := c.ShouldBindJSON(&airline); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Call the CreateAirline function to create the airline in DynamoDB
		if err := CreateAirline(airline, svc); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusCreated, gin.H{"message": "Airline created successfully"})
	})

	r.GET("/airlines", func(c *gin.Context) {
		airlines, err := GetAllAirlines(svc)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, airlines)
	})

	r.POST("/airports", func(c *gin.Context) {
		var airport Airport

		if err := c.ShouldBindJSON(&airport); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if err := CreateAirport(airport, svc); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusCreated, gin.H{"message": "Airport created successfully"})
	})

	r.GET("/airports", func(c *gin.Context) {
		airports, err := GetAllAirports(svc)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, airports)
	})

	r.Run()

}
