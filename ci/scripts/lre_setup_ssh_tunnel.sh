#!/usr/bin/env bash
#
# Copyright (c) 2021, Oracle and/or its affiliates.
# Licensed under the Universal Permissive License v 1.0 as shown at https://oss.oracle.com/licenses/upl.
#

if [ -z "${ssh_private_key_path}" ] ; then
    echo "ssh_private_key_path env var must be set!"
    exit 1
fi
if [ -z "${ssh_public_key_path}" ] ; then
    echo "ssh_public_key_path env var must be set!"
    exit 1
fi
if [ -z "${dev_lre_compartment_id}" ] ; then
    echo "dev_lre_compartment_id env var must be set!"
    exit 1
fi
if [ -z "${KUBECONFIG}" ] ; then
    echo "KUBECONFIG env var must be set!"
    exit 1
fi


BASTION_ID=$(oci bastion bastion list \
            --compartment-id "${dev_lre_compartment_id}" --all \
            | jq -r '.data[0]."id"')

SESSION_ID=$(oci bastion session create-port-forwarding \
   --bastion-id $BASTION_ID \
   --display-name br-test-pf-session \
   --ssh-public-key-file ${ssh_public_key_path} \
   --key-type PUB \
   --target-private-ip 10.196.0.58 \
   --target-port 6443)

echo "Waiting for $SESSION_ID to start"
sleep 15

COMMAND=`oci bastion session get  --session-id=${SESSION_ID} | \
  jq '.data."ssh-metadata".command' | \
  sed 's/"//g' | \
  sed 's|<privateKey>|${ssh_private_key_path}|g' | \
  sed 's|<localPort>|6443|g'`
echo ${COMMAND}
eval ${COMMAND}
