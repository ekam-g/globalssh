#!/usr/bin/bash

cd || exit

sudo rm -rf globalssh-main

rm globalssh.zip

curl -L https://github.com/carghai/globalssh/archive/refs/heads/main.zip > globalssh.zip || exit

unzip  globalssh.zip || exit

cd globalssh-main || exit

make init

make deps

make install

mv globalssh ~

cd || exit

sudo rm -rf globalssh-main

rm globalssh.zip

echo "Welcome to globalssh, rerun to update this script. If you are on mac it may not add to path"