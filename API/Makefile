build:
	@echo "Building your Go application..."
	set GOOS=linux
	go build -o bin/main main.go airline.go airport.go database.go flight.go flightsection.go utils.go seat.go
package:
	@echo "Packaging your application..."
	mkdir bin
	set GOOS=linux
	go build -ldflags="-s -w" -o bin/main main.go airline.go airport.go database.go flight.go flightsection.go utils.go seat.go
	tar -czvf bin/api.zip -C bin api
	del /Q bin\api
clean:
	@echo "Cleaning up..."
	del /Q bin\*
