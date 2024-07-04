FROM ubuntu:22.04
LABEL authors="wayyoungboy"
WORKDIR /build_path
COPY .  /build_path/
WORKDIR /build_path/
RUN mkdir video
ENTRYPOINT ["./gom3u8_linux" ]