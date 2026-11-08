// analytics – sends pipeline and business events to a configurable endpoint.
//
// Usage:
//
//	go run ./cmd/analytics --type technical --payload '{"job":"build","status":"success"}'
//	go run ./cmd/analytics --type business  --payload '{"event":"deploy","variant":"blog"}'
//
// Environment:
//
//	ANALYTICS_ENDPOINT – HTTP endpoint to POST events to (optional, defaults to stdout)
//	ANALYTICS_TOKEN    – Bearer token for authentication (optional)
package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"os"
	"time"
)

func main() {
	typ := flag.String("type", "", "Event type: technical or business")
	payloadRaw := flag.String("payload", "", "JSON payload string")
	endpoint := flag.String("endpoint", "", "HTTP endpoint to POST events to")
	token := flag.String("token", "", "Bearer token for authentication")
	flag.Parse()

	if *typ == "" || *payloadRaw == "" {
		fmt.Fprintln(os.Stderr, "Usage: analytics --type <technical|business> --payload '<json>' [--endpoint <url>] [--token <bearer>]")
		os.Exit(2)
	}

	validTypes := map[string]bool{"technical": true, "business": true}
	if !validTypes[*typ] {
		fmt.Fprintf(os.Stderr, "Invalid --type: %q. Must be technical or business.\n", *typ)
		os.Exit(2)
	}

	// Override from env if flags not set
	if *endpoint == "" {
		*endpoint = os.Getenv("ANALYTICS_ENDPOINT")
	}
	if *token == "" {
		*token = os.Getenv("ANALYTICS_TOKEN")
	}

	var payload map[string]interface{}
	if err := json.Unmarshal([]byte(*payloadRaw), &payload); err != nil {
		fmt.Fprintf(os.Stderr, "Invalid --payload JSON: %v\n", err)
		os.Exit(2)
	}

	host, _ := os.Hostname()
	doc := map[string]interface{}{
		"type":      *typ,
		"payload":   payload,
		"timestamp": time.Now().UTC().Format(time.RFC3339),
		"hostname":  host,
	}

	body, err := json.Marshal(doc)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to marshal event: %v\n", err)
		os.Exit(1)
	}

	if *endpoint != "" {
		client := &http.Client{Timeout: 10 * time.Second}
		if err := sendEvent(*endpoint, *token, body, client); err != nil {
			fmt.Fprintf(os.Stderr, "Failed to send event: %v\n", err)
			os.Exit(1)
		}

		fmt.Println("Analytics event recorded successfully")
	} else {
		// No endpoint configured: log to stdout
		fmt.Printf("[analytics] %s\n", string(body))
	}
}

func sendEvent(endpoint, token string, body []byte, client *http.Client) error {
	req, err := http.NewRequest(http.MethodPost, endpoint, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= http.StatusBadRequest {
		return fmt.Errorf("server returned %d", resp.StatusCode)
	}
	return nil
}
