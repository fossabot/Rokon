# Why is boss music playing?
FROM fedora:rawhide

# We need to create a non-root user

# Install the necessary packages

RUN dnf install -y golang gtk4-devel gobject-introspection-devel git @development-tools @development-libs && dnf clean all -y

# Set the working directory

RUN useradd -m -s /bin/bash dev
RUN echo "dev ALL=(ALL) NOPASSWD: ALL" > /etc/sudoers.d/dev
# add them to the wheel and sudo groups
RUN usermod -aG wheel dev
RUN usermod -aG sudo dev

USER dev

WORKDIR /app

# Copy the source code

COPY --chown=dev:dev . .


# Drop developer into their development environment

CMD ["/usr/bin/bash"]
