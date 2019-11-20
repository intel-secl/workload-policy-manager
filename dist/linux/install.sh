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
WPM_CA_CONFIGURATION=${WPM_CA_CONFIGURATION:-/etc/${APPLICATION}/cacerts/}
WPM_CA_JWT_DIR=${WPM_CA_JWT_DIR:-/etc/${APPLICATION}/jwt/}
WPM_LOGS=${WPM_LOGS:-/var/log/${APPLICATION}}
INSTALL_LOG_FILE=${INSTALL_LOG_FILE:-$WPM_LOGS/install.log}
WPM_LOG_LEVEL=${WPM_LOG_LEVEL:-INFO}


# Deployment phase
# 2. load installer environment file, if present
if [ -f ~/wpm.env  ]; then
  echo_info "Loading environment variables from $(cd ~ && pwd)/wpm.env"
  . ~/wpm.env
  env_file_exports=$(cat ~/wpm.env | grep -E '^[A-Z0-9_]+\s*=' | cut -d = -f 1)
  if [ -n "$env_file_exports" ]; then eval export $env_file_exports; fi
else
  echo_failure "No environment file"
fi

echo_info "Clearing install logs and writing to it..."
# before we start, clear the install log (directory must already exist; created above)
mkdir -p $(dirname $INSTALL_LOG_FILE)
if [ $? -ne 0 ]; then
  echo_failure "Cannot write to log directory: $(dirname $INSTALL_LOG_FILE)"
  exit 1
fi
date > $INSTALL_LOG_FILE
if [ $? -ne 0 ]; then
  echo_failure "Cannot write to log file: $INSTALL_LOG_FILE"
  exit 1
fi

echo_info "Creating application directories and assigning permissions...."  | tee -a $INSTALL_LOG_FILE
# 8. create application directories (chown will be repeated near end of this script, after setup)
for directory in $WPM_CONFIGURATION $WPM_LOGS $WPM_BIN $WPM_CA_CONFIGURATION $WPM_CA_JWT_DIR; do
  # mkdir -p will return 0 if directory exists or is a symlink to an existing directory or directory and parents can be created
  mkdir -p $directory
  if [ $? -ne 0 ]; then
    echo_failure "Cannot create directory: $directory" | tee -a $INSTALL_LOG_FILE
    exit 1
  fi
  chmod 700 $directory
done

# if an existing wpm is already running, stop it while we install
existing_wpm=`which wpm 2>/dev/null`
if [ -f "$existing_wpm" ]; then
 echo_success "Workload Policy Manager is already installed."  | tee -a $INSTALL_LOG_FILE
 exit 0
fi

cp -f $APPLICATION $WPM_BIN/wpm
ln -sfT $WPM_BIN/wpm $WPM_SYMLINK
echo_success "WPM installation complete"  | tee -a $INSTALL_LOG_FILE

# exit wpm setup if WPM_NOSETUP is set
if [ -n "$WPM_NOSETUP" ]; then
  echo "WPM_NOSETUP value is set. So, skipping the wpm setup task." | tee -a $INSTALL_LOG_FILE
  exit 0;
fi

# 33. wpm setup
wpm setup | tee -a $INSTALL_LOG_FILE

#Install secure docker daemon with wpm only if WPM_WITH_SECURE_DOCKER_DAEMON is enabled in wpm.env
if [ "$WPM_WITH_CONTAINER_SECURITY" == "y" ] || [ "$WPM_WITH_CONTAINER_SECURITY" == "Y" ] || [ "$WPM_WITH_CONTAINER_SECURITY" == "yes" ]; then
  which docker 2>/dev/null
  if [ $? -ne 0 ]; then
    echo "Docker is not installed" | tee -a $INSTALL_LOG_FILE
    exit 1
  fi
  which cryptsetup 2>/dev/null
  if [ $? -ne 0 ]; then
    echo "Installing cryptsetup" | tee -a $INSTALL_LOG_FILE
    yum install -y cryptsetup | tee -a $INSTALL_LOG_FILE
  fi 
  echo "Installing secure docker daemon" | tee -a $INSTALL_LOG_FILE
  systemctl stop docker
  mkdir -p $WPM_HOME/secure-docker-daemon/backup
  cp /usr/bin/docker $WPM_HOME/secure-docker-daemon/backup/
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
  echo "Restarting docker" | tee -a $INSTALL_LOG_FILE
  systemctl daemon-reload
  systemctl start docker
  cp uninstall-secure-docker-daemon.sh $WPM_HOME/secure-docker-daemon/
fi

echo "Installation completed." | tee -a $INSTALL_LOG_FILE
