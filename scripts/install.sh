#!/usr/bin/env bash

orig_dir=$(pwd)
cd /tmp
curl -L https://github.com/foundev/go-init/archive/refs/heads/main.zip -o main.zip
unzip main.zip
rm main.zip
cd go-init-main
./scripts/build
echo "copying binary to /usr/local/bin/go-init need sudo permissions to write"
sudo cp ./bin/go-init /usr/local/bin/
cd ..
rm -fr go-init-main
cd $orig_dir
