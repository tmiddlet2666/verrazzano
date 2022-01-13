#!/bin/bash
# Copyright (c) 2022, Oracle and/or its affiliates.
# Licensed under the Universal Permissive License v 1.0 as shown at https://oss.oracle.com/licenses/upl.

set -x

SCRIPT_DIR=$(cd $(dirname "$0"); pwd -P)
$SCRIPT_DIR/init.sh

pushd $1
$SCRIPT_DIR/terraform workspace select $2
$SCRIPT_DIR/terraform init
$SCRIPT_DIR/terraform destroy -auto-approve -no-color
