// SeatMap.tsx
import React, { useState, useEffect } from "react";
import axios from "axios";
import { Button, Typography } from "@mui/material";

interface Seat {
  id: string;
  Row: number;
  Col: number;
  IsBooked: boolean;
  FlightSectionID: string;
  FlightNumber: string;
}

interface SeatMapProps {
  flightNumber: string;
}

const SeatMap: React.FC<SeatMapProps> = ({ flightNumber }) => {
  const [seats, setSeats] = useState<Seat[]>([]);
  const [selectedSeats, setSelectedSeats] = useState<string[]>([]);

  useEffect(() => {
    const fetchSeats = async () => {
      try {
        const response = await axios.get<Seat[]>(
          `http://127.0.0.1:3000/seats/flight/${flightNumber}`
        );
        setSeats(response.data);
      } catch (error) {
        console.error("Error fetching seats:", error);
      }
    };

    fetchSeats();
  }, [flightNumber]);

  const updateSeatStatus = async (seatID: string) => {
    try {
      // Assuming you have a FlightSectionID to send in the request
      const selectedSeat = seats.find((seat) => seat.id === seatID);
      const flightSectionID = selectedSeat?.FlightSectionID;

      await axios.put(`http://127.0.0.1:3000/seats/${seatID}`, {
        IsBooked: true,
        FlightSectionID: flightSectionID,
      });
    } catch (error) {
      console.error("Error updating seat status:", error);
    }
  };

  const handleSeatClick = (seatID: string) => {
    const selectedSeat = seats.find((seat) => seat.id === seatID);
    if (selectedSeat?.IsBooked) {
      return;
    }

    const isSeatSelected = selectedSeats.includes(seatID);
    const updatedSelectedSeats = isSeatSelected
      ? selectedSeats.filter((selected) => selected !== seatID)
      : [...selectedSeats, seatID];

    setSelectedSeats(updatedSelectedSeats);
  };

  const handleSubmit = () => {
    selectedSeats.forEach((seatID) => {
      updateSeatStatus(seatID);
    });

    // Clear the selected seats after submitting
    setSelectedSeats([]);
  };
  // Extract unique rows and sort them
  const uniqueRows = Array.from(new Set(seats.map((seat) => seat.Row))).sort();
  // Extract unique columns and sort them
  const uniqueCols = Array.from(new Set(seats.map((seat) => seat.Col))).sort();

  return (
    <div style={{ textAlign: "center", marginTop: "20px" }}>
      <Typography variant="h4" align="center" mt={3} mb={3}>
        Seat map for flight {flightNumber}
      </Typography>

      <div
        style={{
          display: "flex",
          alignItems: "center",
          flexDirection: "column",
        }}
      >
        <div style={{ display: "flex" }}>
          {/* Empty space for row labels */}
          <div style={{ width: "40px" }}></div>
          {/* Column Labels */}
          {uniqueCols.map((col) => (
            <div
              key={col}
              style={{
                width: "40px",
                textAlign: "center",
                fontWeight: "bold",
                fontSize: "14px",
                margin: "2px",
              }}
            >
              {col}
            </div>
          ))}
        </div>
        {/* Seats */}
        {uniqueRows.map((row) => (
          <div key={row} style={{ display: "flex" }}>
            {/* Row Label */}
            <div
              style={{
                width: "40px",
                textAlign: "center",
                fontWeight: "bold",
                fontSize: "14px",
              }}
            >
              {row}
            </div>
            {/* Seats */}
            {uniqueCols.map((col) => {
              const seat = seats.find((s) => s.Row === row && s.Col === col);
              return (
                <div
                  key={seat?.id || `${row}-${col}`}
                  style={{
                    width: "40px",
                    height: "40px",
                    backgroundColor: seat?.IsBooked
                      ? "red"
                      : selectedSeats.includes(seat?.id || "")
                      ? "blue"
                      : "green",
                    cursor: seat?.IsBooked ? "not-allowed" : "pointer",
                    color: "white",
                    display: "flex",
                    alignItems: "center",
                    justifyContent: "center",
                    fontWeight: "bold",
                    fontSize: "14px",
                    margin: "3px",
                  }}
                  onClick={() => seat && handleSeatClick(seat.id)}
                >
                  {seat?.Row}
                </div>
              );
            })}
          </div>
        ))}
      </div>
      <Button
        variant="contained"
        color="primary"
        onClick={handleSubmit}
        disabled={selectedSeats.length === 0}
        style={{ marginTop: "20px" }}
      >
        Book selected seats
      </Button>
    </div>
  );
};

export default SeatMap;
