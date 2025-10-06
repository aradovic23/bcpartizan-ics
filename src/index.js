const express = require('express');
const cron = require('node-cron');
require('dotenv').config();

const config = require('./config');
const { fetchAllSchedules, getMockSchedule } = require('./scraper');
const { generateCalendar } = require('./icsGenerator');
const { loadCache, saveCache } = require('./cache');

const app = express();

let cachedGames = null;

async function refreshSchedule() {
  console.log('Refreshing schedule...');
  try {
    let games = await fetchAllSchedules();
    
    if (games.length === 0) {
      console.log('No games fetched from real sources, using mock data');
      games = getMockSchedule();
    }
    
    cachedGames = games;
    await saveCache(games);
    console.log(`Schedule refreshed: ${games.length} games found`);
    return games;
  } catch (error) {
    console.error('Error refreshing schedule:', error.message);
    return cachedGames || getMockSchedule();
  }
}

async function getGames() {
  if (cachedGames) {
    return cachedGames;
  }
  
  const cached = await loadCache();
  if (cached) {
    cachedGames = cached;
    return cached;
  }
  
  return await refreshSchedule();
}

app.get('/calendar.ics', async (req, res) => {
  try {
    const games = await getGames();
    const icsContent = generateCalendar(games);
    
    res.setHeader('Content-Type', 'text/calendar; charset=utf-8');
    res.setHeader('Content-Disposition', 'attachment; filename="partizan-schedule.ics"');
    res.send(icsContent);
  } catch (error) {
    console.error('Error generating calendar:', error);
    res.status(500).send('Error generating calendar');
  }
});

app.get('/games', async (req, res) => {
  try {
    const games = await getGames();
    res.json({ count: games.length, games });
  } catch (error) {
    console.error('Error fetching games:', error);
    res.status(500).json({ error: 'Error fetching games' });
  }
});

app.get('/refresh', async (req, res) => {
  try {
    const games = await refreshSchedule();
    res.json({ message: 'Schedule refreshed', count: games.length });
  } catch (error) {
    console.error('Error refreshing:', error);
    res.status(500).json({ error: 'Error refreshing schedule' });
  }
});

app.get('/', (req, res) => {
  res.send(`
    <html>
      <head><title>Partizan Basketball Schedule</title></head>
      <body style="font-family: Arial, sans-serif; max-width: 800px; margin: 50px auto; padding: 20px;">
        <h1>🏀 KK Partizan Schedule Calendar</h1>
        <p>Subscribe to Partizan basketball games across all competitions.</p>
        
        <h2>Calendar Subscription URL:</h2>
        <code style="background: #f4f4f4; padding: 10px; display: block; margin: 10px 0;">
          ${req.protocol}://${req.get('host')}/calendar.ics
        </code>
        
        <h3>How to subscribe:</h3>
        <ul>
          <li><strong>Apple Calendar:</strong> File → New Calendar Subscription → Paste URL</li>
          <li><strong>Google Calendar:</strong> Settings → Add Calendar → From URL → Paste URL</li>
          <li><strong>Outlook:</strong> Add Calendar → Subscribe from web → Paste URL</li>
        </ul>
        
        <h3>Features:</h3>
        <ul>
          <li>✅ All Partizan games across Euroleague, ABA League, and domestic competitions</li>
          <li>✅ Game title includes competition name and teams</li>
          <li>✅ Venue location information</li>
          <li>✅ 30-minute reminder before each game</li>
          <li>✅ Automatically updates every 6 hours</li>
        </ul>
        
        <p><a href="/games">View upcoming games (JSON)</a> | <a href="/refresh">Force refresh</a></p>
      </body>
    </html>
  `);
});

cron.schedule(config.cacheRefreshInterval, () => {
  refreshSchedule();
});

refreshSchedule();

app.listen(config.port, () => {
  console.log(`Partizan ICS Calendar Server running on port ${config.port}`);
  console.log(`Calendar URL: http://localhost:${config.port}/calendar.ics`);
});
