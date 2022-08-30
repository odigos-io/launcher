FROM gcc:9.2 AS build
ENV DEBIAN_FRONTEND noninteractive
RUN apt-get update && apt-get install -y cmake libgtest-dev libboost-test-dev && rm -rf /var/lib/apt/lists/* 
COPY . /launcher
WORKDIR /launcher
RUN mkdir build && cd build && cmake .. && make

FROM busybox
RUN mkdir /kv-launcher
COPY --from=build /launcher/build/HelloLIEF-prefix/src/HelloLIEF-build/HelloLIEF /kv-launcher/launch
RUN chmod -R go+r /kv-launcher