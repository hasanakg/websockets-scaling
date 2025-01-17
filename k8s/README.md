# K8S version

This configuration will use a service interfacing a deployment with 3 replicas of socket-server application.
The back part of Redis is handled by a classic Redis Master/Slave configuration.

## Websocket Deployment/Service (folder wsk/)

The service is packed and exposed from static version via a Docker image deployed on private registry.
The exposed port of the container / pod is 5000. The service can be exposed by a Traefik ingress or via a `NodePort` on port 30000.

The key point of the service is `sessionAffinity: ClientIP` which will mimic Sticky connection for the service.

* Note: In general Round Robin among PODs is not guaranteed using a simple solution like NodePort. In order to mimic Round Robin policy, this configuration should be added:

      sessionAffinityConfig:
        clientIP:
          timeoutSeconds: 10

## Redis Deployment/Service (folder redis/)

A Master-Slave K8S for Redis solution that is maintaining in-sync websockets through Pub/Sub, using the endpoint:

* **redis-master.default.svc.cluster.local`**

## Traefik Ingress Route (folder traefik/)

A Traefik `Ingress Route` can be used to proxy request towards websocket service.

## Setup and Run

### microservices

    # Create with create script
    ./create.sh

### k8s dashboard address
Get token from create.sh output and login with following link
https://localhost:31443/#/login

### client

The client configuration is different depending on cluster setup:

* `NodePort Service`: *<http://localhost:30000>*
* `Traefik Ingress`: *<http://localhost:30080>* (also path is `wsk/`)

For any new client open a new terminal and run:

    node ../stack/client_socket.js
