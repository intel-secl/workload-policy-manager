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

# before we start, clear the install log (directory must already exist; created above)
mkdir -p $workspace
if [ $? -ne 0 ]; then
  echo_failure "Cannot write to log directory: $(dirname $INSTALL_LOG_FILE)"
  exit 1
fi

echo $
echo "Created target/wpm1.0 directory"

mkdir -p ~/.tmp/
cp install/*.sh $workspace
go build -o $workspace/$goExecutableName ./cmd/wpm/main.go
cp $workspace/$goExecutableName /usr/local/bin/


# installer name
projectNameVersion=`basename "${workspace}"`
# where to save the installer (parent of directory containing files)
targetDir=`dirname "${workspace}"`

if [ -z "$workspace" ]; then
  echo "Usage: $0 <workspace>"
  echo "Example: $0 /path/to/AttestationService-0.5.1"
  echo "The self-extracting installer AttestationService-0.5.1.bin would be created in /path/to"
  exit 1
fi

if [ ! -d "$workspace" ]; then echo "Cannot find workspace '$workspace'"; exit 1; fi

# ensure all executable files in the target folder have the x bit set
chmod +x $workspace/*.sh

# check for the makeself tool
makeself=`which makeself`
if [ -z "$makeself" ]; then
    echo "Missing makeself tool"
    exit 1
fi

export TMPDIR=~/.tmp
$makeself --follow --nocomp "$workspace" "$targetDir/${projectNameVersion}.bin" "$projectNameVersion" ./setup.sh
rm -rf $TMPDIR
