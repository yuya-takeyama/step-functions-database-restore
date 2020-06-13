package main

import (
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/yuya-takeyama/step-functions-database-restore/pkg"
)

func main() {
	lambda.Start(pkg.Restore)
}
