# SRE Take Home Project

## Development

Using the development docker the image Dockerfile.dev:
```bash
# Build
podman build --pull --rm -f 'Dockerfile.dev' -t 'sretakehome:dev' '.'
#enter in container:
podman run --rm -it -p 3088:3088 -v $(pwd):/app sretakehome:dev
#compile and run
go build -o bin/key-server.go main.go
bin/key-server.go --max-size 2048 --srv-port 3088

# Run tests
go test ./keyserver
```

### Compose file

Testing integration with prometheus
```bash
# Build
podman-compose build
# Run
podman-compose up
```
## Deployment

Deployment via helm charts using local docker images. In a real scenario Docker images would be obtained via repository (i.e. docker hub)

To import local images in k3s run:
```bash
podman save pcosta/keyserver:1.0 | sudo k3s ctr images import -
```

To deploy the full service (keyserver + prometheus)
```bash
 helm install keyserver-release keyserverapp/
```

### Accessing the services

The keyserver service is exposed via NodePort. On launch the Kubernetes control plane will allocate the service to a port inside the default range 30000-32767. In a real world scenario the service could be exposed to the public internet using an Ingress.

To check the allocated port run the following command after deploying: 
```bash
kubectl -n default get svc | grep keyserver-service | awk '{print $5}' | cut -d: -f2
```

The remaining services (Prometheus, AlertManager, ...) use ClusterIP, meaning they are only reachable from within the cluster. To expose those services use the port-forwarding, for example to expose prometheus in port 9090:
```bash
kubectl --namespace default port-forward <prometheus-pod-name> 9090
```

## Monitoring notifications

Alert manager is used to notify and alert undesired behaviours. The alert manager configuration is done in **values.yaml** and four types of alerts were defined:

- High 5xx errors (application errors)
- High 400 errors (users trying to generate keys with incorrect settings)
- High memory usage
- Readiness probe failure (failed deployments)

As is the alerting is done only on the alert manager interface, but in a real scenario we could associate a slack channel or email address to send those alerts.
