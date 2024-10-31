# Maintainer: Brycen G <brycengranville@outlook.com>

# This is the 'multiarch' flavor of Rokon. The normal Docker image uses the latest stable version of Fedora Linux. This image as you can see, uses Debian Testing. This was created to allow Rokon Docker images to run other architectures like i386. It then uses the tarball which has all the dependencies inside (besides Mesa, sadly) to run anywhere on the platform provided there's Mesa on it. Or once Rokon has command line support, Mesa is not really a hard dependency.

FROM debian:testing AS builder

WORKDIR /app
COPY . .

ENV DEBIAN_FRONTEND=noninteractive
ENV LANG=C.UTF-8
ENV LC_ALL=C.UTF-8



RUN apt update
RUN apt install -y build-essential git libgtk-4-dev libgirepository1.0-dev make golang-go
RUN apt clean && apt autoremove




RUN make PACKAGED=true TBPKGFMT=docker NOTB=1 SANITYCHECK=0 tarball

FROM debian:testing AS runner

ENV DEBIAN_FRONTEND=noninteractive
ENV LANG=C.UTF-8
ENV LC_ALL=C.UTF-8

RUN apt update
RUN apt install -y mesa-vulkan-drivers mesa-opencl-icd mesa-vdpau-drivers libegl-mesa0 libgles2



RUN apt clean && apt autoremove


WORKDIR /app

COPY --from=builder /app/tarball .


CMD ["./rokon"]
