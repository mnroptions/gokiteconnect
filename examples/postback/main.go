package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	kiteconnect "github.com/zerodha/gokiteconnect/v4"
)

// Logger for structured logging
type Logger struct {
	file *os.File
}

func NewLogger(filename string) (*Logger, error) {
	file, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return nil, err
	}
	return &Logger{file: file}, nil
}

func (l *Logger) Log(data interface{}) {
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	jsonData, _ := json.MarshalIndent(data, "", "  ")
	
	logEntry := map[string]interface{}{
		"timestamp": timestamp,
		"data":      json.RawMessage(jsonData),
	}
	
	logJSON, _ := json.MarshalIndent(logEntry, "", "  ")
	
	// Log to file
	fmt.Fprintf(l.file, "%s\n", logJSON)
	
	// Log to console
	fmt.Printf("=== Postback Received at %s ===\n", timestamp)
	fmt.Printf("%s\n", jsonData)
	fmt.Printf("=====================================\n\n")
}

func (l *Logger) Close() {
	l.file.Close()
}

// PostbackServer handles incoming postback requests
type PostbackServer struct {
	logger *Logger
	port   string
}

func NewPostbackServer(port string, logFile string) (*PostbackServer, error) {
	logger, err := NewLogger(logFile)
	if err != nil {
		return nil, err
	}
	
	return &PostbackServer{
		logger: logger,
		port:   port,
	}, nil
}

// handlePostback processes incoming postback requests
// Note: Postbacks only contain order updates for orders placed through 
// this specific Kite Connect app (API key/secret), not orders from 
// Kite web/mobile or other Connect apps.
func (ps *PostbackServer) handlePostback(w http.ResponseWriter, r *http.Request) {
	// Only accept POST requests
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST method is allowed", http.StatusMethodNotAllowed)
		return
	}

	// Read the request body
	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Printf("Error reading request body: %v", err)
		http.Error(w, "Error reading request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	// Log request details
	requestInfo := map[string]interface{}{
		"method":        r.Method,
		"url":          r.URL.String(),
		"headers":      r.Header,
		"content_type": r.Header.Get("Content-Type"),
		"user_agent":   r.Header.Get("User-Agent"),
		"remote_addr":  r.RemoteAddr,
		"body_size":    len(body),
	}

	// Try to parse as JSON if body is not empty
	var postbackData interface{}
	if len(body) > 0 {
		// First try to parse as Kite order update
		var order kiteconnect.Order
		if err := json.Unmarshal(body, &order); err == nil {
			postbackData = map[string]interface{}{
				"parsed_as":   "kite_order",
				"order_data":  order,
				"raw_body":    string(body),
			}
		} else {
			// Try to parse as generic JSON
			var genericData interface{}
			if err := json.Unmarshal(body, &genericData); err == nil {
				postbackData = map[string]interface{}{
					"parsed_as":   "generic_json",
					"json_data":   genericData,
					"raw_body":    string(body),
				}
			} else {
				// Store as raw text
				postbackData = map[string]interface{}{
					"parsed_as": "raw_text",
					"raw_body":  string(body),
				}
			}
		}
	} else {
		postbackData = map[string]interface{}{
			"parsed_as": "empty_body",
			"raw_body":  "",
		}
	}

	// Create complete log entry
	logEntry := map[string]interface{}{
		"request_info": requestInfo,
		"postback":     postbackData,
	}

	// Log the postback
	ps.logger.Log(logEntry)

	// Send success response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	response := map[string]string{
		"status":  "success",
		"message": "Postback received and logged successfully",
	}
	json.NewEncoder(w).Encode(response)
}

// handleHealth provides a health check endpoint
func (ps *PostbackServer) handleHealth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	response := map[string]interface{}{
		"status":    "healthy",
		"timestamp": time.Now().Format(time.RFC3339),
		"service":   "kite-postback-server",
	}
	json.NewEncoder(w).Encode(response)
}

// handleRoot provides basic information about the server
func (ps *PostbackServer) handleRoot(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	html := `
<!DOCTYPE html>
<html>
<head>
    <title>Kite Connect Postback Server</title>
    <style>
        body { font-family: Arial, sans-serif; margin: 40px; }
        .header { color: #333; }
        .endpoint { background-color: #f5f5f5; padding: 10px; margin: 10px 0; border-radius: 5px; }
        .method { color: #0066cc; font-weight: bold; }
    </style>
</head>
<body>
    <h1 class="header">Kite Connect Postback Server</h1>
    <p>This server is ready to receive postback notifications from Zerodha Kite Connect API.</p>
    
    <div style="background-color: #e8f4f8; padding: 15px; margin: 15px 0; border-left: 4px solid #2196F3; border-radius: 4px;">
        <strong>Important:</strong> Postbacks only notify about orders placed through your specific Kite Connect app (API key/secret). 
        You will NOT receive notifications for orders placed through Kite web, mobile app, or other Connect apps.
    </div>
    
    <h2>Available Endpoints:</h2>
    <div class="endpoint">
        <span class="method">POST</span> <strong>/postback</strong> - Receives and logs postback notifications
    </div>
    <div class="endpoint">
        <span class="method">GET</span> <strong>/health</strong> - Health check endpoint
    </div>
    <div class="endpoint">
        <span class="method">GET</span> <strong>/</strong> - This information page
    </div>
    
    <h2>Configuration:</h2>
    <p>Configure your Kite Connect app's postback URL to point to: <code>http://your-server-domain:%s/postback</code></p>
    
    <h2>Logs:</h2>
    <p>All postback data is logged to the console and saved to <code>postback_logs.json</code></p>
</body>
</html>`
	
	fmt.Fprintf(w, html, ps.port)
}

// Start starts the postback server
func (ps *PostbackServer) Start() {
	http.HandleFunc("/postback", ps.handlePostback)
	http.HandleFunc("/health", ps.handleHealth)
	http.HandleFunc("/", ps.handleRoot)

	fmt.Printf("🚀 Kite Connect Postback Server starting on port %s\n", ps.port)
	fmt.Printf("📡 Postback endpoint: http://localhost:%s/postback\n", ps.port)
	fmt.Printf("🏥 Health check: http://localhost:%s/health\n", ps.port)
	fmt.Printf("📋 Logs will be saved to: postback_logs.json\n")
	fmt.Printf("🔗 Open http://localhost:%s in your browser for more info\n\n", ps.port)

	log.Fatal(http.ListenAndServe(":"+ps.port, nil))
}

// Close cleans up resources
func (ps *PostbackServer) Close() {
	ps.logger.Close()
}

func main() {
	port := "8080"
	logFile := "postback_logs.json"

	// Override port from environment variable if set
	if envPort := os.Getenv("PORT"); envPort != "" {
		port = envPort
	}

	// Override log file from environment variable if set
	if envLogFile := os.Getenv("LOG_FILE"); envLogFile != "" {
		logFile = envLogFile
	}

	// Create and start the postback server
	server, err := NewPostbackServer(port, logFile)
	if err != nil {
		log.Fatalf("Failed to create postback server: %v", err)
	}
	defer server.Close()

	server.Start()
}