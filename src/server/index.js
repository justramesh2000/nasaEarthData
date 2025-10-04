const express = require("express");
const cors = require("cors");
const axios = require("axios");
const fs = require('fs');
const csv = require('csv-parser');
const geolib = require('geolib')

const app = express();
const PORT = process.env.PORT || 5050;
app.use(cors());
app.use(express.json());

let airQualityCache = null;

let data = [];
fs.createReadStream('data.csv')
  .pipe(csv())
  .on('data', (row) => {
    row.latitude = parseFloat(row.latitude);
    row.longitude = parseFloat(row.longitude);
    data.push(row);
  })
  .on('end', () => {
    console.log(`Loaded ${data.length} records from csv.`);
  });

app.get("/", (req, res) => {
  res.send("Nasa weather API is running");
});

app.get("/api/data", (req, res) => {
  const { lat, lon } = req.query;
  if (!lat || !lon) {
    return res
      .status(400)
      .json({ error: "Latitude and Longitude are required" });
  }
  const cacheKey = `${lat},${lon}`;
  const cacheEntry = airQualityCache[cacheKey];
  const now = Date.now();
  const isFresh =
    cacheEntry &&
    now - new Date(cacheEntry.lastUpdated.getTime()) < 5 * 60 * 1000;
  if (isFresh) {
    return res.json({
      ...cacheEntry,
      cached: true,
    });
  }

  try {
    const apiKey = process.env.AIR_API_KEY;
    const url = `https://api.example.com/air?lat=${lat}&lon=${lon}&key=${apiKey}`;
    const response = await.axios.get(url);
  } catch (error) {
    console.log("Failed to fetch air quality data:", error.message);
    res.status(500).json({ error: "Failed to fetch air quality data" });
  }
});

app.get('/api/csvdata', (req,res) => {
  const {lat, lng, radius} = req.query;
  if (!lat || !lng || !radius) {
    return res.status(400).json({error: 'Missing lat, lng, or radius'});
  }

  const center = { latitude: parseFloat(lat), longitude: parseFloat(lng) };
  const radiusMiles = parseFloat(radius);
  const radiusMeters = radiusMiles * 1609.34;

  const filtered = data.filter((point) => {
    return geolib.isPointWithinRadius(
      { latitude: point.latitude, longitude: point.longitude },
      center,
      radiusMeters
    );
  });
  res.json(filtered)
})

app.listen(PORT, () => {
  console.log(`Weather server is listening on port ${PORT}`);
});
