FROM ubuntu:bionic
ADD built-check /opt/resource/check
ADD built-in /opt/resource/in

RUN apt update && DEBIAN_FRONTEND=noninteractive apt install -y tzdata && rm -rf /var/lib/apt/lists/*
