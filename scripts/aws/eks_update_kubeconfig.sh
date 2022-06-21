#/bin/bash
#--------------------------------------------------------------
# Set kubeconfig
#--------------------------------------------------------------
KUBECTL_VERSION=v1.22.0
REGION=$1
CLUSTER=$2
PROFILE=${3:-default}

#--------------------------------------------------------------
# Set kubeconfig
#--------------------------------------------------------------
if [ -z "$(command -v kubectl)" ]; then
    curl -LO https://storage.googleapis.com/kubernetes-release/release/${KUBECTL_VERSION}/bin/linux/amd64/kubectl
    chmod +x ./kubectl
    mv ./kubectl /usr/local/bin/kubectl
    echo 'source <(kubectl completion bash)' >>~/.bashrc
fi
if [ "${PROFILE}" == "default"]; then
    if [ -n "${AWS_ACCESS_KEY_ID}" ] && [ -n "${AWS_SECRET_ACCESS_KEY}" ] && [ -n "${AWS_DEFAULT_REGION}" ]; then
        aws configure set aws_access_key_id "${AWS_ACCESS_KEY_ID}" --profile "${PROFILE}"
        aws configure set aws_secret_access_key "${AWS_SECRET_ACCESS_KEY}" --profile "${PROFILE}"
        aws configure set region "${AWS_DEFAULT_REGION}" --profile "${PROFILE}"
    else
        aws configure -p "${PROFILE}"
    fi
else
    aws configure -p "${PROFILE}"
fi

aws eks --region "${REGION}" update-kubeconfig --name "${CLUSTER}"
