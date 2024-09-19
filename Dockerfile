FROM alpine:edge AS builder

RUN apk add --no-cache alpine-sdk go gtk4.0-dev gobject-introspection-dev

WORKDIR /app

COPY . .

RUN make build

FROM alpine:edge AS runner

WORKDIR /app

RUN apk add --no-cache gtk4.0 gobject-introspection mesa-gles


COPY --from=builder /app/rokon .

# Run the application

CMD ["./rokon"]
