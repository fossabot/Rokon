# Maintainer: Brycen G <brycengranville@outlook.com>

# This is the 'multiarch' flavor of Rokon. The normal Docker image uses the latest stable version of Fedora Linux. This image as you can see, it uses Debian Testing. This was created to allow Rokon Docker images to run other architectures like i386 while keeping the dependencies up to date as possible since this does not rely on the host's libraries. It then uses the tarball which has all the dependencies inside (besides Mesa, sadly) to run anywhere on the platform provided there's Mesa on it. Or once Rokon has command line support, Mesa is not really a hard dependency.

FROM debian:testing AS builder

WORKDIR /app
COPY . .

# Breaks building on armhf!
RUN rm cgosymbolizer_linux.go || true

ENV DEBIAN_FRONTEND=noninteractive
ENV LANG=C.UTF-8
ENV LC_ALL=C.UTF-8
ENV CC=clang
ENV CXX=clang++
ENV CFLAGS="-O0 -w -fno-strict-aliasing -gline-tables-only"

RUN apt update && apt full-upgrade -y
RUN apt install -y git libgtk-4-dev libgirepository1.0-dev make golang-go clang
RUN apt clean && apt autoremove


RUN make PACKAGED=true TBPKGFMT=docker NOTB=1 tarball

FROM debian:testing-slim AS runner


ENV DEBIAN_FRONTEND=noninteractive
ENV LANG=C.UTF-8
ENV LC_ALL=C.UTF-8

WORKDIR /app

# libGLES is a hard dependency for the GUI.
# Without ca-certitifcates, telemetry fails to send.
RUN apt update && apt full-upgrade -y && apt install -y libgles2 ca-certificates && apt clean && apt autoremove

COPY --from=builder /app/tarball .


CMD ["./rokon"]
