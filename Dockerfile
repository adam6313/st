FROM golang:latest as binder
LABEL stage=builder

COPY . /api
WORKDIR /api

ARG COMMIT
ARG DATE
ARG TAG

# Input the origin package path

ENV GO111MODULE=on

RUN CGO_ENABLED=1 GOOS=linux go build -v -a -installsuffix cgo \
    -ldflags "-X main.VERSION=$TAG -X main.COMMIT=$COMMIT -X main.BUILD=$DATE" \
    -o api

##
FROM ubuntu:21.04
ENV TZ Asia/Taipei


RUN apt update -y && \ 
    apt install -y wget && \ 
    apt install -y gnupg2 && \ 
    GRPC_HEALTH_PROBE_VERSION=v0.3.1 && \
    wget -qO/bin/grpc_health_probe https://github.com/grpc-ecosystem/grpc-health-probe/releases/download/${GRPC_HEALTH_PROBE_VERSION}/grpc_health_probe-linux-amd64 && \
    chmod +x /bin/grpc_health_probe

EXPOSE 5000

COPY --from=binder /api/api /api

Add watermark.png /watermark.png


WORKDIR /
ENTRYPOINT ["./api"]
CMD  ["version"]
