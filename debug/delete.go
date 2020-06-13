package main

import (
	"fmt"

	"github.com/k0kubun/pp"
	"github.com/yuya-takeyama/step-functions-database-restore/pkg"
)

func main() {
	fmt.Println("DEBUG start")
	err := pkg.Delete(&pkg.ResoteInput{
		OldDBClusterIdentifier:  "hoge-cluster",
		OldDBInstanceIdentifier: "hoge-cluster-a01",
	})

	pp.Println(err)

	fmt.Println("DEBUG end")
}
