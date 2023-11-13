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
} from "@mui/material";
import axios from "axios";
import {DateTimePicker,TimePicker,LocalizationProvider} from "@mui/x-date-pickers"
import { AdapterMoment } from "@mui/x-date-pickers/AdapterMoment";
import moment, { Duration, Moment } from "moment";

interface FlightFormProps {
  open: boolean;
  onClose: () => void;
}

interface FlightData {
    flightNumber: string;
    flightSectionID: string[];
    originAirport: string;
    destinationAirport: string;
    departureDate: Moment;
    flightTimeDate:Moment;
    flightTime: Duration| null;
  }

interface AirportData{
    id:string,
    code:string
}

  const modalStyles={
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


const FlightForm: React.FC<FlightFormProps> = ({ open, onClose }) => {
  const [flightData, setFlightData] = useState<FlightData>({
    flightNumber: "",
    flightSectionID: [],
    originAirport: "",
    destinationAirport: "",
    departureDate: moment(),
    flightTimeDate:moment(),
    flightTime: moment.duration(0),
  });
  const [airports, setAirports] = useState<AirportData[]>([]);
  const [flightSections, setFlightSections] = useState<any[]>([]);
  const [loadingAirports, setLoadingAirports] = useState<boolean>(true);

  useEffect(() => {
    const fetchAirports = async () => {
      try {
        const response = await axios.get("http://127.0.0.1:3000/airports");
        setAirports(response.data);
        setLoadingAirports(false);
      } catch (error) {
        console.error("Error fetching airports:", error);
        setLoadingAirports(false);
      }
    };
    const fetchFlightSections = async () => {
        try {
          const response = await axios.get('http://127.0.0.1:3000/flightsections');
          setFlightSections(response.data);
        } catch (error) {
          console.error('Error fetching flight sections:', error);
        }
      };
    

    fetchAirports();
    fetchFlightSections();

  }, []);
 
  const handleDateChange = (date: moment.Moment | null) => {
    setFlightData((prevData) => ({
      ...prevData,
      departureDate: date || moment(),
    }));
  };

  const handleFlightTimeChange = (newTime: Moment|null ) => {
    setFlightData((prevData) => ({
        ...prevData,
        flightTime: moment.duration(newTime?.format("hh:mm")) || moment.duration(1),
      }));        
    };      

  const handleSelectChange = (event: SelectChangeEvent) => {
    setFlightData({
        ...flightData,
        [event.target.name]: event.target.value,
      });
  };

  const handleInputChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    setFlightData({
      ...flightData,
      [e.target.name]: e.target.value,
    });
  };
  function millisecondsToNanoseconds(milliseconds: any): number {
    const nanosecondsPerMillisecond = 1e6; // 1 millisecond = 1,000,000 nanoseconds
    return milliseconds * nanosecondsPerMillisecond;
  }
  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    try {
        const flightTimeFormated=millisecondsToNanoseconds(flightData.flightTime!.asMilliseconds)
      await axios.post("http://127.0.0.1:3000/flights", {
        flightNumber:flightData.flightNumber,
        flightSectionID:flightData.flightSectionID,
        originAirport:flightData.originAirport,
        destinationAirport:flightData.destinationAirport,
        departureDate:flightData.departureDate,
        flightTime: flightTimeFormated});
      onClose(); // Close the modal if the request is successful
    } catch (error) {
      console.error("Error creating flight:", error);
      // Handle error, show error message, etc.
    }
  };

  return (
    <Modal open={open} onClose={onClose}> 
      <Box sx={modalStyles}>
        <Typography variant="h6" gutterBottom>
          Create Flight
        </Typography>
        <form onSubmit={handleSubmit}>
          <TextField
            label="Flight Number"
            name="flightNumber"
            value={flightData.flightNumber}
            onChange={handleInputChange}
            required
          />
       <br />
          <br />
          <FormControl fullWidth>
            <InputLabel htmlFor="destinationAirport">Origin Airport</InputLabel>
          <Select
              label="Origin Airport"
              name="originAirport"
              value={flightData.originAirport}
              onChange={handleSelectChange}
              required
            >
              {loadingAirports ? (
                <MenuItem value="" disabled>
                  Loading airports...
                </MenuItem>
              ) : (
                airports.map((airport) => (
                  <MenuItem key={airport.id} value={airport.code}>
                    {airport.code}
                  </MenuItem>
                ))
              )}
            </Select>
            </FormControl>
            <br />
          <br />
          <FormControl fullWidth>
            <InputLabel htmlFor="destinationAirport">Destination Airport</InputLabel>
            <Select
              label="Destination Airport"
              name="destinationAirport"
              value={flightData.destinationAirport}
              onChange={handleSelectChange}
              required
            >
              {loadingAirports ? (
                <MenuItem value="" disabled>
                  Loading airports...
                </MenuItem>
              ) : (
                airports.map((airport) => (
                  <MenuItem key={airport.id} value={airport.code}>
                    {airport.code}
                  </MenuItem>
                ))
              )}
            </Select>
          </FormControl>
          <br />
          <br />
          <FormControl fullWidth>
          <InputLabel htmlFor="destinationAirport">Flight Section ID's</InputLabel>
          <Select
            multiple
            value={flightData.flightSectionID}
            onChange={(e) => setFlightData((prevData) => ({ ...prevData, flightSectionID: e.target.value as string[] }))}
          >
            {/* Render options based on the fetched flight sections */}
            {flightSections.map((flightSection) => (
              <MenuItem key={flightSection.seatClass} value={flightSection.id}>
                {flightSection.seatClass}
              </MenuItem>
            ))}
          </Select>
          </FormControl>
          <br/>
          <br/>
          <LocalizationProvider dateAdapter={AdapterMoment}>

          <DateTimePicker
            label="Departure Date"
            value={flightData.departureDate}
            onChange={handleDateChange}
          />
          </LocalizationProvider>

    <br />
          <br />
          <LocalizationProvider dateAdapter={AdapterMoment}>
          <TimePicker
            label="Flight Time"
            value={flightData.flightTimeDate}
            onChange={handleFlightTimeChange}
            views={['hours', 'minutes']} format="hh:mm" 
          />
              </LocalizationProvider>

          <br />
          <br />
          <Button variant="contained" type="submit">
            Submit
          </Button>
        </form>
      </Box>
    </Modal>
  );
};

export default FlightForm;
