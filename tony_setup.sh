#! /bin/bash
set -x

# Clean up any leftovers from last time.
kubectl delete -f mysql-resource.yaml
kubectl delete -f mysql-operator.yaml
kubectl delete crd mysqls.tallen.io

# Just go ahead and rebuild everything for good measure.
./codegen.sh
CGO_ENABLED=0 GOOS=linux go build
eval $(minikube docker-env)
docker build -t mysql-operator:0.1 .

# Start the operator.
kubectl create -f mysql-operator.yaml
sleep 5
kubectl get pod -l app=mysql-operator
