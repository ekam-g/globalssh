#!/bin/bash

# Define the release file name and the project name
RELEASE_FILE="releases"
PROJECT_NAME="global_ssh"

# Delete Old File
rm -r $RELEASE_FILE

# Create the release file
mkdir $RELEASE_FILE

# Get the list of supported platforms
PLATFORMS=("aix/ppc64" "android/386" "android/amd64" "android/arm" "android/arm64" "darwin/386" "darwin/amd64" "darwin/arm" "darwin/arm64" "dragonfly/amd64" "freebsd/386" "freebsd/amd64" "freebsd/arm" "freebsd/arm64" "illumos/amd64" "ios/arm64" "js/wasm" "linux/386" "linux/amd64" "linux/arm" "linux/arm64" "linux/mips" "linux/mips64" "linux/mips64le" "linux/mipsle" "linux/ppc64" "linux/ppc64le" "linux/riscv64" "linux/s390x" "netbsd/386" "netbsd/amd64" "netbsd/arm" "netbsd/arm64" "openbsd/386" "openbsd/amd64" "openbsd/arm" "openbsd/arm64" "plan9/386" "plan9/amd64" "solaris/amd64" "windows/386" "windows/amd64" "windows/arm" "windows/arm64")
# Loop through the platforms and build the project
for PLATFORM in "${PLATFORMS[@]}"
do
  PLATFORM_SPLIT=(${PLATFORM//\// })
  OS=${PLATFORM_SPLIT[0]}
  ARCH=${PLATFORM_SPLIT[1]}

  # Set the output name based on the platform
  if [ "$OS" == "windows" ]; then
    OUTPUT_NAME="${PROJECT_NAME}_${OS}_${ARCH}.exe"
  else
    OUTPUT_NAME="${PROJECT_NAME}_${OS}_${ARCH}"
  fi

  # Build the project for the platform
  echo "Building $OUTPUT_NAME ..."
  GOOS=$OS GOARCH=$ARCH go build -o $OUTPUT_NAME

  if [ "$OS" == "windows" ]; then
    zip ${OUTPUT_NAME}.zip $OUTPUT_NAME
    mv $OUTPUT_NAME.zip $RELEASE_FILE
  else
    # Compress the output file
    tar -czf ${OUTPUT_NAME}.tar.gz $OUTPUT_NAME
    # Delete the old file 
    mv $OUTPUT_NAME.tar.gz $RELEASE_FILE
  fi
  rm $OUTPUT_NAME

done


echo "All platforms have been built and compressed successfully."
