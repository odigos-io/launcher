FROM golang:1.19 AS build
WORKDIR /app
COPY . .
RUN apt-get update && apt-get install -y nasm
RUN nasm payload/mmap.asm
RUN CGO_ENABLED=0 go build -o launcher

FROM busybox
WORKDIR /kv-launcher
COPY --from=build /app/launcher /kv-launcher/launch
RUN chmod -R go+r /kv-launcher