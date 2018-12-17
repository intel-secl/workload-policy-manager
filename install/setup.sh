#!/bin/bash

# Postconditions:
# * exit with error code 1 only if there was a fatal error:
#   functions.sh not found (must be adjacent to this file in the package)
#   


#####


# WARNING:
# *** do NOT use TABS for indentation, use SPACES
# *** TABS will cause errors in some linux distributions

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




# application defaults (these are not configurable and used only in this script so no need to export)
DEFAULT_WPM_HOME=/opt/wpm
DEFAULT_WPM_USERNAME=wpm

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
export WPM_BIN=${WPM_BIN:-$WPM_HOME/bin}
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

# echo "Switching to wpm user"
# # 7. determine if we are installing as root or non-root, create groups and users accordingly
# if [ "$(whoami)" == "root" ]; then
#   # create a wpm user if there isn't already one created
#   WPM_USERNAME=${WPM_USERNAME:-$DEFAULT_WPM_USERNAME}
#   if ! getent passwd $WPM_USERNAME 2>&1 >/dev/null; then
#     useradd --comment "Workload Policy Manager" --home $WPM_HOME --system --shell /bin/false $WPM_USERNAME
#     usermod --lock $WPM_USERNAME
#     # note: to assign a shell and allow login you can run "usermod --shell /bin/bash --unlock $WPM_USERNAME"
#   fi
# else
#   # already running as wpm user
#   WPM_USERNAME=$(whoami)
#   if [ ! -w "$WPM_HOME" ] && [ ! -w $(dirname $WPM_HOME) ]; then
#     WPM_HOME=$(cd ~ && pwd)
#   fi
#   echo_warning "Installing as $WPM_USERNAME into $WPM_HOME"  
# fi

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
# chown $WPM_USERNAME:$WPM_USERNAME $INSTALL_LOG_FILE
logfile=$INSTALL_LOG_FILE

echo "Create application directories and assign permissions"
# 8. create application directories (chown will be repeated near end of this script, after setup)
for directory in $WPM_HOME $WPM_CONFIGURATION $WPM_ENV $WPM_LOGS $WPM_BIN; do
  # mkdir -p will return 0 if directory exists or is a symlink to an existing directory or directory and parents can be created
  mkdir -p $directory
  if [ $? -ne 0 ]; then
    echo "Cannot create directory: $directory"
    exit 1
  fi
  # chown -R $WPM_USERNAME:$WPM_USERNAME $directory
  chmod 700 $directory
done

echo "Writing env file into config file"
# write the env variables to a config file
if [ -d "$WPM_CONFIGURATION" ]
then
    touch $WPM_CONFIGURATION_FILE
    echo "$WPM_ENV_FILE" > "$WPM_CONFIGURATION_FILE"
fi

# chown $WPM_USERNAME:$WPM_USERNAME $WPM_CONFIGURATION_FILE

echo "Adding wpm bin to path variable"
# ensure we have our own wpm programs in the path
export PATH=$WPM_BIN:$PATH

# ensure that trousers and tpm tools are in the path
export PATH=$PATH:/usr/sbin:/usr/local/sbin

profile_dir=$HOME
if [ "$(whoami)" == "root" ] && [ -n "$WPM_USERNAME" ] && [ "$WPM_USERNAME" != "root" ]; then
  profile_dir=$WPM_HOME
fi

# if an existing wpm is already running, stop it while we install
# existing_wpm=`which wpm 2>/dev/null`
# if [ -f "$existing_wpm" ]; then
  # echo "Workload Policy Manager is already installed."
  # exit 0
# fi

echo "store directory layout in env file"
# 11. store directory layout in env file
echo "# $(date)" > $WPM_ENV/wpm-layout
echo "WPM_HOME=$WPM_HOME" >> $WPM_ENV/wpm-layout
echo "WPM_CONFIGURATION=$WPM_CONFIGURATION" >> $WPM_ENV/wpm-layout
echo "WPM_BIN=$WPM_BIN" >> $WPM_ENV/wpm-layout
echo "WPM_LOGS=$WPM_LOGS" >> $WPM_ENV/wpm-layout

# 12. store wpm username in env file
echo "# $(date)" > $WPM_ENV/wpm-username
echo "WPM_USERNAME=$WPM_USERNAME" >> $WPM_ENV/wpm-username

# store the auto-exported environment variables in temporary env file
# to make them available after the script uses sudo to switch users;
# we delete that file later
echo "# $(date)" > $WPM_ENV/wpm-setup
for env_file_var_name in $env_file_exports
do
  eval env_file_var_value="\$$env_file_var_name"
  echo "export $env_file_var_name='$env_file_var_value'" >> $WPM_ENV/wpm-setup
done

# # add bin and sbin directories in wpm home directory to path
# bin_directories=$(find_subdirectories ${WPM_HOME} bin; find_subdirectories ${WPM_HOME} sbin)
# bin_directories_path=$(join_by : ${bin_directories[@]})
# for directory in ${bin_directories[@]}; do
#   chmod -R 700 $directory
# done
# export PATH=$bin_directories_path:$PATH
# appendToUserProfileFile "export PATH=${bin_directories_path}:\$PATH" $profile_name

# copy the go executable to WPM_BIN directory and /usr/local/bin
# bin_location=/usr/local/bin/
# wpm_bin=wpm
# wpm_bin_dir=/opt/wpm/bin/
# echo "Harshitha1"
# echo $wpm_bin
#cp wpm /opt/wpm/bin/
# cp $wpm_bin $bin_location
# echo "Harshitha2"
# cp $wpm_bin $wpm_bin_dir
# chown -R $WPM_USERNAME:$WPM_USERNAME $WPM_HOME
#chmod 755 $WPM_BIN/*


# 16. symlink wpm
# if prior version had control script in /usr/local/bin, delete it
# if [ "$(whoami)" == "root" ] && [ -f /usr/local/bin/wpm ]; then
  # rm /usr/local/bin/wpm
# fi
# EXISTING_WPM_COMMAND=`which wpm 2>/dev/null`
# if [ -n "$EXISTING_WPM_COMMAND" ]; then
  # rm -f "$EXISTING_WPM_COMMAND"
# fi
# link /usr/local/bin/wpm -> /opt/wpm/bin/wpm
# ln -s $WPM_BIN/wpm.sh /usr/local/bin/wpm
# if [[ ! -h $WPM_BIN/wpm ]]; then
#   ln -s $WPM_BIN/wpm.sh $WPM_BIN/wpm
# fi

# Ensure we have given wpm access to its files
# for directory in $WPM_HOME $WPM_CONFIGURATION $WPM_ENV $WPM_LOGS; do
#   echo "chown -R $WPM_USERNAME:$WPM_USERNAME $directory" >>$logfile
#   chown -R $WPM_USERNAME:$WPM_USERNAME $directory 2>>$logfile
# done

# Make the logs dir owned by wpm user
# chown -R $WPM_USERNAME:$WPM_USERNAME $WPM_LOGS/

# 29. ensure the wpm owns all the content created during setup
# for directory in $WPM_HOME $WPM_CONFIGURATION $WPM_BIN $WPM_ENV $WPM_LOGS; do
#   chown -R $WPM_USERNAME:$WPM_USERNAME $directory
# done

# exit wpm setup if WPM_NOSETUP is set
if [ -n "$WPM_NOSETUP" ]; then
  echo "WPM_NOSETUP value is set. So, skipping the wpm setup task."
  exit 0;
fi

echo "WPM installation complete"
# 33. wpm setup
WPM_SETUP_TASKS="create-envelope-key register-envelope-key-with-kbs"
wpm setup --all 

