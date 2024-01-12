# Golang Test Apps

## Description

`/http` - simple http server
`/websockets` - simple websocket server

## Usage

```bash
# Install dependencies
brew upgrade && brew update
brew install kubectl
brew install docker
brew install kind

# Create local cluster and registry
./kind-with-registry.sh

# Build demo app images (example for http)
# Build the client image
(cd ./http/client && \
  docker build -t 'http-client' -t 'http-client:latest' . && \
  # tag and publish to the local registry
  docker tag 'lua-tracing-poc-client:latest' 'localhost:5001/http-client:latest' && \
  docker push 'localhost:5001/http-client:latest')

# Build the server image
(cd ./http/server && \
  docker build -t 'http-server' -t 'http-server:latest' . && \
  # tag and publish to the local registry
  docker tag 'http-server:latest' 'localhost:5001/http-server:latest' && \
  docker push 'localhost:5001/http-server:latest')
```

Now you can use the test images in your cluster by referring to them in the manifests as
```bash
# client image
localhost:5001/http-client:latest
# server image
localhost:5001/http-server:latest
```