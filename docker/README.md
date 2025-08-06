# Running CasaOS in Docker Containers

This repository provides Docker configurations to run [CasaOS](https://github.com/IceWhaleTech/CasaOS) in containers on any system that supports Docker.

## 🐳 Why Docker?

CasaOS is primarily designed for Linux systems, but many users want to experience this powerful personal cloud OS on different platforms or in isolated environments. This Docker approach allows you to:

- ✅ Run CasaOS on Windows, macOS, and Linux
- ✅ Isolate CasaOS in containers for safety and testing
- ✅ Easy setup, backup, and teardown
- ✅ Multiple configuration options
- ✅ No need to modify your host system

## 📋 Prerequisites

### Required Software

- **Docker**: [Install Docker](https://docs.docker.com/get-docker/)
- **Docker Compose**: [Install Docker Compose](https://docs.docker.com/compose/install/)

### Hardware Requirements

- 4GB RAM minimum (8GB recommended)
- 20GB free disk space
- x86_64 or ARM64 architecture

### Platform-Specific Installation

#### Windows

```bash
# Install Docker Desktop
# Download from: https://docs.docker.com/desktop/windows/install/
```

#### macOS

```bash
# Install Docker Desktop
brew install --cask docker

# Install Docker Compose (usually included)
brew install docker-compose
```

#### Linux

```bash
# Install Docker
curl -fsSL https://get.docker.com -o get-docker.sh
sh get-docker.sh

# Install Docker Compose
sudo curl -L "https://github.com/docker/compose/releases/download/v2.20.0/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose
sudo chmod +x /usr/local/bin/docker-compose
```

## 🚀 Quick Start

### 1. Clone or Download

```bash
git clone <this-repository>
cd casaos-docker
```

### 2. Start CasaOS (Interactive)

```bash
chmod +x start-casaos.sh
./start-casaos.sh
```

### 3. Manual Start Options

```bash
# Option 1: Ubuntu-based (recommended)
docker-compose up -d

# Option 2: Debian-based
docker-compose --profile debian up -d

# Option 3: Custom build
docker-compose -f docker-compose.build.yml up -d

# Option 4: Simple configuration
docker-compose -f docker-compose.simple.yml up -d
```

## 🌐 Accessing CasaOS

After starting, CasaOS will be available at:

- **HTTP**: <http://localhost>
- **HTTPS**: <https://localhost>

> **Note**: For simple configuration users:
>
> - HTTP: <http://localhost:8080>
> - HTTPS: <https://localhost:8443>

> **Initial Setup**: First startup may take 2-5 minutes. If the web interface doesn't load immediately, wait a few minutes and refresh.

## 📁 Configuration Files

| File | Description |
|------|-------------|
| `docker-compose.yml` | Main configuration (Ubuntu + Debian profiles) |
| `docker-compose.simple.yml` | Simplified configuration with bridge networking |
| `docker-compose.build.yml` | Custom build configuration |
| `Dockerfile` | Custom CasaOS image definition |
| `install-casaos.sh` | Container installation script |
| `start-casaos.sh` | Interactive startup script |

## 🔧 Management Commands

### Container Management

```bash
# Check status
docker-compose ps

# View logs
docker-compose logs -f casaos-ubuntu

# Restart container
docker-compose restart casaos-ubuntu

# Stop everything
docker-compose down

# Enter container
docker exec -it casaos-ubuntu bash
```

### CasaOS Management

```bash
# Check CasaOS service status
docker exec casaos-ubuntu systemctl status casaos

# View CasaOS logs
docker exec casaos-ubuntu journalctl -u casaos -f

# Restart CasaOS service
docker exec casaos-ubuntu systemctl restart casaos
```

## 💾 Data Persistence

Your data is stored in Docker volumes and local directories:

- **Docker Volumes**: `casaos_data`, `casaos_config`, `casaos_share`
- **Local Directories**: `./data`, `./logs`, `./config`

## 🐞 Troubleshooting

### Common Issues

#### Docker not running

```bash
# Check Docker status
docker info

# Start Docker service (Linux)
sudo systemctl start docker

# Start Docker Desktop (Windows/macOS)
# Use Docker Desktop application
```

#### Port conflicts

```bash
# Check what's using port 80
# Linux/macOS:
sudo lsof -i :80
# Windows:
netstat -ano | findstr :80

# Use simple configuration with different ports
docker-compose -f docker-compose.simple.yml up -d
# Then access via http://localhost:8080
```

#### Container won't start

```bash
# Check Docker logs
docker-compose logs casaos-ubuntu

# Check system resources
docker system df
docker system prune  # Clean up if needed
```

#### CasaOS web interface not loading

```bash
# Wait longer (up to 5 minutes for first start)
# Check if CasaOS service is running
docker exec casaos-ubuntu systemctl status casaos

# Check container logs
docker exec casaos-ubuntu journalctl -u casaos -f
```

### Platform-Specific Tips

#### Windows

- Ensure WSL2 is properly configured for Docker Desktop
- Check Windows Defender firewall settings
- Use PowerShell or Command Prompt as Administrator if needed

#### macOS

- Allocate more resources to Docker Desktop (Settings → Resources)
- Close unnecessary applications to free up system resources
- Ensure Docker Desktop is allowed in Security & Privacy settings

#### Linux

- Ensure your user is in the docker group: `sudo usermod -aG docker $USER`
- Check SELinux/AppArmor if containers fail to start
- Verify sufficient disk space for Docker volumes

## 🔒 Security Considerations

- Container runs in privileged mode (required for CasaOS functionality)
- Docker socket is mounted (needed for container management)
- Host networking is used for full CasaOS functionality
- Consider firewall rules if exposing to network

## 📈 Monitoring

### Container Resource Usage

```bash
# CPU and RAM usage
docker stats casaos-ubuntu

# Disk usage
docker exec casaos-ubuntu df -h
```

### CasaOS System Monitoring

CasaOS web interface provides built-in system monitoring capabilities.

## 🌟 Features Tested

- ✅ Web interface access
- ✅ App store and Docker app installation
- ✅ File management
- ✅ System monitoring
- ✅ User management
- ⚠️ Hardware monitoring (limited in containers)
- ⚠️ Some system-level features may be restricted

## 🤝 Contributing

This is a community contribution to help users run CasaOS in Docker containers across different platforms. Contributions are welcome!

### To Contribute

1. Fork this repository
2. Test your changes on your platform
3. Submit a pull request with clear description
4. Include platform-specific considerations if any

## 🎯 About This Project

I created this Docker configuration because I wanted to test and use CasaOS on different platforms without modifying my host system. Since CasaOS doesn't provide official Docker images, this solution creates containers that install CasaOS automatically.

**Use cases:**

- Testing CasaOS before committing to a full installation
- Running CasaOS on non-Linux systems
- Isolating CasaOS for development or testing
- Easy backup and restoration of CasaOS instances

Feel free to use, modify, and improve this configuration for your needs!

## 📖 Related Links

- [CasaOS Official Repository](https://github.com/IceWhaleTech/CasaOS)
- [CasaOS Documentation](https://wiki.casaos.io/)
- [CasaOS Discord Community](https://discord.gg/knqAbbBbeX)
- [Docker Documentation](https://docs.docker.com/)

## 📄 License

This Docker configuration is provided under MIT license. CasaOS itself has its own license - see the [official repository](https://github.com/IceWhaleTech/CasaOS/blob/main/LICENSE).

---

**⚠️ Important Notes:**

- This is an unofficial Docker configuration for CasaOS
- CasaOS doesn't officially support running in Docker containers
- Some features may not work exactly as on native Linux installations
- This configuration is for development, testing, and personal use
- For production use, consider a dedicated Linux system
