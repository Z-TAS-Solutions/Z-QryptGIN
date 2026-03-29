package utils

import (
	"strings"

	"github.com/gin-gonic/gin"
)

// ParseDeviceName extracts a human-readable device name from a User-Agent string.
func ParseDeviceName(ua string) string {
	if ua == "" {
		return "Unknown Device"
	}

	// Default to full UA if we can't parse, truncated
	fallback := ua
	if len(fallback) > 50 {
		fallback = fallback[:47] + "..."
	}

	var os, browser string

	// Detect OS
	if strings.Contains(ua, "Windows NT 10.0") {
		os = "Windows 10/11"
	} else if strings.Contains(ua, "Windows NT") {
		os = "Windows"
	} else if strings.Contains(ua, "Mac OS X") {
		os = "macOS"
	} else if strings.Contains(ua, "Android") {
		os = "Android"
	} else if strings.Contains(ua, "iPhone") || strings.Contains(ua, "iPad") {
		os = "iOS"
	} else if strings.Contains(ua, "Linux") {
		os = "Linux"
	}

	// Detect Browser
	if strings.Contains(ua, "Edg/") || strings.Contains(ua, "Edge/") {
		browser = "Edge"
	} else if strings.Contains(ua, "OPR/") || strings.Contains(ua, "Opera/") {
		browser = "Opera"
	} else if strings.Contains(ua, "Chrome/") {
		browser = "Chrome"
	} else if strings.Contains(ua, "Firefox/") {
		browser = "Firefox"
	} else if strings.Contains(ua, "Safari/") && !strings.Contains(ua, "Chrome/") {
		browser = "Safari"
	}

	if os != "" && browser != "" {
		return browser + " on " + os
	} else if os != "" {
		return os + " Device"
	} else if browser != "" {
		return browser
	}

	return fallback
}

// GetClientIP gets a clean IP address from the Gin context
func GetClientIP(c *gin.Context) string {
	ip := c.ClientIP()
	if ip == "::1" || ip == "127.0.0.1" {
		return "localhost"
	}
	return ip
}
