#!/bin/bash

#Copy all the vanilla docker daemon binaries from backup to /usr/bin/ and reconfigure the docker.service file to support vanilla docker
systemctl stop docker.service
cp -f /opt/wpm/secure-docker-daemon/backup/* /usr/bin/
sed -i 's/^ExecStart=.*/ExecStart=\/usr\/bin\/dockerd\ \-H\ unix\:\/\/ /' /lib/systemd/system/docker.service
systemctl daemon-reload
systemctl start docker.service
