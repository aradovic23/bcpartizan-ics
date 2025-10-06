# Agent Guidelines for bcpartizan-ics

## Build/Run/Test Commands
- **Start server**: `npm start` or `node src/index.js`
- **Dev mode**: `npm run dev`
- **No tests configured** - project does not have a test suite

## Code Style

**Language**: JavaScript (Node.js), no TypeScript

**Imports**: Use CommonJS (`require`/`module.exports`)

**Error Handling**: Use try-catch blocks, log errors with `console.error()`, return fallback data (mock schedule) on failures

**Naming**: camelCase for functions/variables, SCREAMING_SNAKE_CASE for constants (e.g., `CACHE_FILE`)

**Formatting**: 2-space indentation, single quotes for strings where possible

**Async/Await**: Prefer async/await over promises, use `Promise.all()` for parallel operations

**Configuration**: All config in `src/config.js`, use `process.env` with defaults

**File Structure**: Modular - separate files for scraper, ICS generation, cache, config

**Comments**: Minimal - code should be self-documenting
