#!/bin/bash

#Copy all the vanilla docker daemon binaries from backup to /usr/bin/ and reconfigure the docker.service file to support vanilla docker
systemctl stop docker.service
cp -f /opt/workload-policy-manager/secure-docker-daemon/backup/dockerd /usr/bin/
cp -f /opt/workload-policy-manager/secure-docker-daemon/backup/docker /usr/bin/
cp -f /opt/workload-policy-manager/secure-docker-daemon/backup/daemon.json /etc/docker/ 2>/dev/null
sed -i 's/^ExecStart=.*/ExecStart=\/usr\/bin\/dockerd\ \-H\ unix\:\/\/ /' /lib/systemd/system/docker.service
systemctl daemon-reload
systemctl start docker.service
