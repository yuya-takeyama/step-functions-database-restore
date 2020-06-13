#!/bin/bash

set -eu
set -o pipefail

STAGING_DB_CLUSTER_IDENTIFIER="database-develop"
STAGING_DB_INSTANCE_IDENTIFIER="${STAGING_DB_CLUSTER_IDENTIFIER}-a01"
ENGINE=aurora-postgresql
ENGINE_VERSION=11.6
STAGING_DB_INSTANCE_CLASS="db.t3.medium"
AVAILABILITY_ZONES=ap-northeast-1a
DB_SUBNET_GROUP_NAME=default
SECURITY_GROUP_ID_DEFAULT=sg-8d878fef

CLUSTER_SNAPSHOT_ID="$(aws --profile yuya rds describe-db-cluster-snapshots --db-cluster-identifier database-production --query 'reverse(sort_by(DBClusterSnapshots[],&SnapshotCreateTime))' | jq -r '.[0].DBClusterSnapshotIdentifier')"
aws --profile yuya --region ap-northeast-1 rds restore-db-cluster-from-snapshot \
  --db-cluster-identifier "${STAGING_DB_CLUSTER_IDENTIFIER}" \
  --snapshot-identifier "${CLUSTER_SNAPSHOT_ID}" \
  --engine "${ENGINE}" \
  --engine-version "${ENGINE_VERSION}" \
  --availability-zones "${AVAILABILITY_ZONES}" \
  --db-subnet-group-name "${DB_SUBNET_GROUP_NAME}" \
  --vpc-security-group-ids "${SECURITY_GROUP_ID_DEFAULT}" || :

aws --profile yuya --region ap-northeast-1 rds create-db-instance \
  --db-instance-identifier "${STAGING_DB_INSTANCE_IDENTIFIER}" \
  --db-instance-class "${STAGING_DB_INSTANCE_CLASS}" \
  --engine ${ENGINE} \
  --db-cluster-identifier "${STAGING_DB_CLUSTER_IDENTIFIER}"