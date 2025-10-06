# 🏀 KK Partizan Basketball Schedule - ICS Calendar

Automatically sync KK Partizan basketball games across all competitions to your calendar app. This app fetches schedules from Euroleague, ABA League, and other competitions, and provides an ICS subscription URL you can add to any calendar application.

## Features

- ✅ **Multi-competition support**: Euroleague, ABA League, domestic competitions
- ✅ **Smart event titles**: Includes competition name and matchup (e.g., "Euroleague - Partizan vs Real Madrid")
- ✅ **Venue information**: Full location details for each game
- ✅ **30-minute reminders**: Get notified before each game starts
- ✅ **Auto-refresh**: Schedule updates every 2 days
- ✅ **Caching**: Reduces unnecessary API calls
- ✅ **Real-time data**: Fetches from FlashScore for accurate game info

## Quick Start

### Installation

```bash
npm install
```

### Configuration

Copy the example environment file:

```bash
cp .env.example .env
```

Edit `.env` if you want to customize:
- `PORT`: Server port (default: 3000)
- `CACHE_REFRESH_INTERVAL`: Cron expression for refresh schedule (default: every 2 days at midnight)

### Running the Server

```bash
node src/index.js
```

Or add to `package.json`:

```json
"scripts": {
  "start": "node src/index.js",
  "dev": "nodemon src/index.js"
}
```

Then run:

```bash
npm start
```

The server will start on `http://localhost:3000`

## Usage

### Subscribe to Calendar

1. Start the server
2. Visit `http://localhost:3000` in your browser
3. Copy the calendar subscription URL: `http://localhost:3000/calendar.ics`
4. Add to your calendar app:

#### Apple Calendar
1. Open Calendar app
2. File → New Calendar Subscription
3. Paste the URL
4. Set refresh interval (recommended: hourly or daily)

#### Google Calendar
1. Open Google Calendar
2. Click "+" next to "Other calendars"
3. Select "From URL"
4. Paste the URL

#### Outlook
1. Open Outlook Calendar
2. Add Calendar → Subscribe from web
3. Paste the URL

### API Endpoints

- `GET /` - Home page with instructions
- `GET /calendar.ics` - ICS calendar file (subscribe to this URL)
- `GET /games` - JSON list of upcoming games
- `GET /refresh` - Manually trigger schedule refresh

## Event Details

Each calendar event includes:

- **Title**: `[Competition] - [Home Team] vs [Away Team]`
  - Example: "Euroleague - Partizan vs Real Madrid"
- **Location**: Full venue name and address
- **Duration**: 2 hours (typical game duration)
- **Reminder**: 30 minutes before game starts
- **Status**: Confirmed

## Data Sources

The app fetches schedules from **reliable official sources**:

1. **Euroleague Official API** (`feeds.incrowdsports.com`) ⭐
   - Complete season schedule (all 36-38 regular season games)
   - Official venue information with full addresses
   - Most reliable and comprehensive source
   - Filtered specifically for Partizan games via `teamCode=PAR`

2. **FlashScore** (ABA League only)
   - Real-time fixtures and results for ABA League
   - Rolling 2-3 month window for domestic games
   - Good coverage for recently announced ABA games

### Data Coverage

**Current Implementation:**
- ✅ Fetches **~49 upcoming Partizan games**
  - **36 Euroleague games** (complete regular season through April 2026)
  - **13 ABA League games** (2-3 months ahead)
- ✅ Euroleague: **100% complete season coverage** via official API
- ✅ ABA League: 2-3 months rolling window via FlashScore
- ✅ Auto-refreshes every 2 days to include newly announced ABA games
- ✅ Filters out past games automatically
- ✅ Complete venue information with full addresses

## Project Structure

