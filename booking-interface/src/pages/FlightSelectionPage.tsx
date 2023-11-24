// FlightSelectionPage.tsx
import React, { useState, useEffect } from "react";
import axios from "axios";
import SeatMap from "../components/SeatMap"; // Import the SeatMap component
import NavigationBar from "../components/Navigation";
import Footer from "../components/Footer";

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
      <h2>Select a Flight</h2>
      <ul>
        {flights.map((flight) => (
          <li key={flight.id} onClick={() => handleFlightSelect(flight)}>
            {flight.flightNumber}
          </li>
        ))}
      </ul>

      {selectedFlight && <SeatMap flightNumber={selectedFlight.flightNumber} />}
      <Footer />
    </div>
  );
};

export default FlightSelectionPage;
