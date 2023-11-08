import React from "react";
import AppBar from "@mui/material/AppBar";
import Box from "@mui/material/Box";
import Toolbar from "@mui/material/Toolbar";
import Typography from "@mui/material/Typography";
import Button from "@mui/material/Button";
import Link from "@mui/material/Link";
const NavigationBar: React.FC = () => {
  return (
    <Box sx={{ flexGrow: 1 }}>
      <AppBar position="static">
        <Toolbar>
          <Typography variant="h6" component="div" sx={{ flexGrow: 1 }}>
            Booking Systen
          </Typography>
          <Button color="inherit" LinkComponent={Link} href="/">
            Home
          </Button>
          <Button color="inherit" LinkComponent={Link} href="/airlines">
            Airlines
          </Button>
          <Button color="inherit" LinkComponent={Link} href="/airports">
            Airports
          </Button>
          <Button color="inherit" LinkComponent={Link} href="/flightsections">
            Flight Sections
          </Button>
          <Button color="inherit" LinkComponent={Link} href="/flights">
            Flights
          </Button>
          <Button color="inherit" LinkComponent={Link} href="/seats">
            Seats
          </Button>
        </Toolbar>
      </AppBar>
    </Box>
  );
};

export default NavigationBar;
