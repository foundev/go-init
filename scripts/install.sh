#!/usr/bin/env bash

orig_dir=$(pwd)
cd /tmp
curl -O https://github.com/foundev/go-init/archive/refs/heads/main.zip
unzip main.zip
rm main.zip
cd main
./scripts/build
echo "copying binary to /usr/local/bin/go-init need sudo permissions to write"
sudo cp ./bin/go-init /usr/local/bin/
cd ..
rm -fr main
cd $orig_dir
