FROM golang:1.13 AS builder
WORKDIR /app
COPY . /app/
RUN go get
RUN go build -o budget-tracker

FROM golang:1.13 AS runner
WORKDIR /app

COPY --from=builder /app/budget-tracker /app/budget-tracker
COPY --from=builder /app/swagger.yaml /app/swagger.yaml

ENTRYPOINT [ "/app/budget-tracker" ]