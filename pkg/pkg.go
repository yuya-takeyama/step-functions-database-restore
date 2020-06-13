package pkg

import (
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/rds"
)

// ResoteInput is the input of DeleteOldClusterAndInstance
type ResoteInput struct {
	SourceClusterIdentifier       string
	DestinationClusterIdentifier  string
	DestinationInstanceIdentifier string
	Engine                        string
	EngineVersion                 string
	DBInstanceClass               string
	AvailabilityZones             []string
	DBSubnetGroupName             string
	VpcSecurityGroupIds           []string
	MasterUserPassword            string
}

// Delete deletes the old cluster and the instance
// if they exist
func Delete(input *ResoteInput) error {
	svc := rds.New(session.New())

	deleteInstanceInput := &rds.DeleteDBInstanceInput{
		DBInstanceIdentifier: aws.String(input.DestinationInstanceIdentifier),
		SkipFinalSnapshot:    aws.Bool(true),
	}
	if _, err := svc.DeleteDBInstance(deleteInstanceInput); err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case rds.ErrCodeDBInstanceNotFoundFault:
				fmt.Println("DB Instance does not exist")

			default:
				return fmt.Errorf("Failed to delete DB Instance: %s", err)
			}
		}
	}

	deleteClusterInput := &rds.DeleteDBClusterInput{
		DBClusterIdentifier: aws.String(input.DestinationClusterIdentifier),
		SkipFinalSnapshot:   aws.Bool(true),
	}
	if _, err := svc.DeleteDBCluster(deleteClusterInput); err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case rds.ErrCodeDBClusterNotFoundFault:
				fmt.Println("DB Cluster does not exist")

			default:
				return fmt.Errorf("Failed to delete DB Cluster: %s", err)
			}
		}
	}

	return nil
}

// WaitDeleted returns error unless the old cluster and instance are deleted
func WaitDeleted(input *ResoteInput) error {
	var instanceDeleted bool
	var clusterDeleted bool

	svc := rds.New(session.New())

	describeInstancesInput := &rds.DescribeDBInstancesInput{
		DBInstanceIdentifier: aws.String(input.DestinationInstanceIdentifier),
	}
	if _, err := svc.DescribeDBInstances(describeInstancesInput); err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case rds.ErrCodeDBInstanceNotFoundFault:
				fmt.Println("DB Instance is already deleted")
				instanceDeleted = true

			default:
				return fmt.Errorf("Failed to get the status of DB Instance: %s", err)
			}
		}
	}

	describeClustersInput := &rds.DescribeDBClustersInput{
		DBClusterIdentifier: aws.String(input.DestinationClusterIdentifier),
	}
	if _, err := svc.DescribeDBClusters(describeClustersInput); err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case rds.ErrCodeDBClusterNotFoundFault:
				fmt.Println("DB Cluster is already deleted")
				clusterDeleted = true

			default:
				return fmt.Errorf("Failed to get the status of DB Cluster: %s", err)
			}
		}
	}

	if !instanceDeleted {
		return fmt.Errorf("DB Instance is still not deleted")
	}
	if !clusterDeleted {
		return fmt.Errorf("DB Cluster is still not deleted")
	}

	return nil
}

// Restore restores from a snapshot
func Restore(input *ResoteInput) error {
	svc := rds.New(session.New())

	describeSnapshotsInput := &rds.DescribeDBClusterSnapshotsInput{
		DBClusterIdentifier: aws.String(input.SourceClusterIdentifier),
	}
	describeSnapshotsOutput, err := svc.DescribeDBClusterSnapshots(describeSnapshotsInput)
	if err != nil {
		return fmt.Errorf("Failed get the DB Cluster Snapshots: %s", err)
	}

	if len(describeSnapshotsOutput.DBClusterSnapshots) == 0 {
		return fmt.Errorf("There is no DB Cluster Snapshots")
	}

	latestTime := time.Time{}
	var latestSnapshot *rds.DBClusterSnapshot
	for _, snapshot := range describeSnapshotsOutput.DBClusterSnapshots {
		if *snapshot.Status != "available" {
			continue
		}
		if snapshot.SnapshotCreateTime.After(latestTime) {
			latestTime = *snapshot.SnapshotCreateTime
			latestSnapshot = snapshot
		}
	}
	fmt.Printf("Latest snapshot = %s\n", *latestSnapshot.DBClusterSnapshotIdentifier)

	restoreInput := &rds.RestoreDBClusterFromSnapshotInput{
		DBClusterIdentifier: aws.String(input.DestinationClusterIdentifier),
		SnapshotIdentifier:  latestSnapshot.DBClusterSnapshotIdentifier,
		Engine:              aws.String(input.Engine),
		EngineVersion:       aws.String(input.EngineVersion),
		AvailabilityZones:   aws.StringSlice(input.AvailabilityZones),
		DBSubnetGroupName:   aws.String(input.DBSubnetGroupName),
		VpcSecurityGroupIds: aws.StringSlice(input.VpcSecurityGroupIds),
	}
	if _, err := svc.RestoreDBClusterFromSnapshot(restoreInput); err != nil {
		return fmt.Errorf("Failed to restore: %s", err)
	}

	createInstanceInput := &rds.CreateDBInstanceInput{
		DBClusterIdentifier:  aws.String(input.DestinationClusterIdentifier),
		DBInstanceIdentifier: aws.String(input.DestinationInstanceIdentifier),
		DBInstanceClass:      aws.String(input.DBInstanceClass),
		Engine:               aws.String(input.Engine),
	}
	if _, err := svc.CreateDBInstance(createInstanceInput); err != nil {
		return fmt.Errorf("Failed to create DB Instance: %s", err)
	}

	return nil
}

// WaitAvailable returns error unless the destination cluster and instance are available
func WaitAvailable(input *ResoteInput) error {
	svc := rds.New(session.New())

	describeInstancesInput := &rds.DescribeDBInstancesInput{
		DBInstanceIdentifier: aws.String(input.DestinationInstanceIdentifier),
	}
	describeInstancesOutput, err := svc.DescribeDBInstances(describeInstancesInput)
	if err != nil {
		return fmt.Errorf("Failed to get the DB Instance: %s", err)
	}
	if len(describeInstancesOutput.DBInstances) < 1 {
		return fmt.Errorf("There is no DB Instances")
	}

	instance := describeInstancesOutput.DBInstances[0]
	if *instance.DBInstanceStatus != "available" {
		return fmt.Errorf("DB Instance is still not available: %s", *instance.DBInstanceStatus)
	}

	describeClustersInput := &rds.DescribeDBClustersInput{
		DBClusterIdentifier: aws.String(input.DestinationClusterIdentifier),
	}
	describeClustersOutput, err := svc.DescribeDBClusters(describeClustersInput)
	if err != nil {
		return fmt.Errorf("Failed to get the DB Cluster: %s", err)
	}
	if len(describeClustersOutput.DBClusters) < 1 {
		return fmt.Errorf("There is no DB Clusters")
	}

	cluster := describeClustersOutput.DBClusters[0]
	if *cluster.Status != "available" {
		return fmt.Errorf("DB Cluster is still not available: %s", *cluster.Status)
	}

	return nil
}

// Modify modifies DB Cluster
func Modify(input *ResoteInput) error {
	svc := rds.New(session.New())

	modifyInput := &rds.ModifyDBClusterInput{
		DBClusterIdentifier: aws.String(input.DestinationClusterIdentifier),
		MasterUserPassword:  aws.String(input.MasterUserPassword),
		VpcSecurityGroupIds: aws.StringSlice(input.VpcSecurityGroupIds),
		ApplyImmediately:    aws.Bool(true),
	}
	if _, err := svc.ModifyDBCluster(modifyInput); err != nil {
		return fmt.Errorf("Failed to modify DB Cluster: %s", err)
	}

	return nil
}
