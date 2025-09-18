package handlers

import (
	"fmt"
	"net/http"
	"time"

	"cidr-viewer/models"
	"cidr-viewer/utils"

	"github.com/gin-gonic/gin"
)

// AnalyzeCIDRs handles CIDR analysis requests
func AnalyzeCIDRs(c *gin.Context) {
	var req models.AnalysisRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	}

	// Debug logging
	fmt.Printf("=== BACKEND RECEIVED REQUEST ===\n")
	fmt.Printf("CIDRs: %v\n", req.CIDRs)
	fmt.Printf("VPCCIDRs: %v\n", req.VPCCIDRs)
	fmt.Printf("SubnetCIDRs: %v\n", req.SubnetCIDRs)
	fmt.Printf("================================\n")

	var validCIDRs []models.CIDRRange
	var invalidCIDRs []models.CIDRRange

	// Parse and validate VPC CIDRs
	for _, cidrStr := range req.VPCCIDRs {
		cidr := utils.ParseCIDR(cidrStr)
		cidr.Category = "vpc"
		if cidr.Valid {
			validCIDRs = append(validCIDRs, cidr)
		} else {
			invalidCIDRs = append(invalidCIDRs, cidr)
		}
	}

	// Parse and validate Subnet CIDRs
	for _, cidrStr := range req.SubnetCIDRs {
		cidr := utils.ParseCIDR(cidrStr)
		cidr.Category = "subnet"
		if cidr.Valid {
			validCIDRs = append(validCIDRs, cidr)
		} else {
			invalidCIDRs = append(invalidCIDRs, cidr)
		}
	}

	// Also handle legacy format (all CIDRs in single array)
	for _, cidrStr := range req.CIDRs {
		// Skip if already processed in VPC or Subnet arrays
		alreadyProcessed := false
		for _, vpcCidr := range req.VPCCIDRs {
			if vpcCidr == cidrStr {
				alreadyProcessed = true
				break
			}
		}
		if !alreadyProcessed {
			for _, subnetCidr := range req.SubnetCIDRs {
				if subnetCidr == cidrStr {
					alreadyProcessed = true
					break
				}
			}
		}
		
		if !alreadyProcessed {
			cidr := utils.ParseCIDR(cidrStr)
			if cidr.Valid {
				validCIDRs = append(validCIDRs, cidr)
			} else {
				invalidCIDRs = append(invalidCIDRs, cidr)
			}
		}
	}

	// Find gaps and overlaps
	gaps := utils.FindGaps(validCIDRs)
	overlaps := utils.FindOverlaps(validCIDRs)
	summary := utils.CalculateSummary(validCIDRs, gaps, overlaps)

	response := models.AnalysisResponse{
		ValidCIDRs:   validCIDRs,
		InvalidCIDRs: invalidCIDRs,
		Gaps:         gaps,
		Overlaps:     overlaps,
		Summary:      summary,
	}

	c.JSON(http.StatusOK, response)
}

// ValidateCIDR handles single CIDR validation requests
func ValidateCIDR(c *gin.Context) {
	var req models.ValidationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	}

	cidr := utils.ParseCIDR(req.CIDR)
	c.JSON(http.StatusOK, cidr)
}

// HealthCheck handles health check requests
func HealthCheck(c *gin.Context) {
	response := models.HealthResponse{
		Status:    "healthy",
		Timestamp: time.Now().UTC().Format(time.RFC3339),
	}
	c.JSON(http.StatusOK, response)
}