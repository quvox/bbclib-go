#!/usr/bin/env bash

cd libs

if [ -z $1 ]; then
    git clone https://github.com/openssl/openssl.git ./openssl
    pushd ./openssl
    git checkout f70425d3ac5e4ef17cfa116d99f8f03bbac1c7f2
    ./config && make
    popd

    pushd libbbcsig
    make clean
    make
    popd
    mv libbbcsig.so ../bbclib/

elif [ $1 = "aws" ]; then
    if [ -z `which docker` ]; then
        echo "docker must be installed"
        exit 1
    fi
    cd ami-docker
    cp -RP ../libbbcsig volume/
    bash ami-docker.sh start
    cp volume/libbbcsig.so ../../bbclib/
    exit
fi


