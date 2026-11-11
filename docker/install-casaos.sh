#!/bin/bash
# CasaOS Docker Container Install Script

set -e

echo "🏠 Starting CasaOS installation in Docker container..."

# System update
echo "📦 Updating system packages..."
apt-get update -qq

# Install required packages
echo "🔧 Installing required packages..."
apt-get install -y \
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
    git

# Start systemd
echo "⚙️ Starting systemd services..."
systemctl daemon-reload || true

# Start Docker service
echo "🐳 Starting Docker service..."
service docker start || true

# Install CasaOS
echo "🏠 Installing CasaOS..."
if [ ! -f /usr/local/bin/casaos ]; then
    echo "Downloading and installing CasaOS..."
    # Use official installation script
    curl -fsSL https://get.casaos.io | bash
else
    echo "CasaOS is already installed."
fi

# Start CasaOS service
echo "🚀 Starting CasaOS services..."
systemctl enable casaos || true
systemctl start casaos || true

# Port information
echo ""
echo "✅ CasaOS installation completed!"
echo "🌐 CasaOS web interface access:"
echo "   HTTP:  http://localhost"
echo "   HTTPS: https://localhost"
echo ""
echo "📊 Container status check:"
echo "   CasaOS status: systemctl status casaos"
echo "   CasaOS logs:   journalctl -u casaos -f"
echo ""

# Check CasaOS status
sleep 5
if systemctl is-active --quiet casaos; then
    echo "🎉 CasaOS is running successfully!"
else
    echo "⚠️ CasaOS service failed to start. Try starting manually:"
    echo "   docker exec -it casaos-ubuntu systemctl start casaos"
fi

echo "🔄 Container will continue running..."
