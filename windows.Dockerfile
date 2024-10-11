ARG WINDOWS_VERSION=latest
FROM amitie10g/msys2:${WINDOWS_VERSION} AS builder

ENV MSYSTEM=CLANG64

RUN pacman -S --needed --noprogressbar gtk4 git mingw-w64-clang-x86_64-go mingw-w64-clang-x86_64-gtk4 mingw-w64-clang-x86_64-upx mingw-w64-clang-x86_64-gobject-introspection mingw-w64-clang-x86_64-gdb mingw-w64-clang-x86_64-toolchain make

ENV PATH="C:\\msys64\\clang64\\bin;${PATH}"


WORKDIR C:\app

COPY . .


RUN make TARGET="Rokon.exe" PACKAGED=true EXTRALDFLAGS="-s -w -H windowsgui" EXTRAGOFLAGS="-trimpath" PACKAGEFORMAT=docker build
RUN make TARGET="Rokon.exe" PREFIX="./Rokon" BINDIR="./Rokon" install




RUN upx -f --best --force-overwrite ./*.exe


CMD ["./Rokon.exe"]

