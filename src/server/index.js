const express = require("express");
const cors = require("cors");
const axios = require("axios");

const app = express();
const PORT = process.env.PORT || 5050;
app.use(cors());
app.use(express.json());

let airQualityCache = null;

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

app.listen(PORT, () => {
  console.log(`Weather server is listening on port ${PORT}`);
});
