# Maintainer: Ns2Kracy <2220496937@qq.com>
pkgname=casaos
pkgver=0.4.0
pkgrel=1
pkgdesc="Community-based open source software focused on delivering simple home cloud experience around Docker ecosystem."
arch=('x86_64' 'aarch64' 'armv7h')
url="https://github.com/IceWhaleTech/CasaOS"
license=('APACHE')
depends=(
    'smartmontools' 'parted' 'ntfs-3g' 'net-tools' 'udevil' 'samba' 'cifs-utils' 'docker' 'docker-compose'
    'casaos-cli' 'casaos-ui' 'casaos-app-management' 'casaos-local-storage' 'casaos-user-service' 'casaos-gateway' 'casaos-message-bus'
    )
groups=('casaos')
install="${pkgname}.install"
backup=('etc/casaos/casaos.conf')
source_x86_64=(
    ${url}/releases/download/v${pkgver}/linux-amd64-${pkgname}-v${pkgver}.tar.gz
    ${url}/releases/download/v${pkgver}/linux-amd64-${pkgname}-migration-tool-v${pkgver}.tar.gz
    )
source_aarch64=(
    ${url}/releases/download/v${pkgver}/linux-arm64-${pkgname}-v${pkgver}.tar.gz
    ${url}/releases/download/v${pkgver}/linux-arm64-${pkgname}-migration-tool-v${pkgver}.tar.gz
    )
source_armv7h=(
    ${url}/releases/download/v${pkgver}/linux-arm-7-${pkgname}-v${pkgver}.tar.gz
    ${url}/releases/download/v${pkgver}/linux-arm-7-${pkgname}-migration-tool-v${pkgver}.tar.gz
    )
sha256sums_x86_64=(
    dc30edc8bc69da5db3b2a5943d097ae59d205c1c1fd67fd58344cbb8b5d9abb1
    d6a2670d24b7934ed08815933e602133c356c3a57e0ba6e6f0f59311a5f7f12c
    )
sha256sums_aarch64=(
    48b83900e2d03d62c08b629bf52a062505739d4a85825067fd4a7a21b3a2ec4f 
    3ce1196ec39da92707acfc059161ae0726cbc646ce68fc2e7b2aa51b0e393ff1 
    )
sha256sums_armv7h=(
    747f0f68b374ed7f32de90fbca655f2c6644751e418af3302c4e2200ca955b61
    db71bb5a475c22d2308af34c6310269103ea708739f9de2949444f85a7b048d0
    )

package() {
    _sysdir="${srcdir}/build/sysroot"
	_name="${pkgname#*-}"
	install -Dm755 "${_sysdir}/usr/bin/${pkgname}" "${pkgdir}/usr/bin/${pkgname}"
	install -Dm755 "${_sysdir}/usr/bin/${pkgname}-migration-tool" "${pkgdir}/usr/bin/${pkgname}-migration-tool"
	install -Dm644 "${_sysdir}/etc/casaos/${_name}.conf.sample" "${pkgdir}/etc/casaos/${_name}.conf"
	install -Dm644 "${_sysdir}/usr/lib/systemd/system/${pkgname}.service" "${pkgdir}/usr/lib/systemd/system/${pkgname}.service"
}
