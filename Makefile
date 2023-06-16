
SHELL := /usr/bin/env bash -o errexit -o pipefail -o nounset
MAKEFLAGS += -s
.EXPORT_ALL_VARIABLES:
.DEFAULT_GOAL := help

# You can also you a docker compose and --rm during the docker run commands

.PHONY: help
help: ## Display this help
	@awk 'BEGIN {FS = ":.*##"; printf "  KCD Chennai Prometheus workshop helpers\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z0-9_-]+:.*?##/ { printf "  \033[36m%-35s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

.PHONY: setup
setup: ## Creates a docker network for running the containers
	docker network create kcd-network

.PHONY: build-server
build-server: ## Builid the server image
	docker build -t kcd-server ./server/

.PHONY: run-server
run-server: build-server ## Build and start the server
	docker rm -f kcd-server || true
	docker run --name kcd-server -p 8089:8089 --network kcd-network -d kcd-server

.PHONY: build-client
build-client: ## Builid the client image
	docker build -t kcd-client ./client/

.PHONY: run-client
run-client: build-client ## Build and run the client (make sure to start the server before)
	docker rm -f kcd-client || true
	docker run --name kcd-client --network kcd-network kcd-client

.PHONY: run-prometheus
run-prometheus: ## Running prometheus as a docker container
	docker rm -f kcd-prometheus || true
	docker run --name kcd-prometheus -p 9090:9090 --network kcd-network -d \
	-v $(shell pwd)/prometheus/prometheus.yml:/etc/prometheus/prometheus.yml \
	-v $(shell pwd)/prometheus/rules.yml:/etc/rules/rules.yml \
	prom/prometheus

.PHONY: run-secondary-prometheus
run-secondary-prometheus: ## Running another prometheus instace that collects metrics from pushgateway
	docker rm -f kcd-secondary-prometheus || true
	docker run --name kcd-secondary-prometheus -p 9092:9090 --network kcd-network -d \
	-v $(shell pwd)/prometheus/prometheus-secondary.yml:/etc/prometheus/prometheus.yml \
	prom/prometheus

.PHONY: run-node-exporter
run-node-exporter: ## Running node exporter as a docker container
	docker rm -f kcd-node-exporter || true
	docker run --name kcd-node-exporter -p 9100:9100 --network kcd-network -d \
	prom/node-exporter

.PHONY: run-alertmanager
run-alertmanager: ## Running prometheus alert manager
	docker rm -f kcd-alertmanager || true
	docker run --name kcd-alertmanager --network kcd-network -p 9093:9093 -d \
	-v $(shell pwd)/alertmanager/alertmanager.yml:/etc/alertmanager/alertmanager.yml \
	-v $(shell pwd)/alertmanager/beer_slack_webhook:/etc/slack/beer_slack_webhook \
	-v $(shell pwd)/alertmanager/suggestion_slack_webhook:/etc/slack/suggestion_slack_webhook \
	prom/alertmanager --config.file=/etc/alertmanager/alertmanager.yml

.PHONY: run-pushgateway
run-pushgateway: ## Running prometheus pushgateway
	docker rm -f kcd-pushgateway || true
	docker run --name kcd-pushgateway --network kcd-network -p 9091:9091 -d \
	prom/pushgateway

.PHONY: run-job
run-job: ## Running the job
	docker rm -f kcd-job || true
	docker build -t kcd-job ./job/
	docker run --name kcd-job --network kcd-network kcd-job

.PHONY: delete-job-metrics
delete-job-metrics: ## Delete the metrics pushed by job in pushgateway
	curl -X DELETE http://localhost:9091/metrics/job/kcd-job/

.PHONY: run-grafana
run-grafana: ## Running grafana
	docker rm -f kcd-grafana || true
	docker run --name kcd-grafana --network kcd-network -p 9200:3000 -d \
	grafana/grafana-oss

.PHONY: run-prometheus-federation
run-prometheus-federation: ## Running prometheus in federation pulling metrics from first and secondary prometheus
	docker rm -f kcd-prometheus-federation || true
	docker run --name kcd-prometheus-federation -p 9094:9090 --network kcd-network -d \
	-v $(shell pwd)/prometheus/prometheus-federation.yml:/etc/prometheus/prometheus.yml \
	prom/prometheus --enable-feature=native-histograms

.PHONY: all
all: run-server run-node-exporter run-prometheus run-alertmanager run-grafana run-client ## Run all the containers

.PHONY: clean
clean: ## Deletes all the containers
	docker rm -f kcd-job || true
	docker rm -f kcd-server || true
	docker rm -f kcd-client || true
	docker rm -f kcd-node-exporter || true
	docker rm -f kcd-alertmanager || true
	docker rm -f kcd-grafana || true
	docker rm -f kcd-prometheus || true
	docker rm -f kcd-pushgateway || true
	docker rm -f kcd-secondary-prometheus || true
	docker rm -f kcd-prometheus-federation || true
	docker network rm kcd-network || true
