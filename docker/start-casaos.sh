#!/bin/bash

# CasaOS Docker Container Starter Script
# This script starts CasaOS as a Docker container

set -e

echo "🏠 Starting CasaOS Docker Container..."
echo "📅 Date: $(date)"
echo ""

# Create necessary directories
echo "📁 Creating necessary directories..."
mkdir -p ./data
mkdir -p ./logs
mkdir -p ./config

# Check if Docker is running
echo "🐳 Checking Docker..."
if ! docker info >/dev/null 2>&1; then
    echo "❌ Docker is not running or permission denied!"
    echo "   On macOS, start Docker Desktop: open -a Docker"
    echo "   Or in terminal: brew services start docker"
    exit 1
fi

# Check if Docker Compose is available
if ! command -v docker-compose >/dev/null 2>&1; then
    echo "❌ docker-compose not found!"
    echo "   To install: brew install docker-compose"
    exit 1
fi

echo "✅ Docker and Docker Compose ready!"
echo ""

# Ask user which configuration to use
echo "Which CasaOS configuration would you like to use?"
echo ""
echo "1) 📦 Ready Ubuntu Image (docker-compose.yml) - Fast"
echo "2) 🔨 Custom Build (docker-compose.build.yml) - Your own image"
echo "3) 🐧 Debian Alternative (--profile debian) - For Debian lovers"
echo "4) 🧽 Simple Configuration (docker-compose.simple.yml) - Minimal"
echo ""
read -p "Make your choice [1-4]: " choice

case $choice in
    1)
        COMPOSE_FILE="docker-compose.yml"
        SERVICE_NAME="casaos-ubuntu"
        echo "🚀 Starting with Ubuntu-based configuration..."
        ;;
    2)
        COMPOSE_FILE="docker-compose.build.yml"
        SERVICE_NAME="casaos-custom"
        echo "🔨 Starting with custom build configuration..."
        echo "⏳ First build may take some time..."
        ;;
    3)
        COMPOSE_FILE="docker-compose.yml"
        SERVICE_NAME="casaos-debian"
        PROFILE_ARG="--profile debian"
        echo "🐧 Starting with Debian-based configuration..."
        ;;
    4)
        COMPOSE_FILE="docker-compose.simple.yml"
        SERVICE_NAME="casaos-simple"
        echo "🧽 Starting with simple configuration..."
        ;;
    *)
        echo "❌ Invalid choice! Using default Ubuntu configuration."
        COMPOSE_FILE="docker-compose.yml"
        SERVICE_NAME="casaos-ubuntu"
        ;;
esac

# Start CasaOS
echo ""
echo "📦 Starting CasaOS container..."
echo "📄 File: $COMPOSE_FILE"
echo "🔧 Service: $SERVICE_NAME"
    2)
        COMPOSE_FILE="docker-compose.build.yml"
        SERVICE_NAME="casaos-custom"
        echo "🔨 Custom build konfigürasyonu ile başlatılıyor..."
        echo "⏳ İlk build işlemi biraz zaman alabilir..."
        ;;
    3)
        COMPOSE_FILE="docker-compose.yml"
        SERVICE_NAME="casaos-debian"
        PROFILE_ARG="--profile debian"
        echo "� Debian tabanlı konfigürasyon ile başlatılıyor..."
        ;;
    4)
        COMPOSE_FILE="docker-compose.simple.yml"
        SERVICE_NAME="casaos-simple"
        echo "🧽 Basit konfigürasyon ile başlatılıyor..."
        ;;
    *)
        echo "❌ Geçersiz seçim! Varsayılan olarak Ubuntu konfigürasyonu kullanılıyor."
        COMPOSE_FILE="docker-compose.yml"
        SERVICE_NAME="casaos-ubuntu"
        ;;
esac

# CasaOS'u başlat
echo ""
echo "📦 CasaOS container'ı başlatılıyor..."
echo "📄 Dosya: $COMPOSE_FILE"
echo "🔧 Servis: $SERVICE_NAME"
echo ""

if [ -n "$PROFILE_ARG" ]; then
    docker-compose -f "$COMPOSE_FILE" $PROFILE_ARG up -d
else
    docker-compose -f "$COMPOSE_FILE" up -d
fi

echo ""
echo "⏳ CasaOS'un başlaması bekleniyor (60 saniye)..."
sleep 10

# Container durumunu kontrol et
echo "🔍 Container durumu kontrol ediliyor..."
if docker ps | grep -q "$SERVICE_NAME"; then
    echo "✅ Container başarıyla çalışıyor!"
    
    # CasaOS kurulumunun tamamlanmasını bekle
    echo "⏳ CasaOS kurulumunun tamamlanması bekleniyor..."
    for i in {1..50}; do
        if docker exec "$SERVICE_NAME" systemctl is-active --quiet casaos 2>/dev/null; then
            echo "🎉 CasaOS servisi aktif!"
            break
        fi
        echo -n "."
        sleep 3
    done
    echo ""
    
else
    echo "❌ Container başlatılamadı!"
    echo "🔍 Hata ayıklama için logları kontrol edin:"
    echo "   docker-compose -f $COMPOSE_FILE logs"
    exit 1
fi

echo ""
echo "🎉 CasaOS başarıyla başlatıldı!"
echo ""
echo "📱 Erişim Bilgileri:"
echo "   🌐 Web Arayüzü: http://localhost"
echo "   🔒 HTTPS:       https://localhost"
echo ""
echo "🔧 Yönetim Komutları:"
echo "   📊 Durumu kontrol et: docker-compose -f $COMPOSE_FILE ps"
echo "   📋 Logları göster:    docker-compose -f $COMPOSE_FILE logs -f $SERVICE_NAME"
echo "   🔄 Yeniden başlat:    docker-compose -f $COMPOSE_FILE restart $SERVICE_NAME"
echo "   ⏹️  Durdur:            docker-compose -f $COMPOSE_FILE down"
echo "   🖥️  Container'a gir:   docker exec -it $SERVICE_NAME bash"
echo ""
echo "🐞 Hata Ayıklama:"
echo "   🔍 CasaOS durumu:     docker exec $SERVICE_NAME systemctl status casaos"
echo "   📝 CasaOS logları:    docker exec $SERVICE_NAME journalctl -u casaos -f"
echo ""
echo "⚠️  Not: CasaOS'un tamamen yüklenmesi 2-5 dakika sürebilir."
echo "    Eğer web arayüzü açılmıyorsa, birkaç dakika bekleyip tekrar deneyin."
