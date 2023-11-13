import React, { useEffect, useState } from "react";
import NavigationBar from "../components/Navigation";
import Footer from "../components/Footer";
import { DataGrid, GridColDef } from "@mui/x-data-grid";
import axios from "axios";
import {
  Box,
  Button,
} from "@mui/material";
import FlightForm from "../components/FlightForm";


const Flights: React.FC = () => {
    const [flightData, setFlightData] = useState<any[]>([]);
    const [loading, setLoading] = useState<boolean>(true);
    const [open, setOpen] = React.useState(false);
    const [errorMessage, setErrorMessage] = useState("");
    

    const columns: GridColDef[] = [
      {
        field: "id",
        headerName: "UUID",
        flex: 1,
      },
      {
        field: "flightNumber",
        headerName: "Flight Number",
        flex: 1,
      },
      {
        field: "flightSectionID",
        headerName: "Flight Section ID",
        flex: 1,
        renderCell: (params) => {          
          const sectionIds = params.row.FlightSectionID as string[];
          return (
            <ul>
              {sectionIds.map((sectionId) => (
                <li key={sectionId}>{sectionId}</li>
              ))}
            </ul>
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
      },
      {
        field: "flightTime",
        headerName: "Flight Time",
        flex: 1,
      },
      {
        field: "eta",
        headerName: "ETA",
        flex: 1,
      },
    ];


    const handleOpen = () => {
        setOpen(true);
        setErrorMessage("");
      };
    
      const handleClose = () => setOpen(false);
      const fetchFlightData = async () => {
        try {
          const flightResponse = await axios.get(
            "http://127.0.0.1:3000/flights"
          );
         
          setFlightData(flightResponse.data);
          setLoading(false);          
        } catch (error) {
          setLoading(false);
        }
      };
      useEffect(() => {
        
          fetchFlightData();
      }, []);


    
    return(
    <div>
        <NavigationBar/>
        <Box>
        <Button onClick={handleOpen}>Create Flight</Button>
      </Box>
      <FlightForm open={open} onClose={handleClose} />
      <Box sx={{ height: 400, width: "100%" }}>
        <DataGrid
          rows={flightData}
          columns={columns}
          loading={loading}
          initialState={{
            pagination: {
              paginationModel: {
                pageSize: 5,
              },
            },
          }}
          pageSizeOptions={[5]}
          checkboxSelection
          disableRowSelectionOnClick
        />
      </Box>
        <Footer/>
    </div>
)}

export default Flights;