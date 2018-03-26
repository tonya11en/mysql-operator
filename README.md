# MySQL Operator
A Kubernetes custom resource and Operator that allows a user to describe a trivial single-instance of MySQL. The tasks for [running a single-instance stateful application](https://kubernetes.io/docs/tasks/run-application/run-single-instance-stateful-application/) are done by this operator.

## Pre-reqs
You can run Kubernetes locally with
[Minikube](https://kubernetes.io/docs/getting-started-guides/minikube/).

## Building
```bash
# From the root of repo, ensure all the libraries required are pulled down.
dep ensure

# build the sample operator binary
CGO_ENABLED=0 GOOS=linux go build

# Build the docker container.
docker build -t mysql-operator:0.1 .
docker save mysql-operator:0.1 | eval $(minikube docker-env) # Only needed for minikube.
```

## Using the Operator
```bash
# Start the operator.
kubectl create -f mysql-operator.yaml

# Create a MySQL resource (might take a few seconds).
kubectl create -f mysql-resource.yaml

# Verify the MySQL resource is created by creating a new pod in the cluster and
# running a MySQL client that connects to the newly created service.
kubectl run -it --rm --image=mysql:5.6 --restart=Never mysql-client -- mysql -h mysql -ppassword
```

There is no support for updating the resource because there is nothing in its
spec that will allow for a non-disruptive change to the service. To make
changes to the service, it's recommended to tear it down and redeploy.

## Cleanup
```bash
kubectl delete -f mysql-resource.yaml
kubectl delete -f mysql-operator.yaml
```
