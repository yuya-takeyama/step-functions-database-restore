package pkg

import (
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/rds"
)

func TestPostgresURL(t *testing.T) {
	cluster := &rds.DBCluster{
		Endpoint:       aws.String("postgres-host.aws"),
		Port:           aws.Int64(int64(1234)),
		MasterUsername: aws.String("user"),
		DatabaseName:   aws.String("database"),
	}
	input := &RestoreInput{
		MasterUserPassword: "password",
	}

	expected := "postgresql://user:password@postgres-host.aws:1234/database"
	url := postgresURL(cluster, input)
	if url != expected {
		t.Errorf("Wrong: expected=%s, actual=%s", expected, url)
	}
}
