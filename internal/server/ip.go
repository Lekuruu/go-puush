package server

import (
	"net"
	"net/http"
	"strings"
)

// Cloudflare specific header for original client IP
const cloudflareIPHeader = "CF-Connecting-IP"

// Common reverse proxy headers
var proxyIPHeaders = [...]string{
	"X-Forwarded-For",
	"X-Real-IP",
	"Forwarded",
}

// GetRequestIP tries to extract the real client IP address from the request headers
func GetRequestIP(request *http.Request) string {
	if ip := request.Header.Get(cloudflareIPHeader); ip != "" {
		return ip
	}

	for _, header := range proxyIPHeaders {
		if value := request.Header.Get(header); value != "" {
			// X-Forwarded-For may be a comma-separated list -> take the first entry
			ip := strings.TrimSpace(strings.Split(value, ",")[0])
			if ip != "" {
				return ip
			}
		}
	}

	// Fallback to remote address
	host, _, err := net.SplitHostPort(request.RemoteAddr)
	if err != nil {
		return request.RemoteAddr
	}
	return host
}
