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
	IP              string `json:"ip"`
	Hostname        string `json:"hostname"`
	ReverseIPLookup string `json:"reverseiplookup"`
	FQDN2IP         string `json:"fqdn2ip"`
}

func ReadUserIP(httprequest *http.Request) string {
	// Check X-Forwarded-For header
	if ip := httprequest.Header.Get("X-Forwarded-For"); ip != "" {
		ips := strings.Split(ip, ",")
		return strings.TrimSpace(ips[0])
	}
	// Check X-Real-Ip header
	if ip := httprequest.Header.Get("X-Real-Ip"); ip != "" {
		return strings.TrimSpace(ip)
	}
	// Fallback to RemoteAddr
	ip, _, err := net.SplitHostPort(httprequest.RemoteAddr)
	if err != nil {
		return httprequest.RemoteAddr
	}
	return ip
}

func handler(httpwriter http.ResponseWriter, httprequest *http.Request) {
	clientIP := ReadUserIP(httprequest)
	log.Printf("Client IP: %s", clientIP)

	// Reverse DNS lookup to get hostname
	hostnames, err := net.LookupAddr(clientIP)
	var hostname string
	if err == nil && len(hostnames) > 0 {
		hostname = strings.TrimSuffix(hostnames[0], ".")
	} else {
		log.Printf("Reverse lookup failed for IP %s: %v", clientIP, err)
		hostname = "unknown"
	}

	// Forward DNS lookup from hostname to get IP
	ipAddresses, err := net.LookupHost(hostname)
	var fqdn2ip string
	if err == nil && len(ipAddresses) > 0 {
		fqdn2ip = ipAddresses[0]
	} else {
		log.Printf("Forward lookup failed for hostname %s: %v", hostname, err)
		fqdn2ip = "unknown"
	}

	response := Response{
		IP:              clientIP,
		Hostname:        hostname,
		ReverseIPLookup: hostname,
		FQDN2IP:         fqdn2ip,
	}

	httpwriter.Header().Set("Content-Type", "application/json")

	prettyJSON, err := json.MarshalIndent(response, "", "  ")
	if err != nil {
		http.Error(httpwriter, err.Error(), http.StatusInternalServerError)
		return
	}
	httpwriter.Write(prettyJSON)
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	address := ":" + port
	log.Printf("Starting server on %s", address)
	http.HandleFunc("/", handler)
	if err := http.ListenAndServe(address, nil); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
