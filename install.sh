#!/bin/sh
#------------------------------------------------------------------------
# install commands
#------------------------------------------------------------------------
set -e

# variables
INSTALL_DIR=/usr/local/bin
SCRIPT_DIR=$(cd $(dirname $0); pwd)

# create directory and copy commands
# mkdir -p $INSTALL_DIR
# chmod commands
find $SCRIPT_DIR/cmd -maxdepth 1 -type f | xargs chmod 775

# chmod files
chmod 775 $SCRIPT_DIR/cmd/files
find $SCRIPT_DIR/cmd/files -type f | xargs chmod 644
# copy commands
cp -rp ${SCRIPT_DIR}/cmd/* $INSTALL_DIR
