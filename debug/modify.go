package main

import (
	"fmt"

	"github.com/k0kubun/pp"
	"github.com/yuya-takeyama/step-functions-database-restore/pkg"
)

func main() {
	fmt.Println("DEBUG start")
	err := pkg.Modify(&pkg.RestoreInput{
		SourceClusterIdentifier:       "database-production",
		DestinationClusterIdentifier:  "database-develop",
		DestinationInstanceIdentifier: "database-develop-a01",
		Engine:                        "aurora-postgresql",
		EngineVersion:                 "11.6",
		DBInstanceClass:               "db.t3.medium",
		MasterUserPassword:            "postgres",
	})

	pp.Println(err)

	fmt.Println("DEBUG end")
}
