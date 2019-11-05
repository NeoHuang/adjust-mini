#!/bin/bash

rev=$(git rev-parse HEAD)
docker build  --build-arg backend_target=$1 --build-arg version=${rev} ./ -t neohuang/$1

docker tag neohuang/$1 neohuang/$1:${rev}
docker push neohuang/$1:${rev}

kubectl set image deployment/adjust-server-deployment adjust-server=neohuang/$1:${rev}
