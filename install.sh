#!/bin/bash

cd || exit

sudo rm -rf globalssh-main

rm globalssh.zip

sudo mkdir /usr/local/bin

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

echo "Welcome to Global SSH, If you are on windows it may not add to path"