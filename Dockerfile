FROM alpine:edge AS builder

RUN apk add --no-cache alpine-sdk go gtk4.0-dev gobject-introspection-dev git bash

WORKDIR /app

COPY . .

RUN make PACKAGED=true PACKAGEFORMAT=docker build

FROM alpine:edge AS runner

WORKDIR /app

# Nvidia GPUs are NOT supported with this container!

RUN apk add --no-cache gtk4.0 gobject-introspection mesa-gles


COPY --from=builder /app/rokon .

# Run the application

CMD ["./rokon"]
