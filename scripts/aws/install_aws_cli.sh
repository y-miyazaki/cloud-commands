#!/bin/bash

#--------------------------------------------------------------
# Check OS
#--------------------------------------------------------------
if [ "$(uname)" == 'Darwin' ]; then
    OS='Mac'
elif [ "$(expr substr $(uname -s) 1 5)" == 'Linux' ]; then
    OS='Linux'
elif [ "$(expr substr $(uname -s) 1 10)" == 'MINGW32_NT' ]; then                                                                                           
    OS='Cygwin'
    echo "Your platform ($(uname -a)) is not supported."
    exit 1
else
    echo "Your platform ($(uname -a)) is not supported."
    exit 1
fi

if [ "${OS}" == "Linux" ]; then
    # install aws cli
    curl "https://awscli.amazonaws.com/awscli-exe-linux-x86_64.zip" -o "awscliv2.zip"
    unzip awscliv2.zip
    sudo ./aws/install
    rm -f awscliv2.zip
    rm -rf ./aws
    # https://docs.aws.amazon.com/ja_jp/systems-manager/latest/userguide/session-manager-working-with-install-plugin.html
    echo "install aws ssm plugin..."
    if [ -e "$(command -v yum)" ]; then
        curl "https://s3.amazonaws.com/session-manager-downloads/plugin/latest/linux_64bit/session-manager-plugin.rpm" -o "session-manager-plugin.rpm"
        sudo yum install -y session-manager-plugin.rpm
        session-manager-plugin
    elif [ -e "$(command -v dpkg)" ]; then
        curl "https://s3.amazonaws.com/session-manager-downloads/plugin/latest/ubuntu_64bit/session-manager-plugin.deb" -o "session-manager-plugin.deb"
        sudo dpkg -i session-manager-plugin.deb
        session-manager-plugin
    fi
elif [ "${OS}" == "Mac" ]; then
    # install aws cli
    echo "install aws cli..."
    curl "https://awscli.amazonaws.com/AWSCLIV2.pkg" -o "AWSCLIV2.pkg"
    sudo installer -pkg AWSCLIV2.pkg -target /
    aws --version
    rm -f AWSCLIV2.pkg
    # install aws ssm plugin
    # https://docs.aws.amazon.com/ja_jp/systems-manager/latest/userguide/session-manager-working-with-install-plugin.html
    echo "install aws ssm plugin..."
    curl "https://s3.amazonaws.com/session-manager-downloads/plugin/latest/mac/sessionmanager-bundle.zip" -o "sessionmanager-bundle.zip"
    unzip sessionmanager-bundle.zip
    sudo ./sessionmanager-bundle/install -i /usr/local/sessionmanagerplugin -b /usr/local/bin/session-manager-plugin
    session-manager-plugin
    rm -f sessionmanager-bundle.zip
    rm -rf sessionmanager-bundle
fi
