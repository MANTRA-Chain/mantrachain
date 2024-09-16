#!/usr/bin/env bash

go_arch=$1
go_os=$2
project_name=$3

# copy build directory to the corresponding directory
cp -r dist-merged/${project_name}-${go_os}-${go_arch}* dist/
