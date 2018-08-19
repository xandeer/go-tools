#!/usr/bin/env sh
workspace=$1
name=$2
url=$3

cd $workspace
if [ -d $name ]; then
  cd $name
else
  git clone $url
  cd $name
fi
make build