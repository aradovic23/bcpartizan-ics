# Agent Guidelines for bcpartizan-ics

## Build/Run/Test Commands
- **Build**: `go build -o partizan-ics`
- **Run**: `./partizan-ics`
- **Test**: `go test ./...` (no tests configured yet)
- **Format**: `go fmt ./...`
- **Lint**: `go vet ./...`

## Code Style

**Language**: Go 1.23+

**Package Organization**: 
- Each package in its own directory
- Package name matches directory name
- Main package in root (`main.go`)

**Imports**: 
- Standard library first
- Third-party packages second
- Local packages last
- Group with blank lines

**Error Handling**: 
- Always check errors, never ignore
- Log errors with `log.Printf()` or `log.Println()`
- Return fallback data (mock schedule) on scraper failures
- Use error wrapping with `fmt.Errorf()` when appropriate

**Naming**: 
- Exported: PascalCase (e.g., `FetchAllSchedules`, `Config`)
- Unexported: camelCase (e.g., `fetchEuroleagueSchedule`, `cacheMutex`)
- Constants: PascalCase or camelCase depending on export
- Acronyms: Keep uppercase (e.g., `URL`, `HTTP`, `ICS`)

**Formatting**: 
- Use `gofmt` (tabs for indentation)
- Follow standard Go conventions
- Keep line length reasonable

**Concurrency**:
- Use `sync.RWMutex` for shared cache access
- Goroutines for parallel HTTP requests where appropriate
- Channels for communication between goroutines

**Configuration**: 
- All config in `config/config.go`
- Use `os.Getenv()` with defaults
- Load `.env` with `godotenv.Load()`

**File Structure**: 
- Modular packages: `config/`, `cache/`, `scraper/`, `ics/`, `types/`
- Each package has clear responsibility
- Shared types in `types/` package

**Comments**: 
- Minimal inline comments - code should be self-documenting
- Package-level documentation for exported functions
- Document complex logic only

**Dependencies**:
- `github.com/arran4/golang-ical` - ICS generation
- `github.com/joho/godotenv` - .env file loading
- `github.com/robfig/cron/v3` - Scheduled tasks
