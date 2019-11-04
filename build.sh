#!/bin/bash

rev=$(git rev-parse HEAD)
docker build  --build-arg backend_target=$1 ./ -t neohuang/$1:${rev}
