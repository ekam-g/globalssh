#!/usr/bin/bash

cd || exit

curl -L https://github.com/carghai/globalssh/archive/refs/heads/main.zip > globalssh.zip

unzip  globalssh.zip

cd globalssh-main || exit

make init

make deps

make install

mv globalssh ~

cd || exit

sudo rm -rf globalssh

echo "Welcome to globalssh, rerun to update this script"