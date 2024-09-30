# Use a base image with GTK4 support
FROM ubuntu:latest

# Install necessary dependencies
RUN apt update && apt install -y \
    libgtk-4-dev \
    xvfb \
    scrot \
    dbus-x11 \
    software-properties-common \
    libgirepository1.0-dev \
    ffmpeg \
    libgl1-mesa-dev \
    libosmesa6-dev

RUN add-apt-repository -y ppa:longsleep/golang-backports
RUN apt update && apt install golang -y

# Add your application code
WORKDIR /app
COPY . .

# Build the application
RUN go build -o rokon main.go

ENV DISPLAY=:99

# Command to run the application inside Xvfb and take a screenshot
CMD Xvfb :99 -screen 0 1920x1080x24 & \
    export DISPLAY=:99 && \
    dbus-launch ./rokon & \
    sleep 5 && \
    scrot /app/screenshots/desktop.png
