# adapter-base
adapter-base framework written in Golang to provide stateful RPC and HTTP API's
exposing `scan` and `update` functionality to gateway applications.

This framework will be used as the basis of support for new dependent device
types with only the `update` and `scan` functions needing implementation.

## How it works
* When the `start` command is called a new background worker is created, this
  worker is assigned an ID which is returned by the API. If no background
  workers are available an error message is returned.
* When the `status` command is called with the worker ID the status of the
  worker is returned. If the worker has finished the worker is terminated and
  the final status returned.
* When the `cancel` command is called with the worker ID the worker is
  terminated and the final status returned.

Note: A worker will not be terminated until the final `status` has been
requested, even if it's work is finished.

## RPC command line interface
```
Usage:
  adapter-base [command]

Available Commands:
  help        Help about any command
  scan        Interact with the adapter-base scan endpoint
  server      Start the adapter-base server
  update      Interact with the adapter-base update endpoint

Flags:
  -h, --help      help for adapter-base
  -v, --verbose   verbose output

Use "adapter-base [command] --help" for more information about a command.
```

## HTTP interface
### Scan
#### Start
```
curl -X POST localhost:8080/scan --data '{"name": "test_name", "timeout": 120}' | jq .
{
  "id": "26eeec7d-bac1-43e6-a1dd-3ffed62c375f"
}
```

#### Status
```
curl -X GET localhost:8080/scan/26eeec7d-bac1-43e6-a1dd-3ffed62c375f | jq .
{
  "id": "26eeec7d-bac1-43e6-a1dd-3ffed62c375f",
  "startRequest": {
    "name": "test_name",
    "timeout": "120"
  },
  "results": [
  ],
  "started": "1501759709"
}
```

#### Cancel
```
curl -X DELETE localhost:8080/scan/26eeec7d-bac1-43e6-a1dd-3ffed62c375f | jq .
{
  "id": "26eeec7d-bac1-43e6-a1dd-3ffed62c375f",
  "startRequest": {
    "name": "test_name",
    "timeout": "120"
  },
  "results": [...],
  "started": "1501759709"
}
```

### Update
#### Start
```
curl -X POST localhost:8080/update --data '{"address": "test_address", "payload": "test_payload", "timeout": 120}' | jq .
{
  "id": "7ab875ac-2683-4f19-95ee-7498179729ab"
}
```

#### Status
```
curl -X GET localhost:8080/update/7ab875ac-2683-4f19-95ee-7498179729ab | jq .
{
  "id": "7ab875ac-2683-4f19-95ee-7498179729ab",
  "startRequest": {
    "address": "test_address",
    "payload": "test_payload",
    "timeout": "120"
  },
  "state": "FLASHING",
  "progress": 35,
  "message": "message: 35",
  "started": "1501759227"
}
```

#### Cancel
```
curl -X DELETE localhost:8080/update/7ab875ac-2683-4f19-95ee-7498179729ab | jq .
{
  "id": "89969de7-d29e-4464-b4b8-84eb43a6baaf",
  "startRequest": {
    "address": "test_address",
    "payload": "test_payload",
    "timeout": "120"
  },
  "state": "FLASHING",
  "progress": 54,
  "message": "message: 54",
  "started": "1501759355"
}
```

## Development
* Ensure you have the [protobuf](https://github.com/golang/protobuf)
  dependencies installed
* Install the project dependencies with `glide install`
* Compile the protobuf definitions with `./protoc.sh`
* Build the project with `go build`
