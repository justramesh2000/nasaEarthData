import logo from "./logo.svg";
import "./App.css";
import { useEffect, useState } from "react";
import "leaflet/dist/leaflet.css";
import { MapContainer, TileLayer, Marker, Popup, userMapEvents, useMapEvents } from "react-leaflet";
import L from "leaflet";

function App() {
  const [clickedPosition, setClickedPosition] = useState(null);
  const [userPosition, setUserPosition] = useState(null);
  const [radius, setRadius] = useState(5);
  const [airDataPoints, setAirDataPoints] = useState([]);
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
        setUserPosition([latitude, longitude]);
      },
      () => {
        setError("Failed to get your location");
      }
    );
  }, []);

    const fetchAirData = async (lat, lon, radius = 5) => {
            try {
              const res = await fetch(`/api/csvdata?lat=${lat}&lng=${lon}&radius=${radius}`);
              const data = await res.json();
              setAirDataPoints(data);
            } catch (err) {
              console.error("Failed to fetch air quality data:", err);
              setError("Failed to fetch air quality data");
            }
          };
          
  function MapClickHandler({ radius }) {
    useMapEvents({
      click(e) {
        const { lat, lng } = e.latlng;
        setClickedPosition([lat, lng]);
        fetchAirData(lat, lng, radius);
      }
    });
    return null;
  }        
  const getColor = (aqi) => {
    if (aqi <= 50) return "green";
    if (aqi <= 100) return "yellow";
    if (aqi <= 150) return "orange";
    if (aqi <= 200) return "red";
    return "purple";
  };

  return (
    <div>
      <h1>NASA space App – Weather Visualization</h1>
      {error && <p style={{ color: "red" }}>{error}</p>}
        <div style={{ padding: "1rem" }}>
          <label>
            Radius:
            <select
              value={radius}
              onChange={(e) => setRadius(parseInt(e.target.value))}
              style={{ marginLeft: "0.5rem" }}
            >
              {[5, 6, 7, 8, 9, 10].map((mile) => (
                <option key={mile} value={mile}>
                  {mile} mile{mile > 1 ? "s" : ""}
                </option>
              ))}
            </select>
          </label>
      </div>
      {userPosition ? (
        <MapContainer 
          center={userPosition}
          zoom={8}
          scrollWheelZoom={true}
          style={{ height: "90vh" }}
        >
          <TileLayer
            attribution='&copy; <a href="https://osm.org/copyright">OpenStreetMap</a>'
            url="https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png"
          />
          <MapClickHandler radius={radius}/>

            {/* Show markers for each data point */}
            {Array.isArray(airDataPoints) && airDataPoints.map((point, index) => (
              <Marker
                key={index}
                position={[point.latitude, point.longitude]}
                icon={L.icon({ iconUrl: L.Icon.Default.prototype.options.iconUrl })}
              >
                <Popup>
                  <strong>Air Quality</strong>
                  <br />
                  UV Index: {point.uv_aerosol_index}
                  <br />
                  AOD: {point.aod_354 ?? point.aod_388 ?? point.aod_500 ?? "N/A"}
                  <br />
                  Layer Height: {point.aerosol_layer_height_m} m
                  <br />
                  Time: {new Date(point.time_utc).toLocaleString()}
                  <br />
                  Lat: {point.latitude.toFixed(5)}
                  <br />
                  Lon: {point.longitude.toFixed(5)}
                </Popup>
              </Marker>
            ))}
        </MapContainer>
      ) : (
        <p>Loading map...</p>
      )}
      </div>
  );
}

export default App;
