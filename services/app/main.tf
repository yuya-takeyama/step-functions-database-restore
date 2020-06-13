module "restore-develop" {
  source = "../../modules/restore-aurora"

  db_identifier_prefix = "database"
  db_environment       = "develop"


  engine            = "aurora-postgresql"
  engine_version    = "11.6"
  db_instance_class = "db.t3.medium"

  db_subnet_group_name   = "default"
  vpc_security_group_ids = ["sg-8d878fef"]
  availability_zones     = ["ap-northeast-1a"]

  master_user_password = "postgresql"
}

module "restore-edge" {
  source = "../../modules/restore-aurora"

  db_identifier_prefix = "database"
  db_environment       = "edge"

  engine            = "aurora-postgresql"
  engine_version    = "11.6"
  db_instance_class = "db.t3.medium"

  db_subnet_group_name   = "default"
  vpc_security_group_ids = ["sg-8d878fef"]
  availability_zones     = ["ap-northeast-1a"]

  master_user_password = "postgresql"
}
