#!/usr/bin/env bash
if [ -v BASH_DEBUG ]; then set -x; fi
set -e
set -o pipefail
set -o nounset

echo "start e2e test"
kind create cluster --config kind-config.yaml
sleep 10

cd ..
make
make install

echo "start deploy dns and controller"
kubectl apply -f e2e/nipPod.yaml
kubectl apply -f e2e/dep.yaml

echo "waiting for finish deploy"
sleep 25

echo "start deploy pod test and svc"
kubectl apply -f e2e/pod2.yaml
kubectl apply -f e2e/svc.yaml

echo "start deploy ing"
kubectl apply -f e2e/ing.yaml

sleep 20
echo "try to curll"
curl test.127.0.0.1.nip.io






