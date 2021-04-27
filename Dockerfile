FROM golang:1.16 AS builder
WORKDIR /app
COPY . /app/
RUN go get
RUN go build -o budget-tracker

FROM golang:1.16 AS runner
WORKDIR /app

COPY --from=builder /app/budget-tracker /app/budget-tracker
COPY --from=builder /app/docs/swagger.yaml /app/docs/swagger.yaml

ENTRYPOINT [ "/app/budget-tracker" ]
