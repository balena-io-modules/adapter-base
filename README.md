# adapter-base
adapter-base framework written in Golang to provide stateful RPC and HTTP API's
exposing `scan` and `update` functionality to gateway applications.

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
requested, even if it has finished its work.

## Supported modes
* `simulate` - simulation mode that returns a series of numbers

## Development
* Install [protobuf](https://github.com/golang/protobuf) dependencies
* Build the project with `./build.sh`
