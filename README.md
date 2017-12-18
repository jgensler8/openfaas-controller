# OpenFaaS Controller

This controller reads a CRD and creates functions in an OpenFaas controller

## Compiling

```
glide install
go build
```

## Running

```
$ ./openfaas-controller \
  --kubeconfig ~/.kube/config \
  --openfaas-basepath http://192.168.99.100:31112
```

At the moment, the controller doesn't create the CRD (though it could in the future).

```
kubectl create -f examples/crd-definition.yaml
kubectl create -f function.yaml
```

## Sample Definition

```
apiVersion: "jgensler8.openfaas.apiextensions.k8s.io/v1"
kind: "Function"
metadata:
  name: "samplefunction"
spec:
  service: "myservicetwo"
  image: "nginx:1.13"
  envVars: []
  envProcess: ""
  registryAuth: ""
  network: ""
```