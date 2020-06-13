variable "task_parameters" {
  default = {
    SourceClusterIdentifier       = "database-production",
    DestinationClusterIdentifier  = "database-edge",
    DestinationInstanceIdentifier = "database-edge-a01",
    Engine                        = "aurora-postgresql",
    EngineVersion                 = "11.6",
    DBInstanceClass               = "db.t3.medium",
    AvailabilityZones = [
      "ap-northeast-1"
    ],
    DBSubnetGroupName = "default",
    VpcSecurityGroupIds = [
      "sg-8d878fef"
    ],
    MasterUserPassword = "postgresql"
  }
}

variable "retry_rule" {
  default = [
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

resource "aws_sfn_state_machine" "restore-edge" {
  name     = "restore-database-edge"
  role_arn = "arn:aws:iam::711930837542:role/service-role/StepFunctions-DatabaseRestore-role-ed2ad9da"

  definition = jsonencode({
    Comment = "Restore test",
    StartAt = "Delete",
    States = {
      Delete = {
        Type       = "Task",
        Resource   = "arn:aws:lambda:ap-northeast-1:711930837542:function:step-functions-database-restore-dev-delete",
        Next       = "WaitDeleted",
        Parameters = var.task_parameters
      },
      WaitDeleted = {
        Type       = "Task",
        Resource   = "arn:aws:lambda:ap-northeast-1:711930837542:function:step-functions-database-restore-dev-waitDeleted",
        Next       = "Restore",
        Retry      = var.retry_rule,
        Parameters = var.task_parameters
      },
      Restore = {
        Type       = "Task",
        Resource   = "arn:aws:lambda:ap-northeast-1:711930837542:function:step-functions-database-restore-dev-restore",
        Next       = "WaitAvailable",
        Parameters = var.task_parameters
      },
      WaitAvailable = {
        Type       = "Task",
        Resource   = "arn:aws:lambda:ap-northeast-1:711930837542:function:step-functions-database-restore-dev-waitAvailable",
        Next       = "Modify",
        Retry      = var.retry_rule,
        Parameters = var.task_parameters
      },
      Modify = {
        Type       = "Task",
        Resource   = "arn:aws:lambda:ap-northeast-1:711930837542:function:step-functions-database-restore-dev-modify",
        Next       = "WaitAvailableAfterModify",
        Parameters = var.task_parameters
      },
      WaitAvailableAfterModify = {
        Type       = "Task",
        Resource   = "arn:aws:lambda:ap-northeast-1:711930837542:function:step-functions-database-restore-dev-waitAvailable",
        End        = true,
        Retry      = var.retry_rule,
        Parameters = var.task_parameters
      }
    }
  })
}
