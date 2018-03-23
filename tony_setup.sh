#! /bin/bash
set -x

# Clean up any leftovers from last time.
kubectl delete -f mysql-resource.yaml
kubectl delete -f mysql-operator.yaml
kubectl delete crd mysqls.tallen.io

set -e

# Just go ahead and rebuild everything for good measure.
./codegen.sh
CGO_ENABLED=0 GOOS=linux go build
eval $(minikube docker-env)
docker build -t mysql-operator:0.1 .
docker save mysql-operator:0.1 | eval $(minikube docker-env)

# Start the operator.
kubectl create -f mysql-operator.yaml
sleep 3
kubectl get pod -l app=mysql-operator

# Create the resource.
kubectl create -f mysql-resource.yaml
sleep 3
kubectl get mysqls.tallen.io
