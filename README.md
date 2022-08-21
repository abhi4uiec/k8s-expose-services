# Platform Developer Test

## Tasks

### 1. Expose information on all pods in the cluster

Add an endpoint to the service that exposes the number of pods running in the cluster in namespace `default` per service
and per application group:

```
GET `/services`
[
  {
    "name": "<service>",
    "applicationGroup": "alpha",
    "runningPodsCount": 2
  },
  {
    "name": "<service>",
    "applicationGroup": "beta",
    "runningPodsCount": 1
  },
  ...
]
```

### 2. Expose information on a group of applications in the cluster

Create an endpoint in your service that exposes the pods in the cluster in namespace `default` that are part of the same `applicationGroup`:

```
GET `/services/{applicationGroup}`
[
  {
    "name": "<service>",
    "applicationGroup": "<applicationGroup>",
    "runningPodsCount": 1
  },
  ...
]
```

## How to execute a program

1. Clone the git repo
2. Please apply the services to your local Kubernetes cluster by executing
      
         `kubectl apply -f ./services.yaml`.

3. Launch the program using:

        go run main.go

4. To get details of exposed information on all pods in the cluster, run curl command:

        curl http://localhost:9090/services

   To get details of exposed information on a group of applications in the cluster, run the curl command:

        curl http://localhost:9090/services/{applicationGroup}

        ex. curl http://localhost:9090/services/alpha

## How to execute TC

  * Execute the below command to run tests:
      
        go test -v

## Assumption

  * Program will throw an error if there will be more than one deployment using same label selector. (This is not the correct approach)
  * If the pod doesn't belongs to any application group, then i'm putting the applicationGroup as "none". This will help to fetch all pods having no applicationGroup.

## Improvements

  * Instead of using http server, use https server
    Will need to use cert in that case, or use self signed certificates