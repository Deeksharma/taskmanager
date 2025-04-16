FROM golang:1.23.8-alpine as builder

RUN mkdir /workspace
WORKDIR /workspace

COPY go.mod ./
COPY go.sum ./

RUN go mod download

COPY . ./

RUN go build -o server cmd/taskmanager/main.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates
RUN apk add bash
RUN apk add --no-cache openssh-client ansible git

WORKDIR /root/
RUN chmod -R 755 /root

COPY --from=builder /workspace .

ARG DEFAULT_PORT=80
ENV PORT $DEFAULT_PORT

EXPOSE $PORT

CMD ["./server"]
