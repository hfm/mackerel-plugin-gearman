#!/bin/bash
#
# Original from: http://deeeet.com/writing/2014/05/19/gox/

set -e

DIR=$(cd $(dirname ${0})/.. && pwd)
pushd ${DIR}

VERSION=$(grep "const Version " version.go | sed -E 's/.*"(.+)"$/\1/')
REPO="mackerel-plugin-gearman"

# Run Compile
./scripts/compile.sh

if [ -d pkg ];then
    rm -rf ./pkg/dist
fi

# Package all binary as .zip and .tar.xz
mkdir -p ./pkg/dist/${VERSION}
for PLATFORM in $(find ./pkg -mindepth 1 -maxdepth 1 -type d); do
    PLATFORM_NAME=$(basename ${PLATFORM})
    ARCHIVE_NAME=${REPO}_${VERSION}_${PLATFORM_NAME}

    if [ $PLATFORM_NAME = "dist" ]; then
        continue
    fi

    pushd ${PLATFORM}
    zip ${DIR}/pkg/dist/${VERSION}/${ARCHIVE_NAME}.zip ./*
    tar Jcf ${DIR}/pkg/dist/${VERSION}/${ARCHIVE_NAME}.tar.xz ./*
    popd
done

# Generate shasum
pushd ./pkg/dist/${VERSION}
shasum * > ./${VERSION}_SHASUMS
popd
popd
