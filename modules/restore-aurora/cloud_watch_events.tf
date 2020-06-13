resource "aws_cloudwatch_event_rule" "trigger_restore" {
  name                = "trigger-restore-aurora-${var.db_identifier_prefix}-${var.db_environment}"
  description         = "Trigger restore-aurora-${var.db_identifier_prefix}-${var.db_environment} periodically"
  schedule_expression = var.trigger_schedule
}

resource "aws_cloudwatch_event_target" "trigger-restore" {
  target_id = "step-functions"
  rule      = aws_cloudwatch_event_rule.trigger_restore.name
  arn       = aws_sfn_state_machine.restore.id
  role_arn  = "arn:aws:iam::711930837542:role/cloudwatch-events-trigger-restore-aurora"
}