#!/bin/bash
#
# Copyright (c) 2020, 2021, Oracle and/or its affiliates.
# Licensed under the Universal Permissive License v 1.0 as shown at https://oss.oracle.com/licenses/upl.
#

CURR_SCRIPT_DIR=$(cd $(dirname "$0"); pwd -P)
. $CURR_SCRIPT_DIR/../cluster/init.sh

$SCRIPT_DIR/terraform destroy -auto-approve -no-color
