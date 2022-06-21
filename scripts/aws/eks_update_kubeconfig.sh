#/bin/bash
#--------------------------------------------------------------
# Set kubeconfig
#--------------------------------------------------------------
REGION=$1
CLUSTER=$2

aws eks --region "${REGION}" update-kubeconfig --name "${CLUSTER}"
