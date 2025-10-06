const fs = require('fs').promises;
const path = require('path');
const config = require('./config');

const CACHE_FILE = path.join(__dirname, '../data/cache.json');

async function loadCache() {
  try {
    const data = await fs.readFile(CACHE_FILE, 'utf8');
    const cache = JSON.parse(data);
    
    const now = Date.now();
    if (cache.timestamp && (now - cache.timestamp) < config.cacheTTL) {
      return cache.games;
    }
    
    return null;
  } catch (error) {
    return null;
  }
}

async function saveCache(games) {
  const cache = {
    timestamp: Date.now(),
    games
  };
  
  try {
    await fs.mkdir(path.dirname(CACHE_FILE), { recursive: true });
    await fs.writeFile(CACHE_FILE, JSON.stringify(cache, null, 2));
  } catch (error) {
    console.error('Error saving cache:', error.message);
  }
}

module.exports = {
  loadCache,
  saveCache
};
