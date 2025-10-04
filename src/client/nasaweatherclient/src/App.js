import logo from "./logo.svg";
import "./App.css";
import { useEffect, useState } from "react";
import "leaflet/dist/leaflet.css";
import { MapContainer, TileLayer, Marker, Popup } from "react-leaflet";
import L from "leaflet";

function App() {
  const [position, setPosition] = useState(null);
  const [airData, setAirData] = useState(null);
  const [error, setError] = useState("");

  // Fix Leaflet marker icon issue in React
  delete L.Icon.Default.prototype._getIconUrl;
  L.Icon.Default.mergeOptions({
    iconRetinaUrl:
      "https://unpkg.com/leaflet@1.9.4/dist/images/marker-icon-2x.png",
    iconUrl: "https://unpkg.com/leaflet@1.9.4/dist/images/marker-icon.png",
    shadowUrl: "https://unpkg.com/leaflet@1.9.4/dist/images/marker-shadow.png",
  });

  useEffect(() => {
    if (!navigator.geolocation) {
      setError("Geolocation is not supported by your browser");
      return;
    }

    navigator.geolocation.getCurrentPosition(
      (pos) => {
        const { latitude, longitude } = pos.coords;
        setPosition([latitude, longitude]);

        (async () => {
          const fetchAirData = async () => {
            try {
              const res = await fetch(`/api/data?lat=${latitude}&lon=${longitude}`);
              const data = await res.json();
              setAirData(data);
            } catch (err) {
              console.error("Failed to fetch air quality data:", err);
              setError("Failed to fetch air quality data");
              setAirData({
                aqi: 112,
                pm2_5: 35.2,
                pm10: 56.4,
                status: "Moderate (Dummy)",
              });
            }
          };
          await fetchAirData();
          const interval = setInterval(fetchAirData, 5 * 60 * 1000);
          return () => clearInterval(interval);
        })();
      },
      () => {
        setError("Failed to get your location");
      }
    );
  }, []);

  const getColor = (aqi) => {
    if (aqi <= 50) return "green";
    if (aqi <= 100) return "yellow";
    if (aqi <= 150) return "orange";
    if (aqi <= 200) return "red";
    return "purple";
  };

  return (
    <div>
      <h1>Air Quality</h1>
      {error && <p style={{ color: "red" }}>{error}</p>}

      {position && airData ? (
        <MapContainer
          center={position}
          zoom={13}
          scrollWheelZoom={true}
          style={{ height: "90vh" }}
        >
          <TileLayer
            attribution='&copy; <a href="https://osm.org/copyright">OpenStreetMap</a>'
            url="https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png"
          />
          <Marker position={position}>
            <Popup>
              <strong>Air Quality</strong>
              <br />
              AQI:{" "}
              <span style={{ color: getColor(airData.aqi || 0) }}>
                {airData.aqi || "N/A"}
              </span>
              <br />
              PM2.5: {airData.pm2_5 || "N/A"} µg/m³
              <br />
              PM10: {airData.pm10 || "N/A"} µg/m³
              <br />
              Status: {airData.status || "N/A"}
            </Popup>
          </Marker>
        </MapContainer>
      ) : !error ? (
        <p>Loading air quality data...</p>
      ) : null}
    </div>
  );
}

export default App;
