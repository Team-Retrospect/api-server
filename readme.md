# API Server

A high-powered API server that exposes the Cassandra data in a controlled manner.

## Installation with docker-compose.yml file

Download the [docker-compose.yml](https://github.com/Team-Retrospect/docker-files) file

```sh
docker-compose up
```

## Manual Installation

First, install the dependencies:
- [Make sure Go is installed](https://www.digitalocean.com/community/tutorials/how-to-install-go-on-debian-10)
- Clone this repo and cd into it
- Rename `config.yml.example` to `config.yml` and configure it
- Build the app with `go build`
- Run it with `./cassandra-connector`

Note that the Cassandra instance must be preconfigured. This app does not initialize it, but the schema is available in [data.cql](./data.cql)

## Usage

Send HTTP/S traffic (GET, POST, ...) to the server, using the port specified in the config file.

[See the documentation here](https://retrospect-api.api-docs.io/0.9.0/)

## Testing

Install the dependencies:

```sh
apt update
apt install python3
```

```sh
pip3 install pytest pyyaml requests
```

Run the tests from the root folder with `./tests.sh`

```sh
$ ./tests.sh
============================= test session starts ==============================
platform darwin -- Python 3.9.5, pytest-6.2.4, py-1.10.0, pluggy-0.13.1
rootdir: /Users/nicole/Capstone/Retrospect/api-server
collected 25 items

tests/chapters.py ......                                                 [ 24%]
tests/events.py ......                                                   [ 48%]
tests/snapshots.py ....                                                  [ 64%]
tests/spans.py ........                                                  [ 96%]
tests/trigger_routes.py .                                                [100%]

============================== 25 passed in 9.10s ==============================
```
