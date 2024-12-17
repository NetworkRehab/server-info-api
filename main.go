package main

import (
	"encoding/json"
	"log"
	"net"
	"net/http"
	"os"
	"strings"
)

type Response struct {
	IP       string `json:"ip"`
	Hostname string `json:"hostname"`
}

// This map is used to store the mapping between IP addresses and hostnames.
// This would be replaced with a lookup function to something like infoblox.
var ipHostnameMap = map[string]string{
	"127.0.0.1":   "localhost",
	"192.168.1.2": "host2",
	// ... add more IPs and hostnames as needed ...
}

func ReadUserIP(r *http.Request) string {
	IPAddress := r.Header.Get("X-Real-Ip")
	if IPAddress == "" {
		IPAddress = r.Header.Get("X-Forwarded-For")
		if IPAddress != "" {
			// X-Forwarded-For may contain multiple IPs, take the first one
			parts := strings.Split(IPAddress, ",")
			IPAddress = strings.TrimSpace(parts[0])
		}
	}
	if IPAddress == "" {
		ip, _, err := net.SplitHostPort(r.RemoteAddr)
		if err != nil {
			// Fallback to r.RemoteAddr if SplitHostPort fails
			ip = r.RemoteAddr
		}
		IPAddress = ip
	}
	return IPAddress
}

func handler(w http.ResponseWriter, r *http.Request) {
	ip := ReadUserIP(r)
	log.Printf("Client IP: %s", ip)
	hostname, exists := ipHostnameMap[ip]
	if !exists {
		hostname = "unknown"
	}
	response := Response{IP: ip, Hostname: hostname}
	w.Header().Set("Content-Type", "application/json")

	prettyJSON, err := json.MarshalIndent(response, "", "  ")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(prettyJSON)
}

func main() {
	port := ":8081"
	if p := os.Getenv("PORT"); p != "" {
		port = ":" + p
	}
	log.Printf("Starting server on %s", port)
	http.HandleFunc("/", handler)
	if err := http.ListenAndServe(port, nil); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
