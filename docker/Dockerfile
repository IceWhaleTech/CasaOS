# CasaOS Dockerfile
# Ubuntu 22.04 based CasaOS container

FROM ubuntu:22.04

# Metadata
LABEL maintainer="CasaOS User"
LABEL description="CasaOS Personal Cloud in Docker"
LABEL version="1.0"

# Set timezone
ENV TZ=Europe/Istanbul
ENV DEBIAN_FRONTEND=noninteractive

# Install system packages
RUN apt-get update && apt-get install -y \
    curl \
    wget \
    sudo \
    systemd \
    systemd-sysv \
    init \
    ca-certificates \
    gnupg \
    lsb-release \
    docker.io \
    docker-compose \
    openssl \
    git \
    net-tools \
    htop \
    nano \
    && rm -rf /var/lib/apt/lists/*

# Copy CasaOS installation script
COPY install-casaos.sh /usr/local/bin/install-casaos.sh
RUN chmod +x /usr/local/bin/install-casaos.sh

# Create necessary directories
RUN mkdir -p /var/lib/casaos \
             /etc/casaos \
             /usr/share/casaos \
             /DATA \
             /var/log/casaos

# Install CasaOS
RUN /usr/local/bin/install-casaos.sh || true

# Configure systemd services
RUN systemctl enable systemd-networkd
RUN systemctl enable systemd-resolved

# Expose ports
EXPOSE 80 443

# Volume mount points
VOLUME ["/var/lib/casaos", "/etc/casaos", "/usr/share/casaos", "/DATA"]

# Start systemd
CMD ["/sbin/init"]

# Healthcheck
HEALTHCHECK --interval=30s --timeout=10s --start-period=60s --retries=3 \
    CMD curl -f http://localhost/ || exit 1
