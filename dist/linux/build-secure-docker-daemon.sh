#!/bin/bash

git clone ssh://git@gitlab.devtools.intel.com:29418/sst/isecl/secure_docker_daemon.git 2>/dev/null 

cd secure_docker_daemon
git fetch
git checkout v1.0/develop
git pull

#Build secure docker daemon

make > /dev/null

if [ $? -ne 0 ]; then
  echo "could not build secure docker daemon"
  exit 1
fi
 
echo "Successfully build secure docker daemon" 
