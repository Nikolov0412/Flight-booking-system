import React, { useEffect, useState } from "react";
import NavigationBar from "../components/Navigation";
import Footer from "../components/Footer";
import { DataGrid, GridColDef } from "@mui/x-data-grid";
import axios from "axios";
import {
  Box,
  Button,
  FormControl,
  MenuItem,
  Modal,
  Select,
  SelectChangeEvent,
  TextField,
  Typography,
} from "@mui/material";
import PublishIcon from "@mui/icons-material/Publish";

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
    field: "seatClass",
    headerName: "SeatClass",
    flex: 1,
  },
  {
    field: "numRows",
    headerName: "Number rows",
    flex: 1,
  },
  {
    field: "numCols",
    headerName: "Number Cols",
    flex: 1,
  },
];

const FlightSections: React.FC = () => {
  const [data, setData] = useState<any[]>([]);
  const [loading, setLoading] = useState<boolean>(true);
  const [open, setOpen] = React.useState(false);
  const [errorMessage, setErrorMessage] = useState("");
  const [formData, setFormData] = useState({
    seatClass: "",
    numRows: 0,
    numCols: 0,
  });

  const handleChange = (event: SelectChangeEvent) => {
    setFormData((prevState) => {
      return {
        ...prevState,
        seatClass: event.target.value,
      };
    });
  };

  const handleOpen = () => {
    setOpen(true);
    setErrorMessage("");
  };

  const handleClose = () => setOpen(false);

  const handleSubmit = async (e: any) => {
    e.preventDefault();
    try {
      await axios.post("http://127.0.0.1:3000/flightsections", formData);
      handleClose();
    } catch (error: any) {
      if (error.response && error.response.status === 500) {
        setErrorMessage(
          "There was an error trying to create Flight Section please try again."
        );
      }
    }
  };

  useEffect(() => {
    const fetchData = async () => {
      try {
        const response = await axios.get(
          "http://127.0.0.1:3000/flightsections"
        );
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
          <FormControl fullWidth>
            <Typography id="modal-modal-title" variant="h6" component="h2">
              Flight Section creation form
            </Typography>
            <Typography variant="body1">Seat Class:</Typography>
            <Select
              labelId="demo-simple-select-disabled-label"
              id="demo-simple-select-disabled"
              value={formData.seatClass}
              label="Seat Class"
              onChange={handleChange}
            >
              <MenuItem value={"Economy"}>Economy</MenuItem>
              <MenuItem value={"Business"}>Business</MenuItem>
              <MenuItem value={"First Class"}>First Class</MenuItem>
            </Select>
            <Typography variant="body1">Number of columns</Typography>
            <TextField
              required
              id="outlined-required"
              label="Required"
              value={formData.numCols}
              onChange={(e) =>
                setFormData((prevState) => {
                  return {
                    ...prevState,
                    numCols: Number(e.target.value),
                  };
                })
              }
            />
            <Typography variant="body1">Number of Rows</Typography>

            <TextField
              required
              id="outlined-required"
              label="Required"
              value={formData.numRows}
              onChange={(e) =>
                setFormData((prevState) => {
                  return {
                    ...prevState,
                    numRows: Number(e.target.value),
                  };
                })
              }
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
          </FormControl>
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

export default FlightSections;
