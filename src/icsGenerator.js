const ics = require('ics');
const config = require('./config');

function parseDateTime(dateStr, timeStr) {
  const date = new Date(`${dateStr} ${timeStr}`);
  return {
    year: date.getFullYear(),
    month: date.getMonth() + 1,
    day: date.getDate(),
    hour: date.getHours(),
    minute: date.getMinutes()
  };
}

function generateCalendar(games) {
  const events = games.map(game => {
    const startDateTime = parseDateTime(game.date, game.time);
    
    const title = `${game.competition} - ${game.homeTeam} vs ${game.awayTeam}`;
    
    const description = `${game.competition}\n${game.homeTeam} vs ${game.awayTeam}`;
    
    return {
      start: [
        startDateTime.year,
        startDateTime.month,
        startDateTime.day,
        startDateTime.hour,
        startDateTime.minute
      ],
      duration: { hours: config.defaultGameDuration },
      title: title,
      description: description,
      location: game.location || game.venue,
      status: 'CONFIRMED',
      busyStatus: 'BUSY',
      alarms: [
        {
          action: 'display',
          description: `${title} starts in 30 minutes`,
          trigger: { minutes: 30, before: true }
        }
      ]
    };
  });

  const { error, value } = ics.createEvents(events);

  if (error) {
    console.error('Error generating ICS:', error);
    throw new Error('Failed to generate calendar');
  }

  return value;
}

module.exports = {
  generateCalendar
};
