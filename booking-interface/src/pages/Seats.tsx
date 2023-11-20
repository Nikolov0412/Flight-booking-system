import React, { useEffect, useState } from "react";
import NavigationBar from "../components/Navigation";
import Footer from "../components/Footer";
import { DataGrid, GridColDef } from "@mui/x-data-grid";
import axios from "axios";

const columns: GridColDef[] = [
  {
    field: "id",
    headerName: "UUID",
    flex: 1,
  },
  {
    field: "Row",
    headerName: "Row",
    flex: 1,
  },
  {
    field: "Col",
    headerName: "Col",
    flex: 1,
  },
  {
    field: "IsBooked",
    headerName: "Booked",
    flex: 1,
  },
  {
    field: "FlightSectionId",
    headerName: "Flight Section ID",
    flex: 1,
  },
  {
    field: "FlightNumber",
    headerName: "Flight Number",
    flex: 1,
  },
];
const Seats: React.FC = () => {
  const [data, setData] = useState<any[]>([]);
  const [loading, setLoading] = useState<boolean>(true);

  useEffect(() => {
    const fetchData = async () => {
      try {
        const response = await axios.get("http://127.0.0.1:3000/seats");
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
      <DataGrid
        sx={{ minHeight: "400px " }}
        rows={data}
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
      <Footer />
    </div>
  );
};

export default Seats;
