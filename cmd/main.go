package main

import (
    "context"
    "encoding/json"
    "fmt"
    "log"
    "os"
    "os/signal"
    "syscall"
    "time"

    "libvirt.org/go/libvirt"
)

// Config holds settings for the collector
type Config struct {
    ConnectURI string
    Interval   time.Duration
}

// MetricEvent represents the basic structure of a single event for Elastic Agent.
type MetricEvent struct {
    Timestamp time.Time `json:"@timestamp"`
    ECS       struct {
        Version string `json:"version"`
    } `json:"ecs"`
    Agent struct {
        ID string `json:"id,omitempty"`
    } `json:"agent"`

    // Custom fields for libvirt.domain_metrics
    Libvirt struct {
        Domain struct {
            Name  string `json:"name"`
            UUID  string `json:"uuid"`
            State string `json:"state"`
            ID    uint   `json:"id"`
        } `json:"domain"`
    } `json:"libvirt"`
}

func main() {
    // Configure logging to use stderr for all control/error output
    log.SetOutput(os.Stderr)
    log.SetPrefix("[libvirt_collector] ")

    // Load Configuration
	cfg := Config{
		ConnectURI: os.Getenv("LIBVIRT_CONNECT_URI"),
		Interval:   30 * time.Second, // Default 30s interval
	}
	if cfg.ConnectURI == "" { // FALLBACK for manual testing if not run by Elastic Agent
        // cfg.ConnectURI = "qemu://system"
		cfg.ConnectURI = "qemu+tcp://192.168.0.238/system"
		log.Printf("Info: LIBVIRT_CONNECT_URI environment variable not set. Defaulting to %s", cfg.ConnectURI)
	}
    log.Printf("Collector initialized with URI: %s, Interval: %v", cfg.ConnectURI, cfg.Interval)

    // Setup Context for Graceful Shutdown
    ctx, cancel := context.WithCancel(context.Background())
    defer cancel()

    // Handle OS Signals (SIGINT, SIGTERM)
    signalCh := make(chan os.Signal, 1)
    signal.Notify(signalCh, syscall.SIGINT, syscall.SIGTERM)

    // Start the main collector loop in a goroutine
    go runCollector(ctx, cfg)

    // Block until a signal is received or the context is cancelled
    select {
    case sig := <-signalCh:
        log.Printf("Received signal: %v. Shutting down...", sig)
    case <-ctx.Done():
        log.Println("Context cancelled. Shutting down...")
    }

    // Graceful shutdown delay
    log.Println("Shutdown complete.")
}

// runCollector is the main metric collection loop.
func runCollector(ctx context.Context, cfg Config) {
    // Initial Libvirt Connection
    conn, err := libvirt.NewConnect(cfg.ConnectURI)
    if err != nil {
        log.Fatalf("Fatal: Failed to connect to libvirt at %s: %v", cfg.ConnectURI, err)
        return
    }
    defer conn.Close()
    log.Println("Successfully connected to Libvirt.")

    // Start Ticker
    ticker := time.NewTicker(cfg.Interval)
    defer ticker.Stop()

    // Run the collection immediately upon start
    collectAndPublishMetrics(conn)

    // Main Loop
    for {
        select {
        case <-ctx.Done():
            log.Println("Collector loop stopped.")
            return // Exit the loop and the goroutine
        case <-ticker.C:
            collectAndPublishMetrics(conn)
        }
    }
}

// collectAndPublishMetrics performs the actual data collection and output.
func collectAndPublishMetrics(conn *libvirt.Connect) {
    // Get Domains
    domains, err := conn.ListAllDomains(0)
    
    if err != nil {
        log.Printf("Error listing domains: %v", err)
        return
    }

    for _, dom := range domains {
        defer dom.Free()

        // Get Domain Info
        info, err := dom.GetInfo()
        if err != nil {
            name, _ := dom.GetName()
            log.Printf("Error getting info for domain %s: %v", name, err)
            continue
        }

        // Pass the raw integer state value to the mapper
        stateStr := mapState(info.State) 

        // Create Metric Event
        name, _ := dom.GetName()
        uuidBytes, _ := dom.GetUUID() 
        id, _ := dom.GetID()

        event := MetricEvent{
            Timestamp: time.Now().UTC(),
            ECS:       struct{ Version string `json:"version"` }{Version: "8.11.0"},
            Libvirt: struct {
                Domain struct {
                    Name  string `json:"name"`
                    UUID  string `json:"uuid"`
                    State string `json:"state"`
                    ID    uint   `json:"id"`
                } `json:"domain"`
            }{
                Domain: struct {
                    Name  string `json:"name"`
                    UUID  string `json:"uuid"`
                    State string `json:"state"`
                    ID    uint   `json:"id"`
                }{
                    Name:  name,
                    UUID:  fmt.Sprintf("%x", uuidBytes),
                    State: stateStr,
                    ID:    uint(id), 
                },
            },
        }

        // Marshal to JSON and Print to STDOUT
        jsonEvent, err := json.Marshal(event)
        if err != nil {
            log.Printf("Error marshalling event for domain %s: %v", name, err)
            continue
        }

        fmt.Println(string(jsonEvent))
    }   
}

// Converts the raw libvirt DomainState integer value to a descriptive string.
func mapState(state libvirt.DomainState) string {
    switch int(state) {
    case 0:
        return "no_state"
    case 1:
        return "running"
    case 2:
        return "blocked"
    case 3:
        return "paused"
    case 4:
        return "shutdown"
    case 5:
        return "shutoff"
    case 6:
        return "crashed"
    case 7:
        return "pmsuspended"
    default:
        return "unknown"
    }
}