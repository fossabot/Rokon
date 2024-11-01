# Maintainer: Brycen G <brycengranville@outlook.com>

# This is the 'multiarch' flavor of Rokon. The normal Docker image uses the latest stable version of Fedora Linux. This image as you can see, uses OpenSUSE Tumbleweed. This was created to allow Rokon Docker images to run other architectures like i386. It then uses the tarball which has all the dependencies inside (besides Mesa, sadly) to run anywhere on the platform provided there's Mesa on it. Or once Rokon has command line support, Mesa is not really a hard dependency.

FROM opensuse/tumbleweed:latest AS builder

WORKDIR /app
COPY . .

ENV LANG=C.UTF-8
ENV LC_ALL=C.UTF-8
ENV CC clang
ENV CXX clang++
ENV CFLAGS="-O0 -w -fno-strict-aliasing -gline-tables-only"

RUN zypper in -y git go gtk4-devel gobject-introspection-devel make clang awk



RUN make PACKAGED=true TBPKGFMT=docker NOTB=1 tarball

FROM opensuse/tumbleweed:latest AS runner

ENV LANG=C.UTF-8
ENV LC_ALL=C.UTF-8

WORKDIR /app

COPY --from=builder /app/tarball .


CMD ["./rokon"]
