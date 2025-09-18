package utils

import (
	"fmt"
	"net"
	"sort"

	"cidr-viewer/models"
)

// ParseCIDR validates and parses a CIDR string
func ParseCIDR(cidrStr string) models.CIDRRange {
	cidr := models.CIDRRange{
		Original: cidrStr,
		Valid:    false,
	}

	// Parse the CIDR
	_, ipNet, err := net.ParseCIDR(cidrStr)
	if err != nil {
		cidr.ErrorMsg = fmt.Sprintf("Invalid CIDR format: %v", err)
		return cidr
	}

	// Calculate network details
	cidr.Network = ipNet.IP.String()
	cidr.Mask = net.IP(ipNet.Mask).String()
	
	// Calculate broadcast address
	broadcast := make(net.IP, len(ipNet.IP))
	copy(broadcast, ipNet.IP)
	for i := range ipNet.Mask {
		broadcast[i] |= ^ipNet.Mask[i]
	}
	cidr.Broadcast = broadcast.String()

	// Calculate IP counts
	ones, bits := ipNet.Mask.Size()
	if ones == 0 || bits == 0 {
		cidr.ErrorMsg = "Invalid subnet mask"
		return cidr
	}

	cidr.TotalIPs = 1 << (bits - ones)
	cidr.UsableIPs = cidr.TotalIPs - 2 // Subtract network and broadcast
	if cidr.TotalIPs <= 2 {
		cidr.UsableIPs = 0 // /31 and /32 have special rules
	}

	cidr.Valid = true
	return cidr
}

// IPToInt converts an IP address to integer for calculations
func IPToInt(ip net.IP) uint32 {
	ip = ip.To4()
	if ip == nil {
		return 0
	}
	return uint32(ip[0])<<24 + uint32(ip[1])<<16 + uint32(ip[2])<<8 + uint32(ip[3])
}

// IntToIP converts an integer back to IP address
func IntToIP(i uint32) net.IP {
	return net.IPv4(byte(i>>24), byte(i>>16), byte(i>>8), byte(i))
}

// FindGaps identifies gaps between CIDR ranges
func FindGaps(validCIDRs []models.CIDRRange) []models.Gap {
	if len(validCIDRs) == 0 {
		return []models.Gap{}
	}

	// Create a list of IP ranges
	type ipRange struct {
		start uint32
		end   uint32
	}

	var ranges []ipRange
	for _, cidr := range validCIDRs {
		_, ipNet, err := net.ParseCIDR(cidr.Original)
		if err != nil {
			continue
		}
		
		start := IPToInt(ipNet.IP)
		broadcast := make(net.IP, len(ipNet.IP))
		copy(broadcast, ipNet.IP)
		for i := range ipNet.Mask {
			broadcast[i] |= ^ipNet.Mask[i]
		}
		end := IPToInt(broadcast)
		
		ranges = append(ranges, ipRange{start: start, end: end})
	}

	// Sort ranges by start IP
	sort.Slice(ranges, func(i, j int) bool {
		return ranges[i].start < ranges[j].start
	})

	// Merge overlapping ranges
	merged := []ipRange{}
	for _, r := range ranges {
		if len(merged) == 0 || merged[len(merged)-1].end < r.start-1 {
			merged = append(merged, r)
		} else {
			if merged[len(merged)-1].end < r.end {
				merged[len(merged)-1].end = r.end
			}
		}
	}

	// Find gaps
	var gaps []models.Gap
	for i := 0; i < len(merged)-1; i++ {
		gapStart := merged[i].end + 1
		gapEnd := merged[i+1].start - 1
		
		if gapStart <= gapEnd {
			size := int(gapEnd - gapStart + 1)
			gap := models.Gap{
				StartIP:      IntToIP(gapStart).String(),
				EndIP:        IntToIP(gapEnd).String(),
				Size:         size,
				SuggestedCIDR: SuggestCIDR(gapStart, gapEnd),
			}
			gaps = append(gaps, gap)
		}
	}

	return gaps
}

// SuggestCIDR suggests the best CIDR notation for a gap
func SuggestCIDR(start, end uint32) string {
	size := end - start + 1
	
	// Find the largest power of 2 that fits
	bits := 32
	for i := 0; i < 32; i++ {
		if (1 << i) >= int(size) {
			bits = 32 - i
			break
		}
	}
	
	// Align to proper boundary
	mask := uint32(0xFFFFFFFF) << (32 - bits)
	alignedStart := start & mask
	
	return fmt.Sprintf("%s/%d", IntToIP(alignedStart).String(), bits)
}

// FindOverlaps detects overlapping CIDR ranges
func FindOverlaps(validCIDRs []models.CIDRRange) []models.Overlap {
	var overlaps []models.Overlap
	
	// Early return for small datasets
	if len(validCIDRs) < 2 {
		return overlaps
	}
	
	// For large datasets, we might want to limit the comparison
	// to prevent performance issues
	maxComparisons := 1000
	comparisons := 0
	
	for i := 0; i < len(validCIDRs) && comparisons < maxComparisons; i++ {
		for j := i + 1; j < len(validCIDRs) && comparisons < maxComparisons; j++ {
			comparisons++
			overlap := checkOverlap(validCIDRs[i], validCIDRs[j])
			if overlap.CIDR1 != "" {
				overlaps = append(overlaps, overlap)
			}
		}
	}
	
	return overlaps
}

// checkOverlap checks if two CIDR ranges overlap
func checkOverlap(cidr1, cidr2 models.CIDRRange) models.Overlap {
	_, net1, err1 := net.ParseCIDR(cidr1.Original)
	_, net2, err2 := net.ParseCIDR(cidr2.Original)
	
	if err1 != nil || err2 != nil {
		return models.Overlap{}
	}
	
	// Check if networks overlap
	if net1.Contains(net2.IP) || net2.Contains(net1.IP) {
		overlapType := "partial"
		intersection := cidr1.Original
		
		// Determine overlap type and intersection
		if net1.Contains(net2.IP) && net2.Contains(net1.IP) {
			overlapType = "complete"
		} else if net1.Contains(net2.IP) {
			intersection = cidr2.Original
		}
		
		return models.Overlap{
			CIDR1:        cidr1.Original,
			CIDR2:        cidr2.Original,
			Intersection: intersection,
			Type:         overlapType,
		}
	}
	
	return models.Overlap{}
}

// CalculateSummary generates summary statistics
func CalculateSummary(validCIDRs []models.CIDRRange, gaps []models.Gap, overlaps []models.Overlap) models.Summary {
	totalIPs := 0
	allocatedIPs := 0
	
	for _, cidr := range validCIDRs {
		allocatedIPs += cidr.TotalIPs
	}
	
	availableIPs := 0
	for _, gap := range gaps {
		availableIPs += gap.Size
	}
	
	totalIPs = allocatedIPs + availableIPs
	
	return models.Summary{
		TotalIPs:     totalIPs,
		AllocatedIPs: allocatedIPs,
		AvailableIPs: availableIPs,
		GapCount:     len(gaps),
		OverlapCount: len(overlaps),
	}
}