FROM alpine:edge AS builder

RUN apk add --no-cache alpine-sdk go gtk4.0-dev gobject-introspection-dev bash

WORKDIR /app

COPY . .

<<<<<<< HEAD
RUN go build -v -o rokon .
=======
RUN make PACKAGED=true PACKAGEFORMAT=docker build
>>>>>>> eb15601 (docs(docker): NVIDIA GPUs are not supported with this container)

FROM alpine:edge AS runner

WORKDIR /app

# Nvidia GPUs are NOT supported with this container!

RUN apk add --no-cache gtk4.0 gobject-introspection mesa-gles


COPY --from=builder /app/rokon .

# Run the application

CMD ["./rokon"]
