package models

// CIDRRange represents a CIDR network range with analysis data
type CIDRRange struct {
	Original    string `json:"original"`
	Network     string `json:"network"`
	Mask        string `json:"mask"`
	Broadcast   string `json:"broadcast"`
	TotalIPs    int    `json:"total_ips"`
	UsableIPs   int    `json:"usable_ips"`
	Valid       bool   `json:"valid"`
	Category    string `json:"category,omitempty"` // "vpc", "subnet"
	ErrorMsg    string `json:"error_msg,omitempty"`
}

// Gap represents an unused IP range between CIDR blocks
type Gap struct {
	StartIP      string `json:"start_ip"`
	EndIP        string `json:"end_ip"`
	Size         int    `json:"size"`
	SuggestedCIDR string `json:"suggested_cidr"`
}

// Overlap represents overlapping CIDR ranges
type Overlap struct {
	CIDR1        string `json:"cidr1"`
	CIDR2        string `json:"cidr2"`
	Intersection string `json:"intersection"`
	Type         string `json:"type"` // "partial", "complete"
}

// AnalysisRequest represents the request payload for CIDR analysis
type AnalysisRequest struct {
	CIDRs       []string `json:"cidrs"`
	VPCCIDRs    []string `json:"vpc_cidrs,omitempty"`
	SubnetCIDRs []string `json:"subnet_cidrs,omitempty"`
}

// AnalysisResponse represents the response from CIDR analysis
type AnalysisResponse struct {
	ValidCIDRs   []CIDRRange `json:"valid_cidrs"`
	InvalidCIDRs []CIDRRange `json:"invalid_cidrs"`
	Gaps         []Gap       `json:"gaps"`
	Overlaps     []Overlap   `json:"overlaps"`
	Summary      Summary     `json:"summary"`
}

// Summary provides overall statistics
type Summary struct {
	TotalIPs     int `json:"total_ips"`
	AllocatedIPs int `json:"allocated_ips"`
	AvailableIPs int `json:"available_ips"`
	GapCount     int `json:"gap_count"`
	OverlapCount int `json:"overlap_count"`
}

// ValidationRequest represents a single CIDR validation request
type ValidationRequest struct {
	CIDR string `json:"cidr"`
}

// HealthResponse represents the health check response
type HealthResponse struct {
	Status    string `json:"status"`
	Timestamp string `json:"timestamp"`
}