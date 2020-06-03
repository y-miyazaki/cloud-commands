#!/bin/sh
#------------------------------------------------------------------------
# install commands
#------------------------------------------------------------------------
set -e

# variables
CLOUD_COMMANDS_VERSION=$1
INSTALL_DIR=/usr/local/bin
SCRIPT_DIR=$(cd $(dirname $0); pwd)
echo ${CLOUD_COMMANDS_VERSION}

# download cmd.zip and unzip
curl -L -o cmd.zip https://github.com/y-miyazaki/cloud-commands/releases/download/${CLOUD_COMMANDS_VERSION}/cmd.zip
mkdir -p $SCRIPT_DIR/cloud-commands
ls -la
unzip -d $SCRIPT_DIR/cloud-commands cmd.zip

# create directory and copy commands
# chmod commands
find $SCRIPT_DIR/cloud-commands/cmd -maxdepth 1 -type f | xargs chmod 775

# chmod files
chmod 775 $SCRIPT_DIR/cloud-commands/cmd/files
find $SCRIPT_DIR/cloud-commands/cmd/files -type f | xargs chmod 644
# copy commands
sudo cp -rp $SCRIPT_DIR/cloud-commands/cmd/* $INSTALL_DIR

# remove files
rm -f cmd.zip
rm -rf $SCRIPT_DIR/cloud-commands
