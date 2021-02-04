#!/bin/bash

git clone https://gitlab.devtools.intel.com/sst/isecl/secure-docker-daemon.git 2>/dev/null

cd secure-docker-daemon
git fetch
git checkout v3.5/develop
git pull

#Build secure docker daemon

make > /dev/null

if [ $? -ne 0 ]; then
  echo "could not build secure docker daemon"
  exit 1
fi
 
echo "Successfully build secure docker daemon" 
