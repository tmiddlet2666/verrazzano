#!/bin/bash
#
# Copyright (c) 2022, Oracle and/or its affiliates.
# Licensed under the Universal Permissive License v 1.0 as shown at https://oss.oracle.com/licenses/upl.
#

SCRIPT_DIR=$(cd $(dirname "$0"); pwd -P)
$SCRIPT_DIR/init.sh
$SCRIPT_DIR/terraform init $1 -no-color
$SCRIPT_DIR/terraform plan $1 -no-color

set -o pipefail

# retry 3 times, 30 seconds apart
tries=0
MAX_TRIES=3
while true; do
   tries=$((tries+1))
   echo "terraform apply iteration ${tries}"
   $SCRIPT_DIR/terraform apply $1 -auto-approve -no-color && break
   if [ "$tries" -ge "$MAX_TRIES" ];
   then
      echo "Terraform apply tries exceeded.  Cluster creation has failed!"
      break
   fi
   echo "Deleting Cluster Terraform and applying again"
   ./delete-vcn.sh $1
   sleep 30
done

if [ "$tries" -ge "$MAX_TRIES" ];
then
  exit 1
fi