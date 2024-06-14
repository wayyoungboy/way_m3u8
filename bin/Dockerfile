FROM golang:1.22.4-bullseye
LABEL authors="wayyoungboy"
WORKDIR /build_path
COPY ../*  /build_path/
RUN go build main.go -o ./bin/way_m3u8
WORKDIR /build_path/bin
ENTRYPOINT ["./way_m3u8" ]