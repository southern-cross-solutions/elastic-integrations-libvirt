package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"libvirt.org/go/libvirt"
)

type DomainInfo struct {
	Name     string `json:"name"`
	UUID     string `json:"uuid"`
	State    string `json:"state"`
	VCPU     uint   `json:"vcpu"`
	MemoryKB uint64 `json:"memory_kb"`
}

// Event format for Elastic Agent ingestion
type VMEvent struct {
	Timestamp string       `json:"@timestamp"`
	Libvirt   LibvirtField `json:"libvirt"`
}

type LibvirtField struct {
	VM VMField `json:"vm"`
}

type VMField struct {
	UUID   string     `json:"uuid"`
	Name   string     `json:"name"`
	State  string     `json:"state"`
	VCPU   uint       `json:"vcpu"`
	Memory MemoryData `json:"memory"`
}

type MemoryData struct {
	KB uint64 `json:"kb"`
}

func getDomains() ([]DomainInfo, error) {
	conn, err := libvirt.NewConnect("qemu:///system")
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	domains, err := conn.ListAllDomains(0)
	if err != nil {
		return nil, err
	}

	var result []DomainInfo
    log.Printf("Scraped %d VMs\n", len(domains))
	for _, dom := range domains {
        log.Printf("%+v\n", dom)
		name, _ := dom.GetName()
		uuid, _ := dom.GetUUIDString()
		info, _ := dom.GetInfo()

		result = append(result, DomainInfo{
			Name:     name,
			UUID:     uuid,
			State:    domainStateToString(info.State),
			VCPU:     info.NrVirtCpu,
			MemoryKB: info.Memory,
		})
	}
	return result, nil
}

func domainStateToString(state libvirt.DomainState) string {
	switch state {
	case libvirt.DOMAIN_RUNNING:
		return "running"
	case libvirt.DOMAIN_PAUSED:
		return "paused"
	case libvirt.DOMAIN_SHUTDOWN:
		return "shutdown"
	case libvirt.DOMAIN_SHUTOFF:
		return "shutoff"
	case libvirt.DOMAIN_CRASHED:
		return "crashed"
	default:
		return "unknown"
	}
}

// Convert domains to ECS-compatible events
func domainsToEvents(domains []DomainInfo) []VMEvent {
	events := make([]VMEvent, 0, len(domains))
	now := time.Now().UTC().Format(time.RFC3339)

	for _, d := range domains {
		events = append(events, VMEvent{
			Timestamp: now,
			Libvirt: LibvirtField{
				VM: VMField{
					UUID:  d.UUID,
					Name:  d.Name,
					State: d.State,
					VCPU:  d.VCPU,
					Memory: MemoryData{
						KB: d.MemoryKB,
					},
				},
			},
		})
	}

	return events
}

func main() {
	http.HandleFunc("/v1/domains", func(w http.ResponseWriter, r *http.Request) {
		domains, err := getDomains()
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}

		events := domainsToEvents(domains)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(events)
	})

	log.Println("Serving JSON Elastic Agent endpoint on http://0.0.0.0:8088/v1/domains")
	log.Fatal(http.ListenAndServe(":8088", nil))
}















// package main

// import (
//     "encoding/json"
//     "log"
//     "os"
//     // "time"

//     "libvirt.org/go/libvirt"
// )

// type DomainInfo struct {
//     Name     string `json:"name"`
//     UUID     string `json:"uuid"`
//     State    string `json:"state"`
//     VCPU     uint   `json:"vcpu"`
//     MemoryKB uint64 `json:"memory_kb"`
// }

// type DomainsResponse struct {
//     Domains []DomainInfo `json:"domains"`
// }

// func getDomains() ([]DomainInfo, error) {
//     conn, err := libvirt.NewConnect("qemu:///system")
//     if err != nil {
//         return nil, err
//     }
//     defer conn.Close()

//     domains, err := conn.ListAllDomains(0)
//     if err != nil {
//         return nil, err
//     }

//     var result []DomainInfo
//     for _, dom := range domains {
//         name, _ := dom.GetName()
//         uuid, _ := dom.GetUUIDString()
//         info, _ := dom.GetInfo()

//         result = append(result, DomainInfo{
//             Name:     name,
//             UUID:     uuid,
//             State:    domainStateToString(info.State),
//             VCPU:     info.NrVirtCpu,
//             MemoryKB: info.Memory,
//         })
//     }
//     return result, nil
// }

// func domainStateToString(state libvirt.DomainState) string {
//     switch state {
//     case libvirt.DOMAIN_RUNNING:
//         return "running"
//     case libvirt.DOMAIN_PAUSED:
//         return "paused"
//     case libvirt.DOMAIN_SHUTDOWN:
//         return "shutdown"
//     case libvirt.DOMAIN_SHUTOFF:
//         return "shutoff"
//     case libvirt.DOMAIN_CRASHED:
//         return "crashed"
//     default:
//         return "unknown"
//     }
// }

// func main() {
//     logFile, _ := os.OpenFile("/var/log/libvirt-exporter.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
//     logger := log.New(logFile, "", log.LstdFlags)

//     // for {
//         domains, err := getDomains()
//         if err != nil {
//             logger.Printf("Error fetching domains: %v\n", err)
//         } else {
//             logger.Printf("Fetched %d domains\n", len(domains))
//             for _, d := range domains {
//                 logger.Printf("Domain: %+v\n", d)
//             }

//             // Print JSON to stdout for Elastic Agent
//             resp := DomainsResponse{Domains: domains}
//             enc := json.NewEncoder(os.Stdout)
//             enc.SetIndent("", "  ") // optional pretty-print
//             enc.Encode(resp)
//         }

//     //     time.Sleep(30 * time.Second) // optional interval for local testing
//     // }
// }
