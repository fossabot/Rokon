FROM alpine:edge AS builder



# Install the necessary packages

RUN apk add --no-cache alpine-sdk go gtk4.0-dev gobject-introspection-dev
# Set the working directory

WORKDIR /app

# Copy the source code

COPY . .

# Build the application

RUN go build -v -o rokon .

FROM alpine:edge AS runner

# Copy the binary from the builder image

WORKDIR /app

RUN apk add --no-cache gtk4.0 gobject-introspection mesa-gles


COPY --from=builder /app/rokon .

# Run the application

CMD ["./rokon"]
