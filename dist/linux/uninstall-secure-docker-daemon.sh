#!/bin/bash

# remove all running containers
docker rm -f $(docker container ls –aq)

# unmount and remove the secureoverlay2 layer data
for m in $(mount -t overlay | grep /var/lib/docker/secureoverlay2/ | awk '{print $3}')
do
  umount -d -f -R $m
done
if [ -d "/var/lib/docker/secureoverlay2" ];
then
  umount -d -f -R /var/lib/docker/secureoverlay2
  rm -rf /var/lib/docker/secureoverlay2
fi

# purge the stale layer data
echo y | docker system prune -a 2>/dev/null

#Copy all the vanilla docker daemon binaries from backup to /usr/bin/ and reconfigure the docker.service file to support vanilla docker
systemctl stop docker.service
cp -f /opt/workload-policy-manager/secure-docker-daemon/backup/dockerd /usr/bin/
cp -f /opt/workload-policy-manager/secure-docker-daemon/backup/docker /usr/bin/

# restore original daemon.json else remove current version
if [ -f /opt/workload-policy-manager/secure-docker-daemon/backup/daemon.json ]; then
  cp -f /opt/workload-policy-manager/secure-docker-daemon/backup/daemon.json /etc/docker/daemon.json
else
  rm -f /etc/docker/daemon.json
fi

systemctl daemon-reload
