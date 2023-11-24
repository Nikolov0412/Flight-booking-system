// SeatMap.tsx
import React, { useState, useEffect } from 'react';
import axios from 'axios';

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
        console.error('Error fetching seats:', error);
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
      console.error('Error updating seat status:', error);
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

  return (
    <div style={{ textAlign: 'center' }}>
      <h2>Seat Map for Flight {flightNumber}</h2>
      <div style={{ display: 'flex', flexWrap: 'wrap', justifyContent: 'center' }}>
        {seats.map((seat) => (
          <div
            key={seat.id}
            style={{
              width: '40px',
              height: '40px',
              margin: '5px',
              backgroundColor: seat.IsBooked
                ? 'red'
                : selectedSeats.includes(seat.id)
                ? 'blue'
                : 'green',
              cursor: seat.IsBooked ? 'not-allowed' : 'pointer',
              color: 'white',
              display: 'flex',
              alignItems: 'center',
              justifyContent: 'center',
              fontWeight: 'bold',
              fontSize: '14px',
            }}
            onClick={() => handleSeatClick(seat.id)}
          >
            {seat.Row}-{seat.Col}
          </div>
        ))}
      </div>
      <button onClick={handleSubmit} disabled={selectedSeats.length === 0}>
        Submit
      </button>
    </div>
  );
};

export default SeatMap;
