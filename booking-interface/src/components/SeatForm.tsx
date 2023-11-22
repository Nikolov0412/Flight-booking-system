import React, { useEffect, useState } from "react";
import {
  Button,
  TextField,
  Modal,
  Box,
  Typography,
  Select,
  MenuItem,
  FormControl,
  InputLabel,
  SelectChangeEvent,
  Switch,
} from "@mui/material";
import axios from "axios";

interface SeatFormProps {
  open: boolean;
  onClose: () => void;
}

interface SeatData {
  Row: number;
  Col: number;
  isBooked: boolean;
  FlightSectionId: string;
  FlightNumber: string;
}

interface FlightData {
  id: string;
  flightNumber: string;
  FlightSectionID: string[];
  originAirport: string;
  destinationAirport: string;
  departureDate: any;
  flightTime: any;
  eta: string;
}
const modalStyles = {
  position: "absolute",
  top: "50%",
  left: "50%",
  transform: "translate(-50%, -50%)",
  width: 400,
  bgcolor: "background.paper",
  border: "2px solid #000",
  boxShadow: 24,
  p: 4,
};

const SeatForm: React.FC<SeatFormProps> = ({ open, onClose }) => {
  const [seatData, setSeatData] = useState<SeatData>({
    Row: 0,
    Col: 0,
    isBooked: false,
    FlightSectionId: "",
    FlightNumber: "",
  });
  const [flightsArr, setFlightsArr] = useState<FlightData[]>([]);
  const [loadingFlights, setLoadingFlights] = useState<boolean>(true);
  const [flightSections, setFlightSections] = useState<any[]>([]);
  const [filteredFlightSections, setFilteredFlightSections] = useState<any[]>(
    []
  );

  useEffect(() => {
    const fetchFlights = async () => {
      try {
        const response = await axios.get("http://127.0.0.1:3000/flights");
        setFlightsArr(response.data);
        setLoadingFlights(false);
      } catch (error) {
        console.error("Error fetching flights:", error);
        setLoadingFlights(false);
      }
    };

    const fetchFlightSections = async () => {
      try {
        const response = await axios.get(
          "http://127.0.0.1:3000/flightsections"
        );
        setFlightSections(response.data);
      } catch (error) {
        console.error("Error fetching flight sections:", error);
      }
    };

    fetchFlightSections();
    fetchFlights();
  }, []);

  const handleFlightChange = (event: SelectChangeEvent) => {
    const selectedFlight = flightsArr.find(
      (flight) => flight.flightNumber === event.target.value
    );

    if (selectedFlight) {
      const flightSectionIDs = selectedFlight.FlightSectionID || [];
      const filteredSections = flightSections.filter((section) =>
        flightSectionIDs.includes(section.id)
      );
      setFilteredFlightSections(filteredSections);
      setSeatData((prevData) => ({
        ...prevData,
        FlightNumber: event.target.value as string,
      }));
    }
  };

  const handleInputChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    setSeatData({
      ...seatData,
      [e.target.name]: parseInt(e.target.value, 10) || 0,
    });
  };
  const handleSwitchChange = () => {
    setSeatData({
      ...seatData,
      isBooked: !seatData.isBooked,
    });
  };
  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    try {
      await axios.post("http://127.0.0.1:3000/seats", {
        Row: seatData.Row,
        Col: seatData.Col,
        isBooked: seatData.isBooked,
        FlightSectionID: seatData.FlightSectionId,
        FlightNumber: seatData.FlightNumber,
      });
      onClose(); // Close the modal if the request is successful
    } catch (error) {
      console.error("Error creating Seat:", error);
      // Handle error, show error message, etc.
    }
  };

  return (
    <Modal open={open} onClose={onClose}>
      <Box sx={modalStyles}>
        <Typography variant="h6" gutterBottom>
          Create Seats
        </Typography>
        <form onSubmit={handleSubmit}>
          <TextField
            label="Row"
            name="Row"
            value={seatData.Row}
            onChange={handleInputChange}
            required
          />
          <br />
          <br />
          <TextField
            label="Col"
            name="Col"
            value={seatData.Col}
            onChange={handleInputChange}
            required
          />
          <br />
          <br />
          <FormControl fullWidth>
            <InputLabel htmlFor="FlightNumber">Flight Number</InputLabel>
            <Select
              label="Flight Number"
              name="flightNumber"
              value={seatData.FlightNumber}
              onChange={handleFlightChange}
              required
            >
              {loadingFlights ? (
                <MenuItem value="" disabled>
                  Loading airports...
                </MenuItem>
              ) : (
                flightsArr.map((flight: FlightData) => (
                  <MenuItem key={flight.id} value={flight.flightNumber}>
                    {flight.flightNumber}
                  </MenuItem>
                ))
              )}
            </Select>
          </FormControl>
          <br />
          <br />
          <FormControl fullWidth>
            <InputLabel htmlFor="flightSectionID">
              Flight Section ID's
            </InputLabel>
            <Select
              value={seatData.FlightSectionId}
              onChange={(e) =>
                setSeatData((prevData) => ({
                  ...prevData,
                  FlightSectionId: e.target.value as string,
                }))
              }
              required
            >
              {filteredFlightSections.map((flightSection) => (
                <MenuItem key={flightSection.id} value={flightSection.id}>
                  {flightSection.seatClass}
                </MenuItem>
              ))}
            </Select>
          </FormControl>
          <br />
          <br />
          <FormControl>
            <Typography variant="body2">Book Status:</Typography>
            <Switch checked={seatData.isBooked} onChange={handleSwitchChange} />
          </FormControl>
          <br />
          <br />
          <Button type="submit" variant="contained" color="primary">
            Create Seat
          </Button>
        </form>
      </Box>
    </Modal>
  );
};
export default SeatForm;
