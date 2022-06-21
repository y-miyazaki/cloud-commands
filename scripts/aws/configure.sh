#/bin/bash
#--------------------------------------------------------------
# Set profile.
# https://docs.aws.amazon.com/ja_jp/cli/latest/userguide/cli-configure-quickstart.html#cli-configure-quickstart-config
#--------------------------------------------------------------
PROFILE=${1:-default}

if [ -n "${AWS_ACCESS_KEY_ID}" ] && [ -n "${AWS_SECRET_ACCESS_KEY}" ] && [ -n "${AWS_DEFAULT_REGION}" ]; then
    aws configure set aws_access_key_id "${AWS_ACCESS_KEY_ID}" --profile "${PROFILE}"
    aws configure set aws_secret_access_key "${AWS_SECRET_ACCESS_KEY}" --profile "${PROFILE}"
    aws configure set region "${AWS_DEFAULT_REGION}" --profile "${PROFILE}"
else
    aws configure -p "${PROFILE}"
fi
aws configure list
