# CIDR Viewer - Requirements Document

## Project Overview

The CIDR Viewer is a web application designed to help network administrators and IT professionals visualize IPv4 CIDR (Classless Inter-Domain Routing) ranges, identify gaps in IP allocation, and discover unused IP addresses within their networks.

## Technology Stack

### Frontend
- **Framework**: Svelte
- **Build Tool**: Vite
- **Styling**: CSS3 with modern features (Grid, Flexbox)
- **Visualization**: D3.js or Chart.js for network visualization
- **UI Components**: Custom Svelte components

### Backend
- **Language**: Go (Golang)
- **Framework**: Gin or Echo for REST API
- **Database**: SQLite for development, PostgreSQL for production (optional)
- **Libraries**: 
  - `net` package for IP address manipulation
  - `encoding/json` for API responses
  - CORS middleware for cross-origin requests

## Functional Requirements

### 1. CIDR Input Interface
- **FR-1.1**: Single-page application with a clean, intuitive interface
- **FR-1.2**: Text area or input field for entering multiple CIDR ranges
- **FR-1.3**: Support for IPv4 CIDR notation (e.g., 192.168.1.0/24, 10.0.0.0/8)
- **FR-1.4**: Real-time validation of CIDR format
- **FR-1.5**: Ability to add/remove individual CIDR ranges
- **FR-1.6**: Import functionality for CSV or text files containing CIDR ranges
- **FR-1.7**: Example CIDR ranges provided for demonstration

### 2. CIDR Processing and Analysis
- **FR-2.1**: Parse and validate IPv4 CIDR ranges
- **FR-2.2**: Calculate network address, broadcast address, and usable IP range
- **FR-2.3**: Identify overlapping CIDR ranges
- **FR-2.4**: Detect gaps between CIDR ranges
- **FR-2.5**: Calculate unused IP addresses within specified ranges
- **FR-2.6**: Support for supernet and subnet analysis
- **FR-2.7**: Aggregation of adjacent CIDR ranges

### 3. Visualization Features
- **FR-3.1**: Interactive network topology diagram
- **FR-3.2**: Color-coded visualization showing:
  - Allocated CIDR ranges (green)
  - Gaps/unused ranges (red)
  - Overlapping ranges (orange/yellow)
- **FR-3.3**: Hierarchical tree view of network structure
- **FR-3.4**: Linear timeline view of IP address space
- **FR-3.5**: Zoom and pan functionality for large address spaces
- **FR-3.6**: Tooltip information on hover showing:
  - CIDR range details
  - Number of available IPs
  - Subnet mask information
- **FR-3.7**: Legend explaining color coding and symbols

### 4. Data Export and Reporting
- **FR-4.1**: Export visualization as PNG/SVG
- **FR-4.2**: Generate summary report in PDF format
- **FR-4.3**: Export gap analysis as CSV
- **FR-4.4**: Print-friendly view of analysis results

### 5. User Experience Features
- **FR-5.1**: Responsive design for desktop and tablet devices
- **FR-5.2**: Dark/light theme toggle
- **FR-5.3**: Save and load CIDR configurations locally
- **FR-5.4**: Undo/redo functionality for CIDR modifications
- **FR-5.5**: Search functionality within visualization
- **FR-5.6**: Filter options (show only gaps, show only overlaps, etc.)

## Non-Functional Requirements

### Performance
- **NFR-1.1**: Application should handle up to 1000 CIDR ranges without performance degradation
- **NFR-1.2**: Real-time validation and visualization updates (< 500ms response time)
- **NFR-1.3**: Frontend bundle size should be optimized for fast loading

### Scalability
- **NFR-2.1**: Backend API should support concurrent requests from multiple users
- **NFR-2.2**: Modular architecture allowing for future feature additions

### Security
- **NFR-3.1**: Input validation and sanitization to prevent injection attacks
- **NFR-3.2**: CORS configuration for secure cross-origin requests
- **NFR-3.3**: No sensitive network information stored persistently

### Compatibility
- **NFR-4.1**: Support for modern browsers (Chrome 90+, Firefox 88+, Safari 14+)
- **NFR-4.2**: Cross-platform compatibility (Windows, macOS, Linux)

### Usability
- **NFR-5.1**: Intuitive interface requiring minimal learning curve
- **NFR-5.2**: Accessible design following WCAG 2.1 guidelines
- **NFR-5.3**: Error messages should be clear and actionable

## API Specifications

### Endpoints

#### POST /api/analyze
Analyze CIDR ranges and return gap analysis

