package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
)

type Response struct {
	IP       string `json:"ip"`
	Hostname string `json:"hostname"`
}

// This map is used to store the mapping between IP addresses and hostnames.
// This would be replaced with a lookup function to something like infoblox.
var ipHostnameMap = map[string]string{
	"127.0.0.1": "localhost",
	"192.168.1.2": "host2",
	// ... add more IPs and hostnames as needed ...
}

func ReadUserIP(r *http.Request) string {
	IPAddress := r.Header.Get("X-Real-Ip")
	if IPAddress == "" {
		IPAddress = r.Header.Get("X-Forwarded-For")
	}
	if IPAddress == "" {
		IPAddress = r.RemoteAddr
	}
	return strings.Split(IPAddress, ":")[0]
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
	log.Println("Starting server on :8080")
	http.HandleFunc("/", handler)
	if err := http.ListenAndServe(":8081", nil); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
