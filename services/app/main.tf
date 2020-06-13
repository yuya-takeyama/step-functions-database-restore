module "restore-develop" {
  source = "../../modules/restore-aurora"

  state_machine_name = "restore-develop"

  source_cluster_identifier       = "database-production"
  destination_cluster_identifier  = "database-develop"
  destination_instance_identifier = "database-develop-a01"

  engine            = "aurora-postgresql"
  engine_version    = "11.6"
  db_instance_class = "t3.medium"

  db_subnet_group_name   = "default"
  vpc_security_group_ids = ["sg-8d878fef"]
  availability_zones     = ["ap-northeast-1a"]

  master_user_password = "postgresql"
}

module "restore-edge" {
  source = "../../modules/restore-aurora"

  state_machine_name = "restore-edge"

  source_cluster_identifier       = "database-production"
  destination_cluster_identifier  = "database-edge"
  destination_instance_identifier = "database-edge-a01"

  engine            = "aurora-postgresql"
  engine_version    = "11.6"
  db_instance_class = "t3.medium"

  db_subnet_group_name   = "default"
  vpc_security_group_ids = ["sg-8d878fef"]
  availability_zones     = ["ap-northeast-1a"]

  master_user_password = "postgresql"
}
