#!/bin/bash

git clone -b master https://github.com/beyond-blockchain/libbbcsig.git libs
cd libs

if [ -z $1 ]; then
    bash prepare.sh
elif [ $1 = "aws" ]; then
    bash prepare.sh aws
fi
