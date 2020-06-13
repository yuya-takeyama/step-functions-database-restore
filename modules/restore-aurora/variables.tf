variable "state_machine_name" {
  type = string
}

variable "source_cluster_identifier" {
  type = string
}

variable "destination_cluster_identifier" {
  type = string
}

variable "destination_instance_identifier" {
  type = string
}

variable "engine" {
  type = string
}

variable "engine_version" {
  type = string
}

variable "db_instance_class" {
  type = string
}

variable "availability_zones" {
  type = list(string)
}

variable "db_subnet_group_name" {
  type = string
}

variable "vpc_security_group_ids" {
  type = list(string)
}

variable "master_user_password" {
  type = string
}
