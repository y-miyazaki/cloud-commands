#!/bin/bash
#--------------------------------------------------------------
# Github Create Issue
# If it is not the same title and the state is not open, register it as a Github Issue.
#--------------------------------------------------------------

#--------------------------------------------------------------
# variables
#--------------------------------------------------------------
USER_NAME=$1
REPOSITORY_NAME=$2
TITLE=$3
BODY=$4
ASSIGNEES=$5
GITHUB_ACCESS_TOKEN=$6

pageCount=1
while true
do
    infos=$(curl -s https://api.github.com/repos/${USER_NAME}/${REPOSITORY_NAME}/issues?page=${pageCount} | jq -r '.[] | .a = .title + "@" + .state | .a')
    echo "$infos"
    if [ -z "${infos}" ]; then
        break
    fi
    for info in $infos; do
        IFS=$'@'
        datas=(${info})
        title=`echo ${datas[0]} | tr -d "/"`
        state=`echo ${datas[1]} | tr -d "/"`
        echo "$title/$state"
        if [ "${title}" == "${TITLE}" ] && [ "${state}" == "open" ]; then
            exit 0
        fi
    done
    pageCount=`expr $pageCount + 1`
done

curl -s -X POST \
 -H "Accept: application/vnd.github.v3+json" \
 -H "Authorization: token ${GITHUB_ACCESS_TOKEN}" \
 https://api.github.com/repos/${USER_NAME}/${REPOSITORY_NAME}/issues \
 -d "{\"title\":\"${TITLE}\",\"body\":\"${BODY}\",\"assignees\":[\"${ASSIGNEES}\"]}"
