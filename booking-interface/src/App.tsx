import { BrowserRouter, Routes, Route } from "react-router-dom";
import Home from "./pages/Home";
import Airlines from "./pages/Airlines";
import Airports from "./pages/Airports";
import FlightSections from "./pages/FlightSections";
import Flights from "./pages/Flights";
import Seats from "./pages/Seats";
import FlightSelectionPage from "./pages/FlightSelectionPage";
function App() {
  return (
    <div className="App">
      <BrowserRouter>
        <Routes>
          <Route path="/" element={<Home />} />
          <Route path="/airlines" element={<Airlines />} />
          <Route path="/airports" element={<Airports />} />
          <Route path="/flightsections" element={<FlightSections />} />
          <Route path="/flights" element={<Flights />} />
          <Route path="/seats" element={<Seats />} />
          <Route path="/book" element={<FlightSelectionPage />} />
        </Routes>
      </BrowserRouter>
    </div>
  );
}

export default App;
