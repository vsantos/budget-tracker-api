FROM golang:1.16 AS builder

WORKDIR /budget-tracker-api

# manage dependencies
COPY go.mod .
COPY go.sum .
RUN go mod download

COPY config config
COPY controllers controllers
COPY crypt crypt
COPY docs docs
COPY handlers handlers
COPY models models
COPY observability observability
COPY repository repository
COPY routes routes
COPY server server
COPY services services
COPY main.go .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -o app

FROM golang:1.16 AS runner

WORKDIR /app

COPY --from=builder /budget-tracker-api/app /budget-tracker-api/app
COPY --from=builder /budget-tracker-api/docs/swagger.yaml /budget-tracker-api/docs/swagger.yaml

ENTRYPOINT [ "/budget-tracker-api/app" ]
