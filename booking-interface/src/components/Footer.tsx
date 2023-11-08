import React from "react";
import { Paper, Typography } from "@mui/material";

const Footer: React.FC = () => {
  return (
    <Paper elevation={3} style={{ padding: "20px", textAlign: "center" }}>
      <Typography variant="body2" color="textSecondary">
        &copy; {new Date().getFullYear()} Flight Booking System
      </Typography>
    </Paper>
  );
};

export default Footer;
