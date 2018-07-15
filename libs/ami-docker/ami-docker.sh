#!/bin/bash

docker stop ami
docker rm ami
if [ $# -eq 0 ]; then
    exit 0
fi

if [ $1 = "start" ]; then
  #docker run --name ami -v $(PWD)/docker/volume:/var/data -d amazonlinux:latest /bin/bash -c 'while true; do echo Hello; sleep 1; done'
  docker run --name ami -v $(PWD)/volume:/var/data amazonlinux:latest /bin/bash -c 'cd /var/data && bash build.sh'
fi
