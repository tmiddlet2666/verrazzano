# Copyright (c) 2020, 2021, Oracle and/or its affiliates.
# Licensed under the Universal Permissive License v 1.0 as shown at https://oss.oracle.com/licenses/upl.

variable "compartment_id" {}

variable "cluster_name" {
  default = "oke"
}
variable "label_prefix" {
  default = ""
}
variable "username" {
  default = ""
}
variable "tenancy_name" {
  default = "stevengreenberginc"
}

variable "operating_system_version" {
  default     = "8"
}

variable "kubernetes_version" {
  default = "v1.20.8"
}
variable "allow_worker_ssh_access" {
  default = false
}
variable "worker_type" {
  default = "private"
}
variable "control_plane_type" {
  default = "public"
}
variable "ssh_public_key_path" {
  default = ""
}
variable "ssh_private_key_path" {
  default = ""
}
variable "node_pools" {
  default = {"np1" = {shape="VM.Standard2.4",node_pool_size=4,boot_volume_size=50}}
}
variable "allow_node_port_access" {
  default = false
}
variable "create_bastion_host" {
  default = false
}
variable "create_operator" {
  default = false
}
variable "node_pool_image_id" {}
variable "enable_calico" {}
variable "calico_version" {}
