# Drone Networks Service (DNS)

This is a fully functional drone navigation service, it is designed to help drones quickly locate databanks to upload gathered data across a sector in the galaxy.

## How it Works

Currently the DNS exposes only one endpoint `/v1/locate`. This endpoints requires a `header` value and a json `payload` to work properly. Check below for details.

## How to Run

You can simply run `go run ./cmd/api` to start the server. The server should start on port `8080`. Visit <localhost:8080> to test it.

### Docker

This is my default way of running this program locally, so the rest of this docs will assume you will be doing the same.

There are two `Dockerfiles` in the build folder `Dockerfile.dev` and `Dockerfile.deploy`. The extension suggests their different use cases.

#### Development

You can run the development server using Docker which uses [github.com/codegangsta/gin](gin) to rebuild and restarts the DNS web server, gin relays requests to the DNS server via it's own proxy server that listens on port `3000`. To test this visit [localhost:3000](localhost:3000).

To start the development server run the following commands.

```bash
    $ docker build --rm -f ./build/Dockerfile.dev -t dns .
```

This command build an alpine based docker image, it exposes two ports, one for `gin` and the other for our Go web server.

To run it:

```bash
    $ docker run -it --rm -p 8080:8080 -p 3000:3000 -v $PWD:/go/src/dns dns
```

This commands runs the `dns:latest` image that we built earlier. It also mounts the current working directly volume to the docker container, so that any change we make to our code is reflected within the container.

#### Production

Building and running te production image is similar to development, run the following command to build the production image:

```bash
    $ docker build --rm -f ./build/Dockerfile.deploy -t dns-prod .
```

This commands builds a smaller binary for our productin server.

We can then run the image like this:

```bash
    $ docker run -it --rm -p 8080:8080 dns-prod
```

Visit [localhost:8080](localhost:8080) to test. To run in deamon mode remove the `-it` flag and replace it with `-d` flag.

## API Reference

Once you have the server running successfully you can start making http request.

| | | | |
|-|-|-|-|
| __ENDPOINT__ | __HTTP Verb__ | __Header__ | __PayLoad__ | __Description__
| `/v1/locate` | POST | `requires` that `X-System-Type` header is set to supported systems which is `drone` or `ship` | The payload data are numeric values sent as strings eg. `{ "x": "123.12", "y": "456.56", "z": "789.89", "vel": "20.0" }`.

## Testing

To run test:

```bash
    $ go test ./cmd/api -race -cover
```

This is a relatively small API, so achieving a 100% test coverage was easy.

## TODO

+ [ ] Intgrate tracing capability with Jeager
+ [ ] Write build script or Makefile
+ [ ] Configure CI/CD

## Copyright © All rights reserved

**Atlas Corporation.**
