# KCD Prometheus Workshop

## Contents

- [System Requirements](#system-requirements)
- [Folders](#folders)
  - [Server](#server)
  - [Client](#client)
  - [Prometheus](#prometheus)
  - [Alertmanager](#alertmanager)
  - [Grafana](#grafana)
- [Running](#running)

## System Requirements

Note: Below requirements are based on Mac OS

- Docker version >= 20.10 (Probably works with an older version also)
- Golang version >= 1.20 (Probably works with an older version also only required for IDE)
- GNU Make version >= 3.81 (Probably works with an older version also and if make **doesn't exist run the commands manually**)

## Folders

### Server

Suggestions application server code

### Client

Client code that will invoke the Suggestion application APIs

### Prometheus

Prometheus configuration and rules

### Alertmanager

Alertmanager configuration

### Grafana

Grafana dashboards for the node exporter original file can be found [here](https://grafana.com/grafana/dashboards/1860-node-exporter-full/) and the suggestion Application dashboard

## Running

```bash
# `make` command to see the help for the command
make setup # to setup the docker network
make run-server
make run-node-exporter
make run-prometheus
make run-alertmanager
make run-grafana
make run-pushgateway
make run-prometheus-secondary
make run-prometheus-federation
make run-client # tweak cocurrentRoutines, iterations and delay
# `make clean` to remove the containers and network
# `make all` to run all the containers
```
