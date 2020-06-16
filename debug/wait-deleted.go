package main

import (
	"fmt"

	"github.com/k0kubun/pp"
	"github.com/yuya-takeyama/step-functions-database-restore/pkg"
)

func main() {
	fmt.Println("DEBUG start")
	err := pkg.WaitDeleted(&pkg.RestoreInput{
		DestinationClusterIdentifier:  "database-develop",
		DestinationInstanceIdentifier: "database-develop-a01",
	})

	pp.Println(err)

	fmt.Println("DEBUG end")
}
