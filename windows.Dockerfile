# Maintainer: Brycen G <brycengranville@outlook.com>

# This is the 'windows' flavor of Rokon. It is more specifically for Windows Pro Edition+. If you're running Windows 10 Home, you should use the normal Rokon container (Linux, runs on Windows Home via WSL2)

ARG WINDOWS_VERSION=latest

FROM amitie10g/msys2:latest AS builder
ENV MSYSTEM=CLANG64

# SHELL ["C:\\msys64\\msys2_shell.cmd", "-defterm", "-clang64", "-no-start", "-here"]
ENV CHERE_INVOKING=1
ENV GOROOT="C:\\msys64\\clang64\\lib\\go"


WORKDIR C:\\app

COPY . .


SHELL ["C:\\msys64\\usr\\bin\\env.exe", "/usr/bin/bash", "--login",  "C:\\app\\windows\\wrapper.sh"]


RUN pacman -S --needed --noconfirm --noprogressbar git mingw-w64-clang-x86_64-go mingw-w64-clang-x86_64-gtk4 mingw-w64-clang-x86_64-upx mingw-w64-clang-x86_64-gobject-introspection mingw-w64-clang-x86_64-gdb mingw-w64-clang-x86_64-toolchain make
RUN ls -R

RUN make TARGET="Rokon.exe" PACKAGED=true EXTRALDFLAGS="-s -w -H windowsgui" EXTRAGOFLAGS="-trimpath" PACKAGEFORMAT=docker build
RUN make TARGET="Rokon.exe" PREFIX="./Rokon" BINDIR="./Rokon" install
RUN ldd "Rokon.exe" | { grep "=> /clang64/bin/" || true; }             | cut -d ' ' -f1             | xargs -I{} cp /clang64/bin/{} ./Rokon

WORKDIR C:\\app\\Rokon


RUN upx -f --best --force-overwrite ./*.exe ./*.dll

FROM mcr.microsoft.com/windows/server:ltsc2022 AS runner


COPY --from=builder C:\\app\\Rokon .


WORKDIR C:\\app\\Rokon

CMD ["Rokon.exe"]

