#!/bin/bash

echo "building binary"

go build

echo "removing old binary"

sudo rm /usr/local/bin/global_ssh

echo "moving binary to path"

sudo mv global_ssh /usr/local/bin
