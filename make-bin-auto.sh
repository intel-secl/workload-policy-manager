#!/bin/bash

# Postconditions:
# * exit with error code 1 only if there was a fatal error:
#   functions.sh not found (must be adjacent to this file in the package)
#   


#####


# WARNING:
# *** do NOT use TABS for indentation, use SPACES
# *** TABS will cause errors in some linux distributions

# application defaults (these are not configurable and used only in this script so no need to export)
workspace=./target/wpm-1.0
version=1.0
goExecutableName=wpm
export TMPDIR=~/.tmp

# before we start, clear the install log (directory must already exist; created above)
mkdir -p $workspace
if [ $? -ne 0 ]; then
  echo_failure "Cannot write to log directory: $(dirname $INSTALL_LOG_FILE)"
  exit 1
fi

echo "Created target/wpm1.0 directory"

mkdir -p $TMPDIR
if [ $? -ne 0 ]; then
  echo_failure "Cannot create tmp directory"
  exit 1
fi
cp install/*.sh $workspace
# go build -o $workspace/$goExecutableName ./cmd/wpm/main.go
# cp $workspace/$goExecutableName /usr/local/bin/


# installer name
projectNameVersion=`basename "${workspace}"`
# where to save the installer (parent of directory containing files)
targetDir=`dirname "${workspace}"`

if [ -z "$workspace" ]; then
  echo_info "Usage: $0 <workspace>"
  echo_info "Example: $0 /path/to/wpm-1.0.bin"
  echo_info "The self-extracting installer wpm-1.0.bin would be created in /path/to"
  exit 1
fi

if [ ! -d "$workspace" ]; then echo "Cannot find workspace '$workspace'"; exit 1; fi

# ensure all executable files in the target folder have the x bit set
chmod +x $workspace/*.sh

# check for the makeself tool
makeself=`which makeself`
if [ -z "$makeself" ]; then
    echo_failure "Missing makeself tool"
    exit 1
fi

$makeself --follow --nocomp "$workspace" "$targetDir/${projectNameVersion}.bin" "$projectNameVersion" ./setup.sh
rm -rf $TMPDIR
