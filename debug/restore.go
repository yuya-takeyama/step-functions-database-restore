package main

import (
	"fmt"

	"github.com/k0kubun/pp"
	"github.com/yuya-takeyama/step-functions-database-restore/pkg"
)

func main() {
	fmt.Println("DEBUG start")
	err := pkg.Restore(&pkg.RestoreInput{
		SourceClusterIdentifier:       "database-production",
		DestinationClusterIdentifier:  "database-edge2",
		DestinationInstanceIdentifier: "database-edge2-a01",
		Engine:                        "aurora-postgresql",
		EngineVersion:                 "11.6",
		DBInstanceClass:               "db.t3.medium",
	})

	pp.Println(err)

	fmt.Println("DEBUG end")
}
