# syntax = docker/dockerfile:experimental
#--------------------------------------------------------------
# build image
#--------------------------------------------------------------
ARG IMAGE_VERSION=latest
FROM ymiyazakixyz/golang-build:${IMAGE_VERSION} AS build

#--------------------------------------------------------------
# ARG
#--------------------------------------------------------------
ARG WORKDIR=/workspace
ARG GIT_DOMAIN=github.com
ARG GIT_CLONE_URL
ARG GIT_BRANCH=master

#--------------------------------------------------------------
# workdir
#--------------------------------------------------------------
RUN mkdir -p ${WORKDIR}
WORKDIR  ${WORKDIR}

# clone and go mod download
RUN --mount=type=secret,id=ssh,dst=/root/.ssh/id_rsa \
    echo "Host ${GIT_DOMAIN}\n\tStrictHostKeyChecking no\n\tIdentityFile /root/.ssh/id_rsa\n" >> /root/.ssh/config && \
    ssh-keyscan -H ${GIT_DOMAIN} >> /root/.ssh/known_hosts && \
    git config --global url."git@${GIT_DOMAIN}:".insteadOf "https://${GIT_DOMAIN}/" && \
    # dependencies
    go mod download

CMD []
