// package main

// import (
// 	"encoding/json"
// 	"log"
// 	"net/http"

// 	"libvirt.org/go/libvirt"
// )

// type DomainInfo struct {
// 	Name     string `json:"name"`
// 	UUID     string `json:"uuid"`
// 	State    string `json:"state"`
// }

// type VMField struct {
// 	UUID   string     `json:"uuid"`
// 	Name   string     `json:"name"`
// 	State  string     `json:"state"`
// }

// type VMEvent struct {
// 	VM VMField `json:"vm"`
// }

// type DataWrapper struct {
//     Data []VMEvent `json:"data"`
// }

// func getDomains() ([]DomainInfo, error) {
// 	conn, err := libvirt.NewConnect("qemu:///system")
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer conn.Close()

// 	domains, err := conn.ListAllDomains(0)
// 	if err != nil {
// 		return nil, err
// 	}

// 	var result []DomainInfo
//     log.Printf("Scraped %d VMs\n", len(domains))
// 	for _, dom := range domains {
//         log.Printf("%+v\n", dom)
// 		name, _ := dom.GetName()
// 		uuid, _ := dom.GetUUIDString()
// 		info, _ := dom.GetInfo()

// 		result = append(result, DomainInfo{
// 			Name:     name,
// 			UUID:     uuid,
// 			State:    domainStateToString(info.State),
// 		})
// 	}
// 	return result, nil
// }

// func domainStateToString(state libvirt.DomainState) string {
// 	switch state {
// 	case libvirt.DOMAIN_RUNNING:
// 		return "running"
// 	case libvirt.DOMAIN_PAUSED:
// 		return "paused"
// 	case libvirt.DOMAIN_SHUTDOWN:
// 		return "shutdown"
// 	case libvirt.DOMAIN_SHUTOFF:
// 		return "shutoff"
// 	case libvirt.DOMAIN_CRASHED:
// 		return "crashed"
// 	default:
// 		return "unknown"
// 	}
// }

// // Convert domains to ECS-compatible events
// func domainsToEvents(domains []DomainInfo) DataWrapper {
// 	events := make([]VMEvent, 0, len(domains))

// 	for _, d := range domains {
//         events = append(events, VMEvent{
//             VM: VMField{
//                 UUID:  d.UUID,
//                 Name:  d.Name,
//                 State: d.State,
//             },
//         })
// 	}

// 	return DataWrapper{
//         Data: events,
//     }
// }

// func main() {
// 	http.HandleFunc("/v1/domains", func(w http.ResponseWriter, r *http.Request) {
// 		domains, err := getDomains()
// 		if err != nil {
// 			http.Error(w, err.Error(), 500)
// 			return
// 		}

// 		events := domainsToEvents(domains)
// 		w.Header().Set("Content-Type", "application/json")
// 		json.NewEncoder(w).Encode(events)
// 	})

// 	log.Println("Serving JSON Elastic Agent endpoint on http://0.0.0.0:8088/v1/domains")
// 	log.Fatal(http.ListenAndServe(":8088", nil))
// }







package main

import (
	"encoding/json"
	"log"
	"net/http"

	"libvirt.org/go/libvirt"
)

type DomainInfo struct {
	Name  string `json:"name"`
	UUID  string `json:"uuid"`
	State string `json:"state"`
}

// Struct for the summarized output
type LibvirtSummary struct {
	Libvirt map[string]int `json:"libvirt"`
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
		name, _ := dom.GetName()
		uuid, _ := dom.GetUUIDString()
		info, _ := dom.GetInfo()

		result = append(result, DomainInfo{
			Name:  name,
			UUID:  uuid,
			State: domainStateToString(info.State),
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

// Count the number of domains in each state
func summarizeDomains(domains []DomainInfo) LibvirtSummary {
	counts := map[string]int{
		"running": 0,
		"paused":  0,
		"shutdown": 0,
		"shutoff": 0,
		"crashed": 0,
		"unknown": 0,
	}

	for _, d := range domains {
		if _, exists := counts[d.State]; exists {
			counts[d.State]++
		} else {
			counts["unknown"]++
		}
	}

	return LibvirtSummary{Libvirt: counts}
}

func main() {
	http.HandleFunc("/v1/domains", func(w http.ResponseWriter, r *http.Request) {
		domains, err := getDomains()
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}

		summary := summarizeDomains(domains)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(summary)
	})

	log.Println("Serving JSON Elastic Agent endpoint on http://0.0.0.0:8088/v1/domains")
	log.Fatal(http.ListenAndServe(":8088", nil))
}