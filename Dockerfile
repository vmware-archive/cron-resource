FROM golang:1.13 as resource
COPY . /resource
WORKDIR /resource
RUN ./build.sh

FROM ubuntu:bionic
COPY --from=resource /resource/tmp/build/* /opt/resource/
RUN apt update && DEBIAN_FRONTEND=noninteractive apt install -y tzdata && rm -rf /var/lib/apt/lists/*
