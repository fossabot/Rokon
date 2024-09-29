#! /bin/bash
# Build the docker image with the following command:
# docker build -t rokon .
# Using podman? Replace "docker" with "podman" in the command above.
read Xenv < <(x11docker --xvfb --gpu --wayland --printenv rokon /app/rokon)
echo $Xenv && export $Xenv
# replace "start" with "start-desktop" to forward a desktop environment
xpra start $DISPLAY --use-display \
     --html=on --bind-tcp=localhost:14501 \
     --start-via-proxy=no
