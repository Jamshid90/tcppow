# workspace (GOPATH) configured at /go
FROM golang:1.18 as builder

WORKDIR /app

# Copy the local package files to the container's workspace.
COPY . ./

RUN make build-server

FROM alpine:latest

COPY --from=builder app/bin/server ./server

RUN chmod +x ./server 

CMD ["./server"]
