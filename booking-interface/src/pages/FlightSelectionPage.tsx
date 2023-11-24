import React, { useState, useEffect } from "react";
import axios from "axios";
import SeatMap from "../components/SeatMap";
import NavigationBar from "../components/Navigation";
import Footer from "../components/Footer";
import {
  Typography,
  List,
  ListItem,
  ListItemText,
  Container,
} from "@mui/material";

interface Flight {
  id: string;
  flightNumber: string;
  FlightSectionID: string[];
  originAirport: string;
  destinationAirport: string;
  departureDate: any;
  flightTime: any;
  eta: string;
}

const FlightSelectionPage: React.FC = () => {
  const [flights, setFlights] = useState<Flight[]>([]);
  const [selectedFlight, setSelectedFlight] = useState<Flight | null>(null);

  useEffect(() => {
    const fetchFlights = async () => {
      try {
        const response = await axios.get<Flight[]>(
          "http://127.0.0.1:3000/flights"
        );
        setFlights(response.data);
      } catch (error) {
        console.error("Error fetching flights:", error);
      }
    };

    fetchFlights();
  }, []);

  const handleFlightSelect = (flight: Flight) => {
    setSelectedFlight(flight);
  };

  return (
    <div>
      <NavigationBar />
      <Container>
        <Typography variant="h4" align="center" mt={3} mb={3}>
          Select a Flight
        </Typography>
        <List>
          {flights.map((flight) => (
            <ListItem
              key={flight.id}
              button
              onClick={() => handleFlightSelect(flight)}
            >
              <ListItemText primary={flight.flightNumber} />
            </ListItem>
          ))}
        </List>

        {selectedFlight && (
          <SeatMap flightNumber={selectedFlight.flightNumber} />
        )}
      </Container>
      <Footer />
    </div>
  );
};

export default FlightSelectionPage;
