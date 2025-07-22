package utils

import (
	"net/http"
	"strings"
)

// GetIPAddress extracts the real IP address from the request
func GetIPAddress(r *http.Request) string {
	// Check X-Forwarded-For header first (for proxies)
	forwarded := r.Header.Get("X-Forwarded-For")
	if forwarded != "" {
		// X-Forwarded-For can contain multiple IPs, take the first one
		ips := strings.Split(forwarded, ",")
		if len(ips) > 0 {
			return strings.TrimSpace(ips[0])
		}
	}

	// Check X-Real-IP header (for nginx)
	realIP := r.Header.Get("X-Real-IP")
	if realIP != "" {
		return realIP
	}

	// Fall back to RemoteAddr
	ip := r.RemoteAddr
	if strings.Contains(ip, ":") {
		// Remove port if present
		parts := strings.Split(ip, ":")
		if len(parts) > 1 {
			return strings.Join(parts[:len(parts)-1], ":")
		}
	}
	return ip
}

// GetUserAgent extracts the user agent from the request
func GetUserAgent(r *http.Request) string {
	userAgent := r.Header.Get("User-Agent")
	if len(userAgent) > 500 {
		// Truncate if too long
		return userAgent[:500]
	}
	return userAgent
}