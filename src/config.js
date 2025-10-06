module.exports = {
  port: process.env.PORT || 3000,
  cacheRefreshInterval: process.env.CACHE_REFRESH_INTERVAL || '0 0 */2 * *',
  cacheTTL: 2 * 24 * 60 * 60 * 1000,
  flashscoreABAFixturesUrl: 'https://www.flashscore.com/basketball/europe/aba-league/fixtures/',
  flashscoreABAResultsUrl: 'https://www.flashscore.com/basketball/europe/aba-league/results/',
  euroleagueApiUrl: 'https://feeds.incrowdsports.com/provider/euroleague-feeds/v2/competitions/E/seasons/E2025/games',
  team: 'Partizan',
  teamCodeEuroleague: 'PAR',
  defaultGameDuration: 2
};
