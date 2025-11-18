# Libvirt/KVM Integration for Elastic

## Overview
The `Libvirt/KVM` integration for Elastic enables collection of data such as number of Total, Running, Shut Off, No State, Suspended, Crashed & VM states

## Enabling the integration in Elastic
1. In Kibana go to Management > Integrations
2. In "Search for integrations" search bar type `Libvirt/KVM`
3. Click on `Libvirt/KVM` integration from the search results.
4. Click on Add `Libvirt/KVM` button to add `Libvirt/KVM` integration.

## What do I need to use this integration?
1. A Hypervisor running Libvirt that is onboarded into an Agent Policy.
2. The libvirt-exporter built and running on the hypervisor (build : `go build -o bin/libvirt-exporter -ldflags "-s -w" cmd/main.go`).
3. The libvirt-exporter running and available at a configurable endpoint (run : `/path/to/libvirt-exporter`).

## Metrics Reference

### Libvirt Metrics

The `metrics` data stream provides the states of the Libvirt/KVM hosted Virtual Machines.

#### Example.
```json
{
    "@timestamp": "2025-11-04T07:32:20.094Z",
    "libvirt": {
        "running": 1,
        "shutoff": 1,
        "paused": 0,
        "crashed": 0,
        "shutdown": 0,
        "unknown": 0
    },
    "input": {
        "type": "httpjson"
    },
    "agent": {
        "name": "libvirt-test-host",
        "id": "f1689110-e83a-4e7e-bfae-b2bba5ea79b2",
        "ephemeral_id": "66245cfd-c27a-431a-8437-c9dbfdf6d5d8",
        "type": "filebeat",
        "version": "8.19.6"
    },
    "ecs": {
        "version": "8.0.0"
    },
    "data_stream": {
        "namespace": "default",
        "type": "metrics",
        "dataset": "libvirt.metrics"
    },
    "event": {
        "kind": "metric",
        "module": "libvirt",
        "dataset": "libvirt.metrics"
    },
    "@version": "1",
    "elastic_agent": {
        "id": "f1689110-e83a-4e7e-bfae-b2bba5ea79b2",
        "version": "8.19.6",
        "snapshot": false
    },
    "dimension": "time_series"
}
```

#### Exported fields
|               Field              |                                         Description                                        |       Type       |
|:--------------------------------:|:------------------------------------------------------------------------------------------:|:----------------:|
| @timestamp                       | Event timestamp.                                                                           | date             |
| libvirt.[state]                  | Number of vms that are currently in that state on the Libvirt Hypervisor.                  | long             |
| dimension                        | A "empty" field that is required for time series data.                                     | constant_keyword |
| data_stream.dataset              | Data stream dataset.                                                                       | constant_keyword |
| data_stream.namespace            | Data stream namespace.                                                                     | constant_keyword |
| data_stream.type                 | Data stream type.                                                                          | constant_keyword |
| event.dataset                    | Event dataset.                                                                             | constant_keyword |
| event.module                     | Event module.                                                                              | constant_keyword |
| event.kind                       | Event kind.                                                                                | constant_keyword |
| input.type                       | Type of Filebeat input.                                                                    | constant_keyword |
| agent.name                       | Host name of the hypervisor.                                                               | constant_keyword |