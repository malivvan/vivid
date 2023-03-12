#!/usr/bin/env bash

PACKAGE="github.com/malivvan/vivid"
PUBKEY="MCowBQYDK2VwAyEAVMhMRJeVvkUlmtpAG2aVJpoIWXASCHbKnakqttnJ93U="
VERSION="$(git describe --tags --always --abbrev=0 --match='v[0-9]*.[0-9]*.[0-9]*' 2> /dev/null | sed 's/^.//')"
COMMIT_HASH="$(git rev-parse --short HEAD)"
BUILD_TIMESTAMP=$(date '+%Y-%m-%dT%H:%M:%S')
LDFLAGS=(
  "-s"
  "-w"
  "-X '${PACKAGE}.AppRepo=${PACKAGE}'"
  "-X '${PACKAGE}.AppVersion=${VERSION}'"
  "-X '${PACKAGE}.AppCommit=${COMMIT_HASH}'"
  "-X '${PACKAGE}.AppBuild=${BUILD_TIMESTAMP}'"
  "-X '${PACKAGE}.AppPubkey=${PUBKEY}'"
)
PLATFORMS=(
 # "windows arm64 vivid_windows_arm64.exe"
 # "windows arm vivid_windows_arm.exe"
 # "windows 386 vivid_windows_386.exe"
  "windows amd64 vivid_windows_amd64.exe"
#  "linux arm64 vivid_linux_arm64"
 # "linux arm vivid_linux_arm"
 # "linux 386 vivid_linux_386"
 # "linux amd64 vivid_linux_amd64"
)

for platform in "${PLATFORMS[@]}"; do
    IFS=" " read -a array <<< $platform
    OUTPUT=${array[2]}
    echo "Building ${OUTPUT}"
    GOOS=${array[0]} GOARCH=${array[1]} go build -o ../dist/${OUTPUT} -ldflags="${LDFLAGS[*]}" -trimpath
done


