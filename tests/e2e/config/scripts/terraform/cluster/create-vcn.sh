#!/bin/bash
#
# Copyright (c) 2022, Oracle and/or its affiliates.
# Licensed under the Universal Permissive License v 1.0 as shown at https://oss.oracle.com/licenses/upl.
#

SCRIPT_DIR=$(cd $(dirname "$0"); pwd -P)

moduleLocation=$1
if [ -z "$moduleLocation" ]; then
  echo "Module location must be specified"
fi


$SCRIPT_DIR/init.sh

pushd $moduleLocation

$SCRIPT_DIR/terraform init -no-color
$SCRIPT_DIR/terraform plan -no-color

set -o pipefail

set -x

# retry 3 times, 30 seconds apart
tries=0
MAX_TRIES=3
while true; do
   tries=$((tries+1))
   echo "terraform apply iteration ${tries}"
   $SCRIPT_DIR/terraform apply -auto-approve -no-color && break
   if [ "$tries" -ge "$MAX_TRIES" ];
   then
      echo "Terraform apply tries exceeded.  Cluster creation has failed!"
      break
   fi
   echo "Deleting Cluster Terraform and applying again"
   $SCRIPT_DIR/delete-vcn.sh $moduleLocation
   sleep 30
done

popd

if [ "$tries" -ge "$MAX_TRIES" ];
then
  exit 1
fi
