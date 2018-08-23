package http

import "net/http"

// GetRemoteIP of http Request
func GetRemoteIP(r *http.Request) string {
	ip := r.Header.Get("X-Forwarded-For")
	if len(ip) <= 1 {
		ip = r.RemoteAddr
	}
	return ip
}
