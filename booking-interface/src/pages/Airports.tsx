import React, { useEffect, useState } from "react";
import NavigationBar from "../components/Navigation";
import Footer from "../components/Footer";
import { Box, Button, Modal, Typography, TextField } from "@mui/material";
import PublishIcon from "@mui/icons-material/Publish";
import { DataGrid, GridColDef } from "@mui/x-data-grid";
import axios from "axios";

/* Styles and grid definition */
const modalStyle = {
  position: "absolute" as "absolute",
  top: "50%",
  left: "50%",
  transform: "translate(-50%, -50%)",
  width: 400,
  bgcolor: "background.paper",
  border: "2px solid #000",
  boxShadow: 24,
  p: 4,
};
const columns: GridColDef[] = [
  {
    field: "id",
    headerName: "UUID",
    flex: 1,
  },
  {
    field: "code",
    headerName: "Airport code",
    flex: 1,
  },
];
const Airports: React.FC = () => {
  /* States */
  const [data, setData] = useState<any[]>([]);
  const [loading, setLoading] = useState<boolean>(true);
  const [open, setOpen] = React.useState(false);
  const [airportCode, setAirportCode] = useState("");
  const [errorMessage, setErrorMessage] = useState("");

  /* Handlers and hooks */

  const handleOpen = () => {
    setOpen(true);
    setErrorMessage("");
  };

  const handleClose = () => setOpen(false);

  const handleSubmit = async (e: any) => {
    e.preventDefault();
    try {
      await axios.post("http://127.0.0.1:3000/airports", { code: airportCode });
      handleClose();
    } catch (error: any) {
      if (error.response && error.response.status === 500) {
        setErrorMessage(
          "The following Airport code is already used or it's more than 3 characters in length."
        );
      }
    }
  };
  useEffect(() => {
    const fetchData = async () => {
      try {
        const response = await axios.get("http://127.0.0.1:3000/airports");
        setData(response.data);
        setLoading(false);
      } catch (error) {
        setLoading(false);
      }
    };

    fetchData();
  }, []);

  return (
    <div>
      <NavigationBar />
      <Box>
        <Button onClick={handleOpen}>Create Airport</Button>
      </Box>
      <Modal
        open={open}
        onClose={handleClose}
        aria-labelledby="modal-modal-title"
        aria-describedby="modal-modal-description"
      >
        <Box sx={modalStyle} component="form" onSubmit={handleSubmit}>
          <Typography id="modal-modal-title" variant="h6" component="h2">
            Airport creation form
          </Typography>
          <Typography variant="body1">Airport Code:</Typography>
          <TextField
            required
            id="outlined-required"
            label="Required"
            value={airportCode}
            onChange={(e) => setAirportCode(e.target.value)}
          />
          <br />
          {errorMessage && (
            <Typography variant="body2" color="error">
              {errorMessage}
            </Typography>
          )}
          <br />
          <Button variant="contained" type="submit" endIcon={<PublishIcon />}>
            Submit
          </Button>
        </Box>
      </Modal>
      <Box sx={{ height: 400, width: "100%" }}>
        <DataGrid
          rows={data}
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
      <Footer />
    </div>
  );
};

export default Airports;
