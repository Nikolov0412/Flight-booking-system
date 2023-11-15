import React, { useEffect, useState } from "react";
import NavigationBar from "../components/Navigation";
import Footer from "../components/Footer";
import { DataGrid, GridColDef } from "@mui/x-data-grid";
import axios from "axios";
import { Box, Button } from "@mui/material";
import FlightForm from "../components/FlightForm";
import moment from "moment";
import "moment-timezone";

const Flights: React.FC = () => {
  const [flightData, setFlightData] = useState<any[]>([]);
  const [loading, setLoading] = useState<boolean>(true);
  const [open, setOpen] = React.useState(false);
  const [errorMessage, setErrorMessage] = useState("");

  const formatDuration = (durationInMilliseconds: number) => {
    const durationInSeconds = durationInMilliseconds / 1000;
    const hours = Math.floor(durationInSeconds / 3600);
    const minutes = Math.floor((durationInSeconds % 3600) / 60);

    const formattedHours = hours < 10 ? `0${hours}` : hours;
    const formattedMinutes = minutes < 10 ? `0${minutes}` : minutes;

    return `${formattedHours}:${formattedMinutes}`;
  };

  const columns: GridColDef[] = [
    {
      field: "id",
      headerName: "UUID",
      flex: 2,
    },
    {
      field: "flightNumber",
      headerName: "Flight Number",
      flex: 1,
    },
    {
      field: "flightSectionID",
      headerName: "Flight Section ID",
      flex: 2,
      renderCell: (params) => {
        const sectionIds = params.row.FlightSectionID as string[];
        return (
          <div>
            <ul>
              {sectionIds.map((sectionId) => (
                <li key={sectionId}>{sectionId}</li>
              ))}
            </ul>
          </div>
        );
      },
    },
    {
      field: "originAirport",
      headerName: "Origin Airport",
      flex: 1,
    },
    {
      field: "destinationAirport",
      headerName: "Destination Airport",
      flex: 1,
    },
    {
      field: "departureDate",
      headerName: "Departure Date",
      flex: 1,
      valueFormatter: (params) => moment(params.value).format("DD/MM/YY HH:mm"), // Format the date
    },
    {
      field: "flightTime",
      headerName: "Flight Time",
      flex: 1,
      valueFormatter: (params) => formatDuration(parseInt(params.value, 10)),
    },
    {
      field: "eta",
      headerName: "ETA",
      flex: 1,
      renderCell: (params) => {
        const utcTimeString = params.value; // Assuming params.value is a string in HH:mm format

        // Convert the ETA from UTC to the user's local time using Moment.js
        const userTimeZone = moment.tz.guess(); // Automatically detects the user's timezone
        const userLocalTime = moment
          .utc(`2000-01-01T${utcTimeString}:00Z`)
          .tz(userTimeZone)
          .format("HH:mm");

        return <span>{userLocalTime}</span>;
      },
    },
  ];

  const handleOpen = () => {
    setOpen(true);
    setErrorMessage("");
  };

  const handleClose = () => setOpen(false);
  const fetchFlightData = async () => {
    try {
      const flightResponse = await axios.get("http://127.0.0.1:3000/flights");

      setFlightData(flightResponse.data);
      setLoading(false);
    } catch (error) {
      setLoading(false);
    }
  };
  useEffect(() => {
    fetchFlightData();
  }, []);

  return (
    <div>
      <NavigationBar />
      <Box>
        <Button onClick={handleOpen}>Create Flight</Button>
      </Box>
      <FlightForm open={open} onClose={handleClose} />
      <Box sx={{ height: "auto", width: "100%" }}>
        <DataGrid
          sx={{ minHeight: "400px " }}
          rows={flightData}
          columns={columns}
          loading={loading}
          getRowHeight={() => "auto"}
          initialState={{
            pagination: {
              paginationModel: {
                pageSize: 25,
              },
            },
          }}
          pageSizeOptions={[25]}
          checkboxSelection
          disableRowSelectionOnClick
        />
      </Box>
      <Footer />
    </div>
  );
};

export default Flights;
