#!/usr/bin/env bash

# Enable debugging if BASH_DEBUG is set
if [ -v BASH_DEBUG ]; then set -x; fi
set -e
set -o pipefail
set -o nounset

# Function to check if a command exists
command_exists() {
    command -v "$1" >/dev/null 2>&1
}

# Ensure required commands are available
for cmd in kind kubectl curl; do
    if ! command_exists "$cmd"; then
        echo "Error: $cmd is not installed." >&2
        exit 1
    fi
done

# Function to wait for a condition with a timeout
wait_for_condition() {
    local condition="$1"
    local timeout="$2"
    local interval="$3"

    local start_time
    start_time=$(date +%s)

    while ! eval "$condition"; do
        sleep "$interval"
        if (( $(date +%s) - start_time >= timeout )); then
            echo "Timeout waiting for condition: $condition" >&2
            return 1
        fi
    done
}

echo "Starting end-to-end tests"

# Create Kubernetes cluster with kind
kind create cluster --config kind-config.yaml
echo "Cluster creation initiated, waiting for nodes to be ready..."

# Wait for nodes to be ready
wait_for_condition "kubectl get nodes | grep -q Ready" 120 5
echo "Cluster nodes are ready"

# Build and install the project
echo "Building and installing project"
(
    cd ..
    make
    make install
)

echo "Deploying DNS and controller"
kubectl apply -f e2e/nipPod.yaml
kubectl apply -f e2e/dep.yaml

# Wait for deployments to complete
echo "Waiting for DNS and controller to be ready..."
wait_for_condition "kubectl get pods -l app=nipPod -o jsonpath='{.items[*].status.containerStatuses[*].ready}' | grep -q true" 60 5
wait_for_condition "kubectl get pods -l app=dep -o jsonpath='{.items[*].status.containerStatuses[*].ready}' | grep -q true" 60 5

echo "Deploying test pod and service"
kubectl apply -f e2e/pod.yaml
kubectl apply -f e2e/svc.yaml

# Wait for test pod to be ready
echo "Waiting for test pod to be ready..."
wait_for_condition "kubectl get pods -l app=pod2 -o jsonpath='{.items[*].status.containerStatuses[*].ready}' | grep -q true" 60 5

echo "Deploying ingress"
kubectl apply -f e2e/ing.yaml

# Wait for ingress to be ready
echo "Waiting for ingress to be ready..."
wait_for_condition "kubectl get ingress -o jsonpath='{.items[*].status.loadBalancer.ingress[*].ip}' | grep -q '^[0-9]'" 60 5

echo "Testing ingress with curl"
curl -sS test.127.0.0.1.nip.io || { echo "Curl test failed"; exit 1; }

echo "End-to-end tests completed successfully"