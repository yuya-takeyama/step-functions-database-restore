locals {
  task_parameters = {
    SourceClusterIdentifier       = var.source_cluster_identifier,
    DestinationClusterIdentifier  = var.destination_cluster_identifier,
    DestinationInstanceIdentifier = var.destination_instance_identifier,
    Engine                        = var.engine,
    EngineVersion                 = var.engine_version,
    DBInstanceClass               = var.db_instance_class,
    AvailabilityZones             = var.availability_zones,
    DBSubnetGroupName             = var.db_subnet_group_name,
    VpcSecurityGroupIds           = var.vpc_security_group_ids,
    MasterUserPassword            = var.master_user_password
  }
  retry_rule = [
    {
      ErrorEquals = [
        "States.TaskFailed"
      ],
      IntervalSeconds = 60,
      MaxAttempts     = 60,
      BackoffRate     = 1
    }
  ]
}

resource "aws_sfn_state_machine" "restore" {
  name     = var.state_machine_name
  role_arn = "arn:aws:iam::711930837542:role/service-role/StepFunctions-DatabaseRestore-role-ed2ad9da"

  definition = jsonencode({
    Comment = "Restore test",
    StartAt = "Delete",
    States = {
      Delete = {
        Type       = "Task",
        Resource   = "arn:aws:lambda:ap-northeast-1:711930837542:function:step-functions-database-restore-dev-delete",
        Next       = "WaitDeleted",
        Parameters = local.task_parameters
      },
      WaitDeleted = {
        Type       = "Task",
        Resource   = "arn:aws:lambda:ap-northeast-1:711930837542:function:step-functions-database-restore-dev-waitDeleted",
        Next       = "Restore",
        Retry      = local.retry_rule,
        Parameters = local.task_parameters
      },
      Restore = {
        Type       = "Task",
        Resource   = "arn:aws:lambda:ap-northeast-1:711930837542:function:step-functions-database-restore-dev-restore",
        Next       = "WaitAvailable",
        Parameters = local.task_parameters
      },
      WaitAvailable = {
        Type       = "Task",
        Resource   = "arn:aws:lambda:ap-northeast-1:711930837542:function:step-functions-database-restore-dev-waitAvailable",
        Next       = "Modify",
        Retry      = local.retry_rule,
        Parameters = local.task_parameters
      },
      Modify = {
        Type       = "Task",
        Resource   = "arn:aws:lambda:ap-northeast-1:711930837542:function:step-functions-database-restore-dev-modify",
        Next       = "WaitAvailableAfterModify",
        Parameters = local.task_parameters
      },
      WaitAvailableAfterModify = {
        Type       = "Task",
        Resource   = "arn:aws:lambda:ap-northeast-1:711930837542:function:step-functions-database-restore-dev-waitAvailable",
        End        = true,
        Retry      = local.retry_rule,
        Parameters = local.task_parameters
      }
    }
  })
}
