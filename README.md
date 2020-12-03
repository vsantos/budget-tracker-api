# budget-tracker-api

API REST Backend for a personal budget tracker

## Architecture

<img src="imgs/budget-tracker.png" alt="">

# Developer tools

## Running locally

You can use `docker-compose` to run the entire backend stack locally: `budget-tracker` and `mongodb` (with an initial `admin` user created)

The mongodb served by `docker-compose` has no credentials so it's recommended only for development purposes.

## Swagger API support

This application uses go-swagger to generate swagger specs directly from the code, to run it just:

`make swagger-spec`

The output will be a local `swagger.yaml` file which can be used along with external applications such as insomnia designer. In case of side-by-side UI containers such as `swagger-ui` you can serve the following endpoint: `/swagger.yaml`