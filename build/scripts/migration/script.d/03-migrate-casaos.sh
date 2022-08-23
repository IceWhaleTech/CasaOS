#!/bin/bash

set -e

# functions
__is_version_gt() {
    test "$(echo "$@" | tr " " "\n" | sort -V | head -n 1)" != "$1"
}

__is_migration_needed() {
    local version1
    local version2

    version1="${1}"
    version2="${2}"

    if [ "${version1}" = "${version2}" ]; then
        return 1
    fi

    if [ "CURRENT_VERSION_NOT_FOUND" = "${version1}" ]; then
        return 1
    fi

    if [ "LEGACY_WITHOUT_VERSION" = "${version1}" ]; then
        return 0
    fi

    __is_version_gt "${version2}" "${version1}"
}

BUILD_PATH=$(dirname "${BASH_SOURCE[0]}")/../../..
SOURCE_ROOT=${BUILD_PATH}/sysroot

APP_NAME="casaos"
APP_NAME_FORMAL="CasaOS"

# check if migration is needed
SOURCE_BIN_PATH=${SOURCE_ROOT}/usr/bin
SOURCE_BIN_FILE=${SOURCE_BIN_PATH}/${APP_NAME}

CURRENT_BIN_PATH=/usr/bin
CURRENT_BIN_PATH_LEGACY=/usr/local/bin
CURRENT_BIN_FILE=${CURRENT_BIN_PATH}/${APP_NAME}
CURRENT_BIN_FILE_LEGACY=$(realpath -e ${CURRENT_BIN_PATH_LEGACY}/${APP_NAME} || which ${APP_NAME} || echo CURRENT_BIN_FILE_LEGACY_NOT_FOUND)

SOURCE_VERSION="$(${SOURCE_BIN_FILE} -v)"
CURRENT_VERSION="$(${CURRENT_BIN_FILE} -v || ${CURRENT_BIN_FILE_LEGACY} -v || (stat "${CURRENT_BIN_FILE_LEGACY}" > /dev/null && echo LEGACY_WITHOUT_VERSION) || echo CURRENT_VERSION_NOT_FOUND)"

echo "CURRENT_VERSION: ${CURRENT_VERSION}"
echo "SOURCE_VERSION: ${SOURCE_VERSION}"

NEED_MIGRATION=$(__is_migration_needed "${CURRENT_VERSION}" "${SOURCE_VERSION}" && echo "true" || echo "false")

if [ "${NEED_MIGRATION}" = "false" ]; then
    echo "✅ Migration is not needed."
    exit 0
fi

MIGRATION_SERVICE_DIR=${BUILD_PATH}/scripts/migration/service.d/${APP_NAME}
MIGRATION_LIST_FILE=${MIGRATION_SERVICE_DIR}/migration.list
MIGRATION_PATH=()

CURRENT_VERSION_FOUND="false"

while read -r VERSION_PAIR; do
    if [ -z "${VERSION_PAIR}" ]; then
        continue
    fi

    VER1=$(echo "${VERSION_PAIR}" | cut -d' ' -f1)
    VER2=$(echo "${VERSION_PAIR}" | cut -d' ' -f2)

    if [ "v${CURRENT_VERSION}" = "${VER1// /}" ]; then
        CURRENT_VERSION_FOUND="true"
    fi

    if [ "${CURRENT_VERSION_FOUND}" = "true" ]; then
        MIGRATION_PATH+=("${VER2// /}")
    fi
done < "${MIGRATION_LIST_FILE}"

if [ ${#MIGRATION_PATH[@]} -eq 0 ]; then
    echo "🟨 No migration path found from ${CURRENT_VERSION} to ${SOURCE_VERSION}"
    exit 0
fi

ARCH="unknown"

case $(uname -m) in
    x86_64)
        ARCH="amd64"
        ;;
    aarch64)
        ARCH="arm64"
        ;;
    armv7l)
        ARCH="arm-7"
        ;;
    *)
        echo "Unsupported architecture"
        exit 1
        ;;
esac

pushd "${MIGRATION_SERVICE_DIR}"

{
    for VER2 in "${MIGRATION_PATH[@]}"; do
        MIGRATION_TOOL_URL=https://github.com/IceWhaleTech/"${APP_NAME_FORMAL}"/releases/download/"${VER2}"/linux-"${ARCH}"-"${APP_NAME}"-migration-tool-"${VER2}".tar.gz
        echo "Dowloading ${MIGRATION_TOOL_URL}..."
        curl -sL -O "${MIGRATION_TOOL_URL}"
    done
} || {
    echo "🟥 Failed to download migration tools"
    popd
    exit 1
}

{
    for VER2 in "${MIGRATION_PATH[@]}"; do
        MIGRATION_TOOL_FILE=linux-"${ARCH}"-"${APP_NAME}"-migration-tool-"${VER2}".tar.gz
        echo "Extracting ${MIGRATION_TOOL_FILE}..."
        tar zxvf "${MIGRATION_TOOL_FILE}"

        MIGRATION_TOOL_PATH=build/sysroot/usr/bin/${APP_NAME}-migration-tool
        echo "Running ${MIGRATION_TOOL_PATH}..."
        ${MIGRATION_TOOL_PATH}
    done
} || {
    echo "🟥 Failed to extract and run migration tools"
    popd
    exit 1
}

popd
