#!/bin/bash

# Define the release file name and the project name
RELEASE_FILE="global_ssh.tar.gz"
PROJECT_NAME="global_ssh"

# Extract the release file
tar -xzf $RELEASE_FILE

# Get the list of supported platforms
PLATFORMS=("linux/amd64" "linux/arm64" "linux/arm" "darwin/amd64" "darwin/arm64" "windows/amd64" "windows/arm64")

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
    mv $OUTPUT_NAME releases
  else
    # Compress the output file
    tar -czf ${OUTPUT_NAME}.tar.gz $OUTPUT_NAME
    # Delete the old file 
    rm $OUTPUT_NAME
    mv $OUTPUT_NAME.tar.gz releases
  fi

done


echo "All platforms have been built and compressed successfully."
