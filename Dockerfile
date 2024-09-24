FROM alpine:edge AS builder

RUN apk add --no-cache alpine-sdk go gtk4.0-dev gobject-introspection-dev bash

WORKDIR /app

COPY . .

RUN go build -v -o rokon .

FROM alpine:edge AS runner

WORKDIR /app

RUN apk add --no-cache gtk4.0 gobject-introspection mesa-gles


COPY --from=builder /app/rokon .

# Run the application

CMD ["./rokon"]