```
bcpartizan-ics/
├── src/
│   ├── index.js         # Express server & main app
│   ├── scraper.js       # Schedule fetching logic
│   ├── icsGenerator.js  # ICS calendar generation
│   ├── cache.js         # Caching mechanism
│   └── config.js        # Configuration
├── data/
│   └── cache.json       # Cached schedule data
├── .env                 # Environment variables
├── .gitignore
└── package.json
```

## Development

### Adding New Data Sources

Edit `src/scraper.js` and add a new fetch function:

```javascript
async function fetchNewCompetition() {
  // Your scraping logic
  return games;
}
```

Then add it to `fetchAllSchedules()`:

```javascript
const [euroleague, aba, newComp] = await Promise.all([
  fetchEuroleagueSchedule(),
  fetchABALeagueSchedule(),
  fetchNewCompetition()
]);
```

### Customizing Event Format

Edit `src/icsGenerator.js` to modify:
- Event titles
- Duration
- Reminder timing
- Description format

## Deployment

---

### 📱 How Friends Can Subscribe

**Apple Calendar (iPhone/Mac):**
1. Open Calendar app
2. File → New Calendar Subscription (Mac) or Settings → Accounts → Add Account → Other (iPhone)
3. Paste: `https://YOUR-APP.onrender.com/calendar.ics`
4. Set refresh: Hourly or Daily

**Google Calendar:**
1. Open [calendar.google.com](https://calendar.google.com)
2. Click **"+"** next to "Other calendars"
3. Select **"From URL"**
4. Paste: `https://YOUR-APP.onrender.com/calendar.ics`

**Outlook:**
1. Open Outlook Calendar
2. Add Calendar → Subscribe from web
3. Paste: `https://YOUR-APP.onrender.com/calendar.ics`

---

### ⚠️ Important Notes About Free Tier

**Service Sleep:**
- Free tier spins down after 15 min of inactivity
- First request after sleep takes ~30-60 seconds (cold start)
- **This is fine for calendars!** Calendar apps will wait for the response

**Keeping It Active (Optional):**
If you want faster responses, use [UptimeRobot](https://uptimerobot.com) (free) to ping your service every 10 minutes:
- Add Monitor → HTTP(s)
- URL: `https://YOUR-APP.onrender.com/`
- Interval: 5 minutes

**No Configuration Needed:**
- PORT is automatically set by Render
- All environment variables are pre-configured
- Just deploy and go! 🚀

---

### Alternative Deployment Options

<details>
<summary><b>Railway.app</b></summary>

1. Sign up at [railway.app](https://railway.app)
2. New Project → Deploy from GitHub
3. Select your repository
4. Railway auto-deploys
5. Get public URL from dashboard

**Cost:** $5 credit/month (free tier)
</details>

<details>
<summary><b>Fly.io</b></summary>

1. Install CLI: `brew install flyctl`
2. Run: `fly launch`
3. Follow prompts
4. Deploy: `fly deploy`

**Cost:** Free tier includes 3 VMs
</details>

<details>
<summary><b>Docker (Self-hosted)</b></summary>

Create `Dockerfile`:
```dockerfile
FROM node:18-alpine
WORKDIR /app
COPY package*.json ./
RUN npm install --production
COPY . .
EXPOSE 3000
CMD ["npm", "start"]
```

Build and run:
```bash
docker build -t partizan-ics .
docker run -p 3000:3000 partizan-ics
```
</details>

## Troubleshooting

**No games showing up?**
- Check if the scraper sources are accessible
- The app uses mock data as fallback
- Visit `/games` endpoint to see raw data

**Calendar not updating?**
- Calendar apps have their own refresh intervals
- Try manually refreshing your calendar
- Use `/refresh` endpoint to force update

**Calendar subscription not working?**
- Ensure your server is publicly accessible (not localhost)
- Use HTTPS in production for better compatibility
- Check if your calendar app supports HTTP subscriptions

## License

ISC

## Support

For issues or questions about KK Partizan, visit their official website.

---

**Partizane napadaj!** 🖤🤍
