package main

import (
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/elasticsearchservice"
	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/spf13/cobra"
)

var (
	// Version The version of the application (set by make file)
	Version       = "UNKNOWN"
	defaultRegion = "us-west-2"

	cmdRoot = &cobra.Command{
		Use:   "aws-ek-setup",
		Short: "Manage AWS Elasticsearch Clusters",
		Long:  ``,
	}

	rootOpts struct {
		AWSDebug bool
	}

	elasticSearchSvc = elasticsearchservice.New(newAWSConfig())
	iamSvc           = iam.New(newAWSConfig())
)

func init() {
	cmdRoot.PersistentFlags().BoolVar(&rootOpts.AWSDebug, "aws-debug", false, "Log debug information from aws-sdk-go library")
}

func main() {
	cmdRoot.Execute()
}

func newAWSConfig() *aws.Config {
	c := aws.NewConfig()
	c = c.WithRegion(defaultRegion)
	if rootOpts.AWSDebug {
		c = c.WithLogLevel(aws.LogDebug)
	}
	return c
}

func stderr(msg string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, msg+"\n", args...)
}
