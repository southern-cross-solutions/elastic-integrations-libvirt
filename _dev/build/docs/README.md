# Libvirt/KVM Integration for Elastic

## Overview
The `Libvirt/KVM` integration for Elastic enables collection of data such as number of Total, Running, Shut Off, No State, Suspended, Crashed & VM states

## Enabling the integration in Elastic
1. In Kibana go to Management > Integrations
2. In "Search for integrations" search bar type `Libvirt/KVM`
3. Click on `Libvirt/KVM` integration from the search results.
4. Click on Add `Libvirt/KVM` button to add `Libvirt/KVM` integration.

## What do I need to use this integration?
<!-- TODO -->

## Logs Reference
<!-- TODO -->

<!-- ### {Data stream name}

The `{data stream name}` data stream provides events from Libvirt/KVM.

#### Example.
TODO : Example Json
```json
{
    "@timestamp": "2023-10-31T07:31:24.050Z",
    "agent": {
        "ephemeral_id": "bf237146-2d4b-427b-b731-6dadb1dfdd90",
        "id": "fa60f5ca-bf95-4706-9195-907dd5f9b537",
        "name": "docker-fleet-agent",
        "type": "filebeat",
        "version": "8.4.1"
    },
    "bitwarden": {
        "collection": {
            "external": {
                "id": "external_id_123456"
            },
            "id": "539a36c5-e0d2-4cf9-979e-51ecf5cf6593"
        },
        "object": "collection"
    },
    "data_stream": {
        "dataset": "bitwarden.collection",
        "namespace": "ep",
        "type": "logs"
    },
    "ecs": {
        "version": "8.11.0"
    },
    "elastic_agent": {
        "id": "fa60f5ca-bf95-4706-9195-907dd5f9b537",
        "snapshot": false,
        "version": "8.4.1"
    },
    "event": {
        "agent_id_status": "verified",
        "created": "2023-10-31T07:31:24.050Z",
        "dataset": "bitwarden.collection",
        "ingested": "2023-10-31T07:31:27Z",
        "kind": "event",
        "original": "{\"externalId\":\"external_id_123456\",\"groups\":null,\"id\":\"539a36c5-e0d2-4cf9-979e-51ecf5cf6593\",\"object\":\"collection\"}",
        "type": [
            "info"
        ]
    },
    "input": {
        "type": "httpjson"
    },
    "tags": [
        "preserve_original_event",
        "preserve_duplicate_custom_fields",
        "forwarded",
        "bitwarden-collection"
    ]
}
```

#### Exported fields
TODO : Table of fields, descriptions, types
|               Field              |                                         Description                                        |       Type       |
|:--------------------------------:|:------------------------------------------------------------------------------------------:|:----------------:|
| @timestamp                       | Event timestamp.                                                                           | date             |
| bitwarden.collection.external.id | External identifier for reference or linking this collection to another system.            | keyword          |
| bitwarden.collection.groups      | The associated groups that this collection is assigned to.                                 | nested           |
| bitwarden.collection.id          | The collection's unique identifier.                                                        | keyword          |
| bitwarden.object                 | String representing the object's type. Objects of the same type share the same properties. | keyword          |
| data_stream.dataset              | Data stream dataset.                                                                       | constant_keyword |
| data_stream.namespace            | Data stream namespace.                                                                     | constant_keyword |
| data_stream.type                 | Data stream type.                                                                          | constant_keyword |
| event.dataset                    | Event dataset.                                                                             | constant_keyword |
| event.module                     | Event module.                                                                              | constant_keyword |
| input.type                       | Type of Filebeat input.                                                                    | keyword          |
| log.offset                       | Log offset.                                                                                | long             | -->