# workspace (GOPATH) configured at /go
FROM golang:1.18 as builder

WORKDIR /app

# Copy the local package files to the container's workspace.
COPY . ./

RUN make build-client

FROM alpine:latest

COPY --from=builder app/bin/client ./client

RUN chmod +x ./client 

CMD ["./client"]