**Request Body:**
```json
{
  "cidrs": ["192.168.1.0/24", "10.0.0.0/8", "172.16.0.0/12"]
}
```

**Response:**
```json
{
  "analysis": {
    "valid_cidrs": [...],
    "invalid_cidrs": [...],
    "gaps": [...],
    "overlaps": [...],
    "summary": {
      "total_ips": 16777216,
      "allocated_ips": 16711680,
      "available_ips": 65536
    }
  }
}
```

#### POST /api/validate
Validate CIDR format

**Request Body:**
```json
{
  "cidr": "192.168.1.0/24"
}
```

**Response:**
```json
{
  "valid": true,
  "network": "192.168.1.0",
  "broadcast": "192.168.1.255",
  "mask": "255.255.255.0",
  "total_ips": 256,
  "usable_ips": 254
}
```

#### GET /api/health
Health check endpoint

**Response:**
```json
{
  "status": "healthy",
  "timestamp": "2025-09-17T10:30:00Z"
}
```

## User Interface Mockup Description

### Main Layout
1. **Header**: Application title and theme toggle
2. **Input Panel** (Left side, 30% width):
   - CIDR input textarea
   - Add/Remove buttons
   - Import file button
   - Validation status indicators
3. **Visualization Panel** (Right side, 70% width):
   - Interactive network diagram
   - Zoom controls
   - Filter options
   - Export buttons
4. **Results Panel** (Bottom):
   - Summary statistics
   - Gap analysis table
   - Overlap warnings

### Color Scheme
- **Primary**: Blue (#2563eb)
- **Success/Allocated**: Green (#16a34a)
- **Warning/Overlaps**: Orange (#ea580c)
- **Error/Gaps**: Red (#dc2626)
- **Background**: Light gray (#f8fafc) / Dark gray (#1e293b)

## Data Models

### CIDR Range
```go
type CIDRRange struct {
    Original    string    `json:"original"`
    Network     net.IP    `json:"network"`
    Mask        net.IPMask `json:"mask"`
    Broadcast   net.IP    `json:"broadcast"`
    TotalIPs    int       `json:"total_ips"`
    UsableIPs   int       `json:"usable_ips"`
    Valid       bool      `json:"valid"`
}
```

### Gap Analysis
```go
type Gap struct {
    StartIP     net.IP `json:"start_ip"`
    EndIP       net.IP `json:"end_ip"`
    Size        int    `json:"size"`
    Suggested   string `json:"suggested_cidr"`
}
```

### Overlap Detection
```go
type Overlap struct {
    CIDR1       string `json:"cidr1"`
    CIDR2       string `json:"cidr2"`
    Intersection string `json:"intersection"`
    Type        string `json:"type"` // "partial", "complete"
}
```

## Implementation Phases

### Phase 1: Core Functionality (MVP)
- Basic CIDR input and validation
- Simple gap detection
- Basic visualization (table format)
- REST API with core endpoints

### Phase 2: Enhanced Visualization
- Interactive network diagrams
- Color-coded visualization
- Zoom and pan functionality
- Export capabilities

### Phase 3: Advanced Features
- File import/export
- Advanced filtering
- Overlap detection and resolution
- Performance optimizations

### Phase 4: Polish and Optimization
- UI/UX improvements
- Accessibility features
- Comprehensive testing
- Documentation

## Testing Requirements

### Unit Tests
- CIDR parsing and validation functions
- Gap detection algorithms
- IP address calculation utilities
- API endpoint logic

### Integration Tests
- Frontend-backend API communication
- File import/export functionality
- Visualization rendering

### User Acceptance Tests
- End-to-end user workflows
- Cross-browser compatibility
- Performance under load
- Accessibility compliance

## Deployment Requirements

### Development Environment
- Docker Compose for local development
- Hot reload for both frontend and backend
- Environment variable configuration

### Production Environment
- Container deployment (Docker)
- Reverse proxy configuration (Nginx)
- HTTPS termination
- Health check endpoints
- Logging and monitoring

## Success Criteria

1. **Functional**: All core features implemented and working correctly
2. **Performance**: Handles 1000+ CIDR ranges with sub-second response times
3. **Usability**: Users can complete gap analysis within 2 minutes of first use
4. **Reliability**: 99.9% uptime with proper error handling
5. **Maintainability**: Clean, documented code with comprehensive test coverage

## Future Enhancements

1. IPv6 support
2. Historical analysis and trending
3. Integration with network management tools
4. Multi-user collaboration features
5. Advanced reporting and analytics
6. Mobile application
7. API rate limiting and authentication
8. Database persistence for large datasets