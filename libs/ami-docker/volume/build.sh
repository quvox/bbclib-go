#!/usr/bin/env bash

yum groupinstall -y "Development Tools"

git clone https://github.com/openssl/openssl.git ./openssl
pushd ./openssl
git checkout f70425d3ac5e4ef17cfa116d99f8f03bbac1c7f2
./config && make
popd

pushd libbbcsig
make clean
make
popd
