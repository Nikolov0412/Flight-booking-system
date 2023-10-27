package main

import (
	"context"
	"log"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	ginadapter "github.com/awslabs/aws-lambda-go-api-proxy/gin"
	"github.com/gin-gonic/gin"
)

var svc *dynamodb.DynamoDB
var ginLambda *ginadapter.GinLambdaV2

type Response struct {
	message string `json:"message"`
}

func main() {
	var err error
	svc, err = initDynamoDB()
	log.Printf(svc.Endpoint)
	log.Printf(svc.ServiceName)
	if err != nil {
		log.Printf(err.Error())
		panic(err)
	}
	lambda.Start(Handler)
}
func Handler(ctx context.Context, req events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	return ginLambda.ProxyWithContext(ctx, req)
}

func init() {
	log.Printf("Gin cold start")
	r := gin.Default()
	// Define a route for creating airlines
	r.POST("/airlines", func(c *gin.Context) {
		var airline Airline

		// Bind the request body to the Airline struct
		if err := c.ShouldBindJSON(&airline); err != nil {
			c.AbortWithError(http.StatusBadRequest, err)
			return
		}

		// Call the CreateAirline function to create the airline in DynamoDB
		if err := CreateAirline(airline, svc); err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		c.JSON(http.StatusCreated, Response{message: "Airline created successfully"})
	})

	r.GET("/airlines", func(c *gin.Context) {
		airlines, err := GetAllAirlines(svc)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		c.JSON(http.StatusOK, airlines)
	})
	r.GET("/airlines/:id", func(c *gin.Context) {
		airlineID := c.Param("id")

		airline, err := GetAirlineByID(airlineID, svc)
		if err != nil {
			c.AbortWithError(http.StatusNotFound, err)
			return
		}

		c.JSON(http.StatusOK, airline)
	})

	r.POST("/airports", func(c *gin.Context) {
		var airport Airport

		if err := c.ShouldBindJSON(&airport); err != nil {
			c.AbortWithError(http.StatusBadRequest, err)
			return
		}

		if err := CreateAirport(airport, svc); err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		c.JSON(http.StatusCreated, Response{message: "Airport created successfully"})
	})

	r.GET("/airports", func(c *gin.Context) {
		airports, err := GetAllAirports(svc)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		c.JSON(http.StatusOK, airports)
	})

	r.GET("/airports/:id", func(c *gin.Context) {
		// Get the ID parameter from the URL
		airportID := c.Param("id")

		// Call the GetAirportByID function to retrieve the airport by ID
		airport, err := GetAirportByID(airportID, svc)
		if err != nil {
			c.AbortWithError(http.StatusNotFound, err)
			return
		}

		// Return the airport as JSON response
		c.JSON(http.StatusOK, airport)
	})

	r.POST("/seats", func(c *gin.Context) {
		var seat Seat

		// Bind the request body to the Seat struct
		if err := c.ShouldBindJSON(&seat); err != nil {
			c.AbortWithError(http.StatusBadRequest, err)
			return
		}

		// Call the CreateSeat function to create the seat in DynamoDB
		if err := CreateSeat(seat, svc); err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		c.JSON(http.StatusCreated, Response{message: "Seat created successfully"})
	})
	r.GET("/seats", func(c *gin.Context) {
		seats, err := GetAllSeats(svc)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		c.JSON(http.StatusOK, seats)
	})
	r.GET("/seats/flight/:flightNumber", func(c *gin.Context) {
		flightNumber := c.Param("flightNumber")

		seats, err := GetSeatsByFlightNumber(flightNumber, svc)
		if err != nil {
			c.AbortWithError(http.StatusNotFound, err)
			return
		}

		c.JSON(http.StatusOK, seats)
	})
	r.GET("/seats/flightsection/:flightSectionID", func(c *gin.Context) {
		flightSectionID := c.Param("flightSectionID")

		seats, err := GetSeatsByFlightSectionID(flightSectionID, svc)
		if err != nil {
			c.AbortWithError(http.StatusNotFound, err)
			return
		}

		c.JSON(http.StatusOK, seats)
	})
	r.GET("/seats/:id", func(c *gin.Context) {
		seatID := c.Param("id")

		seat, err := GetSeatByID(seatID, svc)
		if err != nil {
			c.AbortWithError(http.StatusNotFound, err)
			return
		}

		c.JSON(http.StatusOK, seat)
	})

	r.PUT("/seats/:id", func(c *gin.Context) {
		seatID := c.Param("id")

		var updateData struct {
			IsBooked bool `json:"IsBooked"`
		}

		// Bind the request body to the updateData struct
		if err := c.ShouldBindJSON(&updateData); err != nil {
			c.AbortWithError(http.StatusBadRequest, err)
			return
		}

		err := UpdateSeatIsBooked(seatID, updateData.IsBooked, svc)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		c.JSON(http.StatusOK, Response{message: "Seat updated successfully"})
	})

	r.POST("/flightsections", func(c *gin.Context) {
		var flightSection FlightSection

		if err := c.ShouldBindJSON(&flightSection); err != nil {
			c.AbortWithError(http.StatusBadRequest, err)
			return
		}

		if err := CreateFlightSection(flightSection, svc); err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)

			return
		}

		c.JSON(http.StatusCreated, Response{message: "Flight section created successfully"})
	})

	r.GET("/flightsections", func(c *gin.Context) {
		flightSections, err := GetAllFlightSections(svc)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		c.JSON(http.StatusOK, flightSections)
	})

	r.GET("/flightsections/:id", func(c *gin.Context) {
		// Retrieve the flight section ID from the URL parameters
		sectionID := c.Param("id")

		// Call the GetFlightSectionByID function to fetch the flight section
		flightSection, err := GetFlightSectionByID(sectionID, svc)
		if err != nil {
			c.AbortWithError(http.StatusNotFound, err)
			return
		}

		c.JSON(http.StatusOK, flightSection)
	})

	r.POST("/flights", func(c *gin.Context) {
		var flight Flight

		// Bind the request body to the Flight struct
		if err := c.ShouldBindJSON(&flight); err != nil {
			c.AbortWithError(http.StatusBadRequest, err)
			return
		}

		// Call the CreateFlight function to create the flight in DynamoDB
		if err := CreateFlight(flight, svc); err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		c.JSON(http.StatusCreated, Response{message: "Flight created successfully"})
	})
	r.GET("/flights", func(c *gin.Context) {
		flights, err := GetAllFlights(svc)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}
		c.JSON(http.StatusOK, flights)
	})
	r.GET("/flights/origin/:airport", func(c *gin.Context) {
		originAirport := c.Param("airport")
		flights, err := GetFlightsByOriginAirport(originAirport, svc)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}
		c.JSON(http.StatusOK, flights)
	})

	r.GET("/flights/destination/:airport", func(c *gin.Context) {
		destinationAirport := c.Param("airport")
		flights, err := GetFlightsByDestinationAirport(destinationAirport, svc)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}
		c.JSON(http.StatusOK, flights)
	})

	ginLambda = ginadapter.NewV2(r)

}
