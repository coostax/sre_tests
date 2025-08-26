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
