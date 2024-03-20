FROM golang:1.22-alpine as builder

WORKDIR /app

COPY . .

RUN go mod download

RUN go build -o job-portal cmd/main.go

FROM alpine:3.15

RUN apk --no-cache add ca-certificates curl

COPY --from=builder /app/job-portal /
COPY --from=builder /app/config.json /config.json

ENTRYPOINT [ "/job-portal" ]

CMD [ "start", "-c", "/config.json" ]