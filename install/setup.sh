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


# application defaults (these are not configurable and used only in this script so no need to export)
DEFAULT_WPM_HOME=/opt/wpm

# default settings
export WPM_ADMIN_USERNAME=${WPM_ADMIN_USERNAME:-wpm-admin}
export WPM_HOME=${WPM_HOME:-$DEFAULT_WPM_HOME}
WPM_LAYOUT=${WPM_LAYOUT:-home}

# the env directory is not configurable; it is defined as WPM_HOME/env.d and the
# administrator may use a symlink if necessary to place it anywhere else
export WPM_ENV=$WPM_HOME/env.d


# 1. load application environment variables if already defined from env directory
if [ -d $WPM_ENV ]; then
  WPM_ENV_FILES=$(ls -1 $WPM_ENV/*)
  for env_file in $WPM_ENV_FILES; do
    . $env_file
    env_file_exports=$(cat $env_file | grep -E '^[A-Z0-9_]+\s*=' | cut -d = -f 1)
    echo_info $env_file_exports
    if [ -n "$env_file_exports" ]; then eval export $env_file_exports; fi
  done
fi

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

echo_info "Creating directory layout"

# LOCAL CONFIGURATION
directory_layout() {
export WPM_CONFIGURATION=${WPM_CONFIGURATION:-$WPM_HOME/configuration}
export WPM_LOGS=${WPM_LOGS:-$WPM_HOME/logs}
export INSTALL_LOG_FILE=$WPM_LOGS/install.log
export WPM_CONFIGURATION_FILE=$WPM_CONFIGURATION/wpm.properties
}

# 5. define application directory layout
directory_layout

# Output:
# - variable "yum" contains path to yum or empty
yum_detect() {
  yum=`which yum 2>/dev/null`
  if [ -n "$yum" ]; then return 0; else return 1; fi
}

# check for the makeself tool
makeself=`which makeself`
if [ -z "$makeself" ]; then
   if yum_detect; then
      echo_info "Installing yum"
      yum install makeself
    else
      echo_failure "Yum package not detected"
      exit 1
    fi
fi

echo "Clear the install logs and write to it"
# before we start, clear the install log (directory must already exist; created above)
mkdir -p $(dirname $INSTALL_LOG_FILE)
if [ $? -ne 0 ]; then
  echo "Cannot write to log directory: $(dirname $INSTALL_LOG_FILE)"
  exit 1
fi
date > $INSTALL_LOG_FILE
if [ $? -ne 0 ]; then
  echo "Cannot write to log file: $INSTALL_LOG_FILE"
  exit 1
fi

echo "Create application directories and assign permissions"
# 8. create application directories (chown will be repeated near end of this script, after setup)
for directory in $WPM_HOME $WPM_CONFIGURATION $WPM_ENV $WPM_LOGS; do
  # mkdir -p will return 0 if directory exists or is a symlink to an existing directory or directory and parents can be created
  mkdir -p $directory
  if [ $? -ne 0 ]; then
    echo "Cannot create directory: $directory"
    exit 1
  fi
  chmod 700 $directory
done

WPM_ENV_FILE=~/wpm.env
echo_info "Writing env file into config file"
# write the env variables to a config file
if [ -d "$WPM_CONFIGURATION" ]; then
    touch $WPM_CONFIGURATION_FILE
    cat $WPM_ENV_FILE > $WPM_CONFIGURATION_FILE
fi

# if an existing wpm is already running, stop it while we install
existing_wpm=`which wpm 2>/dev/null`
if [ -f "$existing_wpm" ]; then
 echo_success "Workload Policy Manager is already installed."
 exit 0
fi

echo "store directory layout in env file"
# 11. store directory layout in env file
echo "# $(date)" > $WPM_ENV/wpm-layout
echo "WPM_HOME=$WPM_HOME" >> $WPM_ENV/wpm-layout
echo "WPM_CONFIGURATION=$WPM_CONFIGURATION" >> $WPM_ENV/wpm-layout
echo "WPM_LOGS=$WPM_LOGS" >> $WPM_ENV/wpm-layout



# store the auto-exported environment variables in temporary env file
# to make them available after the script uses sudo to switch users;
# we delete that file later
echo "# $(date)" > $WPM_ENV/wpm-setup
for env_file_var_name in $env_file_exports
do
  eval env_file_var_value="\$$env_file_var_name"
  echo "export $env_file_var_name='$env_file_var_value'" >> $WPM_ENV/wpm-setup
done

# exit wpm setup if WPM_NOSETUP is set
if [ -n "$WPM_NOSETUP" ]; then
  echo "WPM_NOSETUP value is set. So, skipping the wpm setup task."
  exit 0;
fi

echo "WPM installation complete"
WPM_BIN_NAME=wpm
cp $WPM_BIN_NAME /usr/local/bin/
# 33. wpm setup
WPM_SETUP_TASKS="create-envelope-key register-envelope-key-with-kbs"
wpm setup --all 

