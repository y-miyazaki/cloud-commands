#!/bin/bash

# var
BASE_URL=https://github.com/
BRANCHES=("develop")
USER_NAME=$1
REPOSITORY_NAME=$2
TOKEN=$3

UPDATED_AT=$(date)

# Get Repositoies from orgnization
full_name=`curl -s \
  -H "Accept: application/vnd.github.v3+json" \
  -H "Authorization: token ${TOKEN}" \
  https://api.github.com/${USER_NAME}/${REPOSITORY_NAME}/repos | jq -r '.[] | .a = .full_name + "@" + .description | .a'`
ORG_IFS=$IFS

# Header
echo -e -n "# List of repository status\n"
echo -e -n "\n"
echo -e -n "Updated: ${UPDATED_AT}\n"
echo -e -n "\n"

# Index
echo -e -n "## Index\n"
echo -e -n "\n"
IFS=$'\n'
for fn in ${full_name}; do
    IFS=$'@'
    datas=(${fn})
    anchor=`echo ${datas[0]} | tr -d "/"`
    if [[ "${datas[1]}" == *Archived* || "${datas[1]}" == *Deprecated* ]]; then
        echo -e -n "- [~~${datas[0]}~~](#${anchor}) (Archived or Deprecated)\n"
    else
        echo -e -n "- [${datas[0]}](#${anchor})\n"
    fi
done
echo -e -n "\n\n\n"

# Repositories
echo -e -n "## Repositories\n"
echo -e -n "\n"
IFS=$'\n'
for fn in ${full_name}; do
    IFS=$'@'
    datas=(${fn})
    comment_out=0
    repository=${datas[0]}
    description=${datas[1]}
    if [[ "${datas[1]}" == *Archived* || "${datas[1]}" == *Deprecated* ]]; then
        description="<font color=gray>${description}</font>"
        comment_out=1
    fi

    # for each repository
    if [ ${comment_out} == 0 ]; then
        echo -e -n "## ${repository}\n"
    else
        echo -e -n "## ~~${repository}~~\n"
    fi
    echo -e -n "${BASE_URL}${datas[0]}  \n"
    echo -e -n "${description}  \n"

    if [ ${comment_out} == 0 ]; then
        branch_base_url="${BASE_URL}${datas[0]}/tree/"

        echo -e -n "- Specific latest branch status \n"
        # list of specify branch
        echo -e -n "  | Branch | Latest commit name | Latest commit date | verified | protected |\n"
        echo -e -n "  | :----- | :----------------- | :----------------- | :------: | :-------- |\n"
        for b in "${BRANCHES[@]}"; do
            branch=`curl -s \
                -H "Accept: application/vnd.github.v3+json" \
                -H "Authorization: token ${TOKEN}" \
                https://api.github.com/repos/${datas[0]}/branches/${b} | jq --arg b ${b} --arg url ${branch_base_url} -r 'if .message then .data="  | " + $b + " | - | - | - | "  else .data= "  | [" + .name + "](" + $url + .name + ") | " + "<img src=\"" + .commit.author.avatar_url + "\" width=\"15\"> [" + .commit.author.login + "](" + .commit.author.html_url + ") | " + .commit.commit.author.date + " | " + if .commit.commit.verification.verified then ":white_check_mark:" else ":x:" end + " | " + (.protected|tostring) + " | " end | .data'`
            echo $branch
        done
        echo -e -n "\n"

        # All branches list
        branch_all=`curl -s \
            -H "Accept: application/vnd.github.v3+json" \
            -H "Authorization: token ${TOKEN}" \
            https://api.github.com/repos/${datas[0]}/branches?per_page=100 | jq --arg url ${branch_base_url} -r '.[] | \
                .data="  | [" + .name + "](" + $url + .name + ") | " + .commit.sha + " | " | .data'`
        echo -e -n "  <details><summary>All branches list(Max 100)</summary>  \n\n"
        echo -e -n "  | Branch | sha |\n"
        echo -e -n "  | :----- | :-- |\n"
        echo $branch_all
        echo -e -n "  </details>\n"
        echo -e -n "\n"

        # for each pulls
        prs=`curl -s \
            -H "Accept: application/vnd.github.v3+json" \
            -H "Authorization: token ${TOKEN}" \
            https://api.github.com/repos/${datas[0]}/pulls?state=open | jq --arg url ${branch_base_url} -r '.[] | \
                if .user.avatar_url then \
                    .a = "<img src=\"" + .user.avatar_url + "\" width=\"15\"> [" + .user.login + "](" + .user.html_url + ")" \
                else \
                    .a = "[" + .user.login + "](" + .user.html_url + ")" \
                end | \
                "  [" + .title + "](" + .html_url + ") | " + .a + " | [" + .base.ref + "](" + $url + .base.ref + ") from [" + .head.ref + "](" + $url + .head.ref + ") | " + .created_at + " | " + .updated_at + " | "'`
                # if .requested_reviewers then \
                #     .b = .requested_reviewers.login \
                # else \
                #     .b = "" \
                # end | \
                # "[" + .title + "](" + .html_url + ") | " + .a + " | " + .b + " | " + .created_at + " | " + .updated_at + " | "'`
        echo -e -n "- All pull requests(status=open)(Max 100)  \n"
        echo -e -n "  | Pull requests | Created user | branch | Created | Updated | \n"
        echo -e -n "  | :------------ | :----------- | :----- | :------ | :------ | \n"
        IFS=$'\n'
        for pr in ${prs}; do
            echo $pr
        done
        echo -e -n "\n"
    fi
    echo -e -n "</br></br></br>\n\n"
done
IFS=$ORG_IFS
