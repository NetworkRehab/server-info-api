package main

import (
	"encoding/json"
	"log"
	"net"
	"net/http"
	"os"
	"strings"
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"flag"
	"encoding/csv"
	"fmt"
)

type Response struct {
	IP              string `json:"ip"`
	Hostname        string `json:"hostname"`
	ReverseIPLookup string `json:"reverseiplookup"`
	FQDN2IP         string `json:"fqdn2ip"`
}

var db *sql.DB

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

func ReverseLookupHostname(ip string) string {
	// Reverse DNS lookup to get hostname
	hostnames, err := net.LookupAddr(ip)
	if err == nil && len(hostnames) > 0 {
		return strings.TrimSuffix(hostnames[0], ".")
	}
	log.Printf("Reverse lookup failed for IP %s: %v", ip, err)
	return "unknown"
}

func ForwardLookupIP(hostname string) string {
	// Forward DNS lookup from hostname to get IP
	ipAddresses, err := net.LookupHost(hostname)
	if err == nil && len(ipAddresses) > 0 {
		return ipAddresses[0]
	}
	log.Printf("Forward lookup failed for hostname %s: %v", hostname, err)
	return "unknown"
}

func initializeDatabase(dbPath string) *sql.DB {
	database, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}

	createTableSQL := `CREATE TABLE IF NOT EXISTS host_ip_map (
		ip TEXT PRIMARY KEY,
		hostname TEXT
	);`

	_, err = database.Exec(createTableSQL)
	if err != nil {
		log.Fatalf("Failed to create table: %v", err)
	}

	return database
}

func getHostnameFromDB(database *sql.DB, ip string) string {
	var hostname string
	query := `SELECT hostname FROM host_ip_map WHERE ip = ?`
	err := database.QueryRow(query, ip).Scan(&hostname)
	if err != nil {
		if err == sql.ErrNoRows {
			hostname = "unknown"
		} else {
			log.Printf("Database query error: %v", err)
			hostname = "unknown"
		}
	}
	return hostname
}

func importCSVToDB(database *sql.DB, csvFilePath string) error {
	file, err := os.Open(csvFilePath)
	if (err != nil) {
		return fmt.Errorf("failed to open CSV file: %v", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	// Read all records
	records, err := reader.ReadAll()
	if err != nil {
		return fmt.Errorf("failed to read CSV file: %v", err)
	}

	tx, err := database.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %v", err)
	}

	stmt, err := tx.Prepare(`INSERT OR REPLACE INTO host_ip_map(ip, hostname) VALUES(?, ?)`)
	if err != nil {
		return fmt.Errorf("failed to prepare statement: %v", err)
	}
	defer stmt.Close()

	for _, record := range records {
		if len(record) < 2 {
			log.Printf("Skipping invalid record: %v", record)
			continue
		}
		ip := strings.TrimSpace(record[0])
		hostname := strings.TrimSpace(record[1])
		_, err = stmt.Exec(ip, hostname)
		if err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to execute statement: %v", err)
		}
	}

	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("failed to commit transaction: %v", err)
	}

	log.Printf("Successfully imported data from %s into the database.", csvFilePath)
	return nil
}

func handler(httpwriter http.ResponseWriter, httprequest *http.Request) {
	clientIP := ReadUserIP(httprequest)
	hostname := getHostnameFromDB(db, clientIP)
	reverseiplookup := ReverseLookupHostname(clientIP)
	fqdn2ip := ForwardLookupIP(hostname)

	log.Printf("Client IP: %s", clientIP)

	response := Response{
		IP:              clientIP,
		Hostname:        hostname,
		ReverseIPLookup: reverseiplookup,
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

	// Add flag for CSV import
	importCSV := flag.String("import", "", "Path to CSV file to import into the database")
	flag.Parse()

	db = initializeDatabase("host_ip_map.db")

	if *importCSV != "" {
		// Import CSV and exit
		err := importCSVToDB(db, *importCSV)
		if err != nil {
			log.Fatalf("Error importing CSV: %v", err)
		}
		return
	}

	http.HandleFunc("/", handler)
	if err := http.ListenAndServe(address, nil); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
