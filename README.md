# CIDR Viewer

A web application for visualizing IPv4 CIDR ranges and identifying gaps and unused IP addresses.

## Architecture

- **Frontend**: Svelte + Vite (Port 5173)
- **Backend**: Go + Gin (Port 8080)

## Quick Start

### Configuration

The frontend can be configured using environment variables. Copy the example configuration:

```bash
cd frontend
cp env.example .env
```

Edit `.env` to set your API base URL:

```env
VITE_API_BASE_URL=http://localhost:8080/api
VITE_API_TIMEOUT=30000
```

### Backend

```bash
cd backend
go run main.go
```

### Frontend

```bash
cd frontend
npm install
npm run dev
```

## Runtime Configuration

The application uses a JSON configuration file that can be modified without rebuilding the application.

Edit `public/app-config.json` to configure the API:

```json
{
  "apiBaseUrl": "http://localhost:8080/api",
  "apiTimeout": 30000
}
```

### Configuration Options

- `apiBaseUrl`: Backend API base URL
- `apiTimeout`: API request timeout in milliseconds

### Environment-Specific Configuration

**Development:**
```json
{
  "apiBaseUrl": "http://localhost:8080/api",
  "apiTimeout": 30000
}
```

**Production:**
```json
{
  "apiBaseUrl": "https://your-production-api.com/api",
  "apiTimeout": 10000
}
```

### Deployment

The same build can be deployed to multiple environments by simply changing the `app-config.json` file:

```bash
# Build once
npm run build

# Deploy to staging
echo '{"apiBaseUrl": "https://staging-api.com/api", "apiTimeout": 15000}' > dist/app-config.json

# Deploy to production  
echo '{"apiBaseUrl": "https://prod-api.com/api", "apiTimeout": 10000}' > dist/app-config.json
```

## Project Structure

```
cidr-viewer/
├── backend/              # Go backend API
│   ├── main.go          # Server entry point
│   ├── handlers/        # HTTP handlers
│   ├── models/          # Data models
│   └── utils/           # Utility functions
├── frontend/            # Svelte frontend
│   ├── src/            # Source code
│   ├── public/         # Static assets
│   └── package.json    # Dependencies
└── REQUIREMENTS.md     # Project requirements
```

## Development Phases

- **Phase 1 (MVP)**: Basic CIDR input, validation, gap analysis
- **Phase 2**: Enhanced visualization with interactive diagrams
- **Phase 3**: Advanced features (file import, export)
- **Phase 4**: Polish and optimization