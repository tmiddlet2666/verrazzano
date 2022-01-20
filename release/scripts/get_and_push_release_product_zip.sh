#!/usr/bin/env bash
#
# Copyright (c) 2022, Oracle and/or its affiliates.
# Licensed under the Universal Permissive License v 1.0 as shown at https://oss.oracle.com/licenses/upl.
#
# Get the product zip uploaded by the periodic job for the given commit hash, and push it to the
# release-specific location

if [ -z "$VERRAZZANO_DEV_VERSION" ] || [ -z "$SHORT_COMMIT_HASH" ] || [ -z "$CLEAN_BRANCH_NAME" ]; then
  echo "Environment variables VERRAZZANO_DEV_VERSION, SHORT_COMMIT_HASH and CLEAN_BRANCH_NAME must be set"
  exit 1
fi

if [ -z "$OCI_OS_NAMESPACE" ] || [ -z "$OCI_OS_BUCKET" ] || [ -z "$WORKSPACE" ]; then
  echo "Environment variables OCI_OS_NAMESPACE, OCI_OS_BUCKET and WORKSPACE must be set"
  exit 2
fi

PERIODIC_ZIPFILE_OBJECT_STORE="${CLEAN_BRANCH_NAME}/${SHORT_COMMIT_HASH}/verrazzano.${VERRAZZANO_DEV_VERSION}.zip"
DOWNLOADED_ZIPFILE="${WORKSPACE}/verrazzano.${VERRAZZANO_DEV_VERSION}.zip"
RELEASE_ZIPFILE_OBJECT_STORE="${CLEAN_BRANCH_NAME}-release/${VERRAZZANO_DEV_VERSION}/verrazzano.zip"

echo "Downloading ${PERIODIC_ZIPFILE_OBJECT_STORE} from object store"
# get the product zip from the location where periodic tests would have pushed it
oci --region us-phoenix-1 os object get --namespace ${OCI_OS_NAMESPACE} -bn ${OCI_OS_BUCKET} --name ${PERIODIC_ZIPFILE_OBJECT_STORE} --file ${DOWNLOADED_ZIPFILE}

echo "Uploading the downloaded file to object store as ${RELEASE_ZIPFILE_OBJECT_STORE}"
# upload product zip in the location for the release
oci --region us-phoenix-1 os object put --force --namespace ${OCI_OS_NAMESPACE} -bn ${OCI_OS_BUCKET} --name ${RELEASE_ZIPFILE_OBJECT_STORE} --file ${DOWNLOADED_ZIPFILE}