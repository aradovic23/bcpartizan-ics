const axios = require('axios');
const config = require('./config');

function getVenueForTeam(teamName) {
  const venues = {
    'Partizan': 'Stark Arena, Belgrade, Serbia',
    'Crvena zvezda': 'Aleksandar Nikolic Hall, Belgrade, Serbia',
    'Real Madrid': 'WiZink Center, Madrid, Spain',
    'Barcelona': 'Palau Blaugrana, Barcelona, Spain',
    'Panathinaikos': 'OAKA, Athens, Greece',
    'Olympiacos': 'Peace and Friendship Stadium, Piraeus, Greece',
    'Fenerbahce': 'Ulker Sports Arena, Istanbul, Turkey',
    'Anadolu Efes': 'Sinan Erdem Dome, Istanbul, Turkey',
    'Maccabi Tel Aviv': 'Menora Mivtachim Arena, Tel Aviv, Israel',
    'Zalgiris': 'Zalgirio Arena, Kaunas, Lithuania',
    'Bayern Munich': 'Audi Dome, Munich, Germany',
    'ALBA Berlin': 'Mercedes-Benz Arena, Berlin, Germany',
    'Milano': 'Mediolanum Forum, Milan, Italy',
    'Virtus Bologna': 'Segafredo Arena, Bologna, Italy',
    'Monaco': 'Salle Gaston Médecin, Monaco',
    'Paris': 'Adidas Arena, Paris, France',
    'Baskonia': 'Buesa Arena, Vitoria-Gasteiz, Spain',
    'Valencia': 'La Fonteta, Valencia, Spain'
  };

  for (const [team, venue] of Object.entries(venues)) {
    if (teamName?.includes(team)) {
      return venue;
    }
  }

  return `${teamName} Arena`;
}

