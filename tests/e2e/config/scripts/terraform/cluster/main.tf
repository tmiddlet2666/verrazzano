# Copyright (c) 2020, 2022, Oracle and/or its affiliates.
# Licensed under the Universal Permissive License v 1.0 as shown at https://oss.oracle.com/licenses/upl.

module "oke" {
  source = "oracle-terraform-modules/oke/oci"
  version = "4.2.2"

  tenancy_id = var.tenancy_id
  user_id = var.user_id
  region = var.region
  home_region = var.home_region
  api_fingerprint = var.api_fingerprint
  api_private_key_path =var.api_private_key_path

  cluster_name = var.cluster_name
  compartment_id = var.compartment_id
  kubernetes_version = var.kubernetes_version
  allow_worker_ssh_access = var.allow_worker_ssh_access
  worker_type = var.worker_type
  control_plane_type = var.control_plane_type
  node_pools =var.node_pools
  allow_node_port_access = var.allow_node_port_access
  create_operator = var.create_operator
  create_bastion_host = var.create_bastion_host
  username = var.username

  enable_calico = var.enable_calico
  calico_version = var.calico_version

  vcn_name = "${var.cluster_name}-vcn"
  vcn_dns_label = var.cluster_name
  label_prefix = var.label_prefix

  email_address = ""

  create_service_account = false
  service_account_cluster_role_binding = ""

  use_signed_images = false
  image_signing_keys = []

  node_pool_image_id = var.node_pool_image_id

  providers = {
    oci.home = oci.home
  }
}
