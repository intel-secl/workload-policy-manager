#!/bin/bash

# WARNING:
# *** do NOT use TABS for indentation, use SPACES (tabs will cause errors in some linux distributions)
# *** do NOT use 'exit' to return from the functions in this file, use 'return' ONLY (exit will cause unit testing hassles)

# TERM_DISPLAY_MODE can be "plain" or "color"
TERM_DISPLAY_MODE=color
TERM_COLOR_GREEN="\\033[1;32m"
TERM_COLOR_CYAN="\\033[1;36m"
TERM_COLOR_RED="\\033[1;31m"
TERM_COLOR_YELLOW="\\033[1;33m"
TERM_COLOR_NORMAL="\\033[0;39m"

# Environment:
# - TERM_DISPLAY_MODE
# - TERM_DISPLAY_GREEN
# - TERM_DISPLAY_NORMAL
echo_success() {
  if [ "$TERM_DISPLAY_MODE" = "color" ]; then echo -en "${TERM_COLOR_GREEN}"; fi
  echo ${@:-"[  OK  ]"}
  if [ "$TERM_DISPLAY_MODE" = "color" ]; then echo -en "${TERM_COLOR_NORMAL}"; fi
  return 0
}

# Environment:
# - TERM_DISPLAY_MODE
# - TERM_DISPLAY_RED
# - TERM_DISPLAY_NORMAL
echo_failure() {
  if [ "$TERM_DISPLAY_MODE" = "color" ]; then echo -en "${TERM_COLOR_RED}"; fi
  echo ${@:-"[FAILED]"}
  if [ "$TERM_DISPLAY_MODE" = "color" ]; then echo -en "${TERM_COLOR_NORMAL}"; fi
  return 1
}

# Environment:
# - TERM_DISPLAY_MODE
# - TERM_DISPLAY_YELLOW
# - TERM_DISPLAY_NORMAL
echo_warning() {
  if [ "$TERM_DISPLAY_MODE" = "color" ]; then echo -en "${TERM_COLOR_YELLOW}"; fi
  echo ${@:-"[WARNING]"}
  if [ "$TERM_DISPLAY_MODE" = "color" ]; then echo -en "${TERM_COLOR_NORMAL}"; fi
  return 1
}

echo_info() {
  if [ "$TERM_DISPLAY_MODE" = "color" ]; then echo -en "${TERM_COLOR_CYAN}"; fi
  echo ${@:-"[INFO]"}
  if [ "$TERM_DISPLAY_MODE" = "color" ]; then echo -en "${TERM_COLOR_NORMAL}"; fi
  return 1
}

############################################################################################################

# default settings
APPLICATION=workload-policy-manager
WPM_HOME=${WPM_HOME:-/opt/${APPLICATION}}
WPM_BIN=${WPM_BIN:-$WPM_HOME/bin}
WPM_SYMLINK=${WPM_SYMLINK:-/usr/local/bin/wpm}
WPM_CONFIGURATION=${WPM_CONFIGURATION:-/etc/${APPLICATION}}
WPM_CA_CONFIGURATION=${WPM_CA_CONFIGURATION:-/etc/${APPLICATION}/certs/trustedca}
WPM_FLAVOR_SIGN_DIR=${WPM_CONFIGURATION}/certs/flavorsign
WPM_KBS_ENVELOPKEY_DIR=${WPM_CONFIGURATION}/certs/kbs
WPM_LOGS=${WPM_LOGS:-/var/log/${APPLICATION}}
WPM_LOG_LEVEL=${WPM_LOG_LEVEL:-INFO}

# Deployment phase
# 2. load installer environment file, if present
if [ -f ~/wpm.env ]; then
  echo_info "Loading environment variables from $(cd ~ && pwd)/wpm.env"
  . ~/wpm.env
  env_file_exports=$(cat ~/wpm.env | grep -E '^[A-Z0-9_]+\s*=' | cut -d = -f 1)
  if [ -n "$env_file_exports" ]; then eval export $env_file_exports; fi
else
  echo_failure "No environment file"
fi

echo_info "Creating application directories and assigning permissions...."
# 8. create application directories (chown will be repeated near end of this script, after setup)
for directory in $WPM_CONFIGURATION $WPM_LOGS $WPM_BIN $WPM_CA_CONFIGURATION $WPM_FLAVOR_SIGN_DIR $WPM_KBS_ENVELOPKEY_DIR; do
  # mkdir -p will return 0 if directory exists or is a symlink to an existing directory or directory and parents can be created
  mkdir -p $directory
  if [ $? -ne 0 ]; then
    echo_failure "Cannot create directory: $directory"
    exit 1
  fi
  chmod 700 $directory
done

# if WPM is already installed then exit
existing_wpm=$(which wpm 2>/dev/null)
if [ -f "$existing_wpm" ]; then
  echo_success "Workload Policy Manager is already installed."
  exit 0
fi

cp -f $APPLICATION $WPM_BIN/wpm
ln -sfT $WPM_BIN/wpm $WPM_SYMLINK

auto_install() {
  local component=${1}
  local cprefix=${2}
  local yum_packages=$(eval "echo \$${cprefix}_YUM_PACKAGES")
  # detect available package management tools. start with the less likely ones to differentiate.
  yum -y install $yum_packages
}

logRotate_clear() {
  logrotate=""
}

logRotate_detect() {
  local logrotaterc=`ls -1 /etc/logrotate.conf 2>/dev/null | tail -n 1`
  logrotate=`which logrotate 2>/dev/null`
  if [ -z "$logrotate" ] && [ -f "/usr/sbin/logrotate" ]; then
    logrotate="/usr/sbin/logrotate"
  fi
}

logRotate_install() {
  LOGROTATE_YUM_PACKAGES="logrotate"
  if [ "$(whoami)" == "root" ]; then
    auto_install "Log Rotate" "LOGROTATE"
    if [ $? -ne 0 ]; then echo_failure "Failed to install logrotate"; exit -1; fi
  fi
  logRotate_clear; logRotate_detect;
    if [ -z "$logrotate" ]; then
      echo_failure "logrotate is not installed"
    else
      echo  "logrotate installed in $logrotate"
    fi
}

logRotate_install

export LOG_ROTATION_PERIOD=${LOG_ROTATION_PERIOD:-weekly}
export LOG_COMPRESS=${LOG_COMPRESS:-compress}
export LOG_DELAYCOMPRESS=${LOG_DELAYCOMPRESS:-delaycompress}
export LOG_COPYTRUNCATE=${LOG_COPYTRUNCATE:-copytruncate}
export LOG_SIZE=${LOG_SIZE:-100M}
export LOG_OLD=${LOG_OLD:-12}

mkdir -p /etc/logrotate.d

if [ ! -a /etc/logrotate.d/wpm ]; then
 echo "/var/log/workload-policy-manager/* {
    missingok
        notifempty
        rotate $LOG_OLD
        maxsize $LOG_SIZE
    nodateext
        $LOG_ROTATION_PERIOD
        $LOG_COMPRESS
        $LOG_DELAYCOMPRESS
        $LOG_COPYTRUNCATE
}" > /etc/logrotate.d/wpm
fi

# exit wpm setup if WPM_NOSETUP is set
if [ "$WPM_NOSETUP" = "true" ]; then
  echo "WPM_NOSETUP value is set to true. So, skipping the wpm setup task."
  exit 0
fi

# 33. wpm setup
wpm setup all
SETUP_RESULT=$?

if [ $SETUP_RESULT -ne 0 ]; then
  echo_failure "Error: WPM setup tasks failed. Exiting."
  exit $SETUP_RESULT
fi

#Install secure docker daemon with wpm only if WPM_WITH_SECURE_DOCKER_DAEMON is enabled in wpm.env
if [ "$WPM_WITH_CONTAINER_SECURITY" = "y" ] || [ "$WPM_WITH_CONTAINER_SECURITY" = "Y" ] || [ "$WPM_WITH_CONTAINER_SECURITY" = "yes" ]; then
  which docker 2>/dev/null
  if [ $? -ne 0 ]; then
    echo_failure "Error: Docker is required for Secure Docker Daemon to be installed!"
    exit 1
  fi
  which cryptsetup 2>/dev/null
  if [ $? -ne 0 ]; then
    echo "Installing cryptsetup"
    yum install -y cryptsetup
    CRYPTSETUP_RESULT=$?
    if [ $CRYPTSETUP_RESULT -ne 0 ]; then
      echo_failure "Error: Secure Docker requires cryptsetup - Install failed. Exiting."
      exit $CRYPTSETUP_RESULT
    fi
  fi
  echo "Installing secure docker daemon"
  systemctl stop docker
  mkdir -p $WPM_HOME/secure-docker-daemon/backup
  cp /usr/bin/docker $WPM_HOME/secure-docker-daemon/backup/
  cp /etc/docker/daemon.json $WPM_HOME/secure-docker-daemon/backup/ 2>/dev/null
  chown -R root:root docker-daemon
  cp -f docker-daemon/docker /usr/bin/
  which /usr/bin/dockerd-ce 2>/dev/null
  if [ $? -ne 0 ]; then
    cp /usr/bin/dockerd $WPM_HOME/secure-docker-daemon/backup/
    cp -f docker-daemon/dockerd-ce /usr/bin/dockerd
  else
    cp /usr/bin/dockerd-ce $WPM_HOME/secure-docker-daemon/backup/
    cp -f docker-daemon/dockerd-ce /usr/bin/dockerd-ce
  fi
  mkdir -p /etc/docker
  cp daemon.json /etc/docker/
  echo "Restarting docker"
  systemctl daemon-reload
  systemctl start docker
  cp uninstall-secure-docker-daemon.sh $WPM_HOME/secure-docker-daemon/
fi

echo "Installation completed."