function parseFlashScoreData(html) {
  const games = [];

  const dataMatch = html.match(/data:\s*`([^`]+)`/);
  if (!dataMatch) return games;

  const dataString = dataMatch[1];
  const gameBlocks = dataString.split('~AA÷');

  for (let i = 0; i < gameBlocks.length; i++) {
    const block = gameBlocks[i];
    if (!block.includes('PTZ') && !block.includes(config.team)) continue;

    const fields = {};
    const pairs = block.split('¬');

    for (const pair of pairs) {
      const [key, value] = pair.split('÷');
      if (key && value) {
        fields[key] = value;
      }
    }

    const nextBlock = i < gameBlocks.length - 1 ? gameBlocks[i + 1] : '';
    const venueMatch = nextBlock.match(/AM÷([^¬]+)/);
    let venue = venueMatch ? venueMatch[1].replace('Neutral location - ', '').replace(/\.$/, '') : null;

    const timestamp = fields['AD'];
    const homeTeam = fields['AE'];
    const awayTeam = fields['AF'];
    const round = fields['ER'];
    const competition = fields['ZA'];

    if ((homeTeam?.includes(config.team) || awayTeam?.includes(config.team)) && timestamp) {
      const date = new Date(parseInt(timestamp) * 1000);

      if (!venue || venue.includes('TBD')) {
        venue = getVenueForTeam(homeTeam);
      }

      games.push({
        competition: competition || 'Euroleague',
        homeTeam: homeTeam || 'Unknown',
        awayTeam: awayTeam || 'Unknown',
        date: date.toISOString().split('T')[0],
        time: date.toTimeString().slice(0, 5),
        venue: venue,
        location: venue,
        round: round || ''
      });
    }
  }

  return games;
}

async function fetchFlashScoreSchedule(url, competition) {
  try {
    const response = await axios.get(url, {
      headers: {
        'User-Agent': 'Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36'
      }
    });

    const games = parseFlashScoreData(response.data);
    return games.map(game => ({ ...game, competition, source: 'flashscore' }));
  } catch (error) {
    console.error(`Error fetching ${competition} schedule from FlashScore:`, error.message);
    return [];
  }
}

function deduplicateGames(games) {
  const seen = new Map();

  for (const game of games) {
    const key = `${game.date}-${game.time}-${game.homeTeam}-${game.awayTeam}`;
    if (!seen.has(key)) {
      seen.set(key, game);
    }
  }

  return Array.from(seen.values());
}

async function fetchEuroleagueSchedule() {
  try {
    const response = await axios.get(config.euroleagueApiUrl, {
      params: { teamCode: config.teamCodeEuroleague },
      headers: {
        'User-Agent': 'Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36'
      }
    });

    if (response.data?.status !== 'success' || !response.data?.data) {
      console.error('Euroleague API returned no data');
      return [];
    }

    const games = response.data.data;
    const now = new Date();

    const futureGames = games
      .filter(game => new Date(game.date) >= now)
      .map(game => ({
        competition: 'Euroleague',
        homeTeam: game.home.abbreviatedName || game.home.name,
        awayTeam: game.away.abbreviatedName || game.away.name,
        date: new Date(game.date).toISOString().split('T')[0],
        time: new Date(game.date).toTimeString().slice(0, 5),
        venue: game.venue?.name && game.venue?.address
          ? `${game.venue.name}, ${game.venue.address}`
          : game.venue?.name || getVenueForTeam(game.home.name),
        location: game.venue?.name && game.venue?.address
          ? `${game.venue.name}, ${game.venue.address}`
          : game.venue?.name || getVenueForTeam(game.home.name),
        round: game.round?.name || '',
        source: 'euroleague-api'
      }));

    return futureGames;
  } catch (error) {
    console.error('Error fetching Euroleague API:', error.message);
    return [];
  }
}

async function fetchABALeagueSchedule() {
  const [fixtures, results] = await Promise.all([
    fetchFlashScoreSchedule(config.flashscoreABAFixturesUrl, 'ABA League'),
    fetchFlashScoreSchedule(config.flashscoreABAResultsUrl, 'ABA League')
  ]);

  const allGames = [...fixtures, ...results];
  const now = new Date();

  const futureGames = allGames.filter(game => {
    const gameDate = new Date(game.date + ' ' + game.time);
    return gameDate >= now;
  });

  return deduplicateGames(futureGames);
}

async function fetchAllSchedules() {
  const [euroleagueGames, abaGames] = await Promise.all([
    fetchEuroleagueSchedule(),
    fetchABALeagueSchedule()
  ]);

  const allGames = [...euroleagueGames, ...abaGames];

  return allGames.sort((a, b) => {
    const dateA = new Date(a.date + ' ' + a.time);
    const dateB = new Date(b.date + ' ' + b.time);
    return dateA - dateB;
  });
}

function getMockSchedule() {
  const now = new Date();
  const games = [];

  for (let i = 0; i < 10; i++) {
    const gameDate = new Date(now);
    gameDate.setDate(gameDate.getDate() + (i * 7));

    const opponents = ['Real Madrid', 'Barcelona', 'Olimpia Milano', 'Fenerbahce', 'Crvena Zvezda', 'Maccabi', 'Bayern Munich', 'Zalgiris'];
    const opponent = opponents[i % opponents.length];
    const isHome = i % 2 === 0;

    games.push({
      competition: i % 3 === 0 ? 'Euroleague' : 'ABA League',
      homeTeam: isHome ? 'Partizan' : opponent,
      awayTeam: isHome ? opponent : 'Partizan',
      date: gameDate.toISOString().split('T')[0],
      time: '20:00',
      venue: isHome ? 'Stark Arena, Belgrade, Serbia' : `${opponent} Arena`,
      location: isHome ? 'Stark Arena, Belgrade, Serbia' : `${opponent} Arena`,
      source: 'mock'
    });
  }

  return games;
}

module.exports = {
  fetchAllSchedules,
  getMockSchedule
};

