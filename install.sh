#!/usr/bin/bash

cd || exit

git clone https://github.com/carghai/globalssh

cd globalssh || exit

make init

make deps

make install

cd || exit

sudo rm -rf globalssh

echo "Welcome to globalssh, rerun to update this script"