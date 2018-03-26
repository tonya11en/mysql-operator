# MySQL Operator
A Kubernetes custom resource and Operator that allows a user to describe a trivial single-instance of MySQL. The tasks for [running a single-instance stateful application](https://kubernetes.io/docs/tasks/run-application/run-single-instance-stateful-application/#accessing-the-mysql-instance) are done by this operator.

## Building
```bash
# From the root of repo, ensure all the libraries required are pulled down.
dep ensure

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

# Delete the MySQL resource.
kubectl delete -f mysql-resource.yaml
```
