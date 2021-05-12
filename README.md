# budget-tracker-api

API REST Backend for a personal budget tracker.

It also has the goal to implement engineering best practices and new technologies for study purposes.

## Architecture

<img src="imgs/budget-tracker.png" alt="">

## Protocols

The backend supports both HTTP/1.1 (with optional TLS) and HTTP/2 (with mandatory TLS) protocols. You can simply play with both `hc.InitHTTPServer()` and `InitHTTP2Server()` methods for each one of the protocols at the `main.go` file.

# Developer tools

## Running locally

You can use `docker-compose` to run the entire backend stack locally: `budget-tracker` and `mongodb` (with an initial `admin` user created)

The mongodb served by `docker-compose` has no credentials so it's recommended only for development purposes.

## Swagger API support

This application uses go-swagger to generate swagger specs directly from the code, to run it just:

`make swagger-spec`

The output will be a local `./docs/swagger.yaml` file which can be used along with external applications such as insomnia designer. In the case of side-by-side UI containers such as `swagger-ui` you can serve the following endpoint: `/swagger.yaml`

## Insomnia

Download `./docs/insomnia.json` file and upload it to your insomnia application to get all endpoints to be tested

# Observability

## Opentelemetry

This application supports opentelemetry and currently uses `jaeger` exporter, `stdout`, and `zipkin` exporters are supported as well.
