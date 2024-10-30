# Maintainer: Brycen G <brycengranville@outlook.com>

# This container *can* work under NVIDIA.
# For example, you could run the command `distrobox create --name rokon --image ghcr.io/brycensranch/rokon --nvidia`

FROM fedora:41 AS builder

RUN dnf install -y \
    make \
    go \
    gtk4-devel \
    gobject-introspection-devel \
    which \
    patchelf \
    upx
RUN dnf clean all



WORKDIR /app
COPY . .

# TBPKGFMT = TARBALL PACKAGE FORMAT (This is used for telemetry and logging purposes, does not affect the package whatsoever)
# NOTB = Prevents the creation of tar.gz files. It's not needed and the container won't use it.
RUN make PACKAGED=true TBPKGFMT=docker NOTB=1 tarball

FROM fedora:41 AS runner

WORKDIR /app

COPY --from=builder /app/tarball .

CMD ["./rokon"]
