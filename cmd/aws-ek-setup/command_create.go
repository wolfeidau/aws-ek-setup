package main

import (
	"bytes"
	"fmt"
	"os"
	"regexp"
	"text/template"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/elasticsearchservice"
	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/spf13/cobra"
)

var (
	cmdCreateCluster = &cobra.Command{
		Use:   "create",
		Short: "Create elasticsearch cluster.",
		Long:  ``,
		Run:   runCmdCreateCluster,
	}
	clusterDomainName    string
	clusterInstanceType  string
	clusterInstanceCount int
	clusterVolumeSize    int

	arnRegex = regexp.MustCompile(`arn:aws:iam::(?P<aws_account_id>\d+):.*`)

	pubSubToAnyTopicPolicy = `{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "Principal": {
        "AWS": "{{.UserARN}}"
      },
      "Action": "es:*",
      "Resource": "arn:aws:es:us-east-1:{{.AccountID}}:{{.DomainName}}/logs/*"
    }
  ]
}`
)

func init() {
	cmdCreateCluster.Flags().StringVar(&clusterDomainName, "name", "", "The name of the cluster domain.")
	cmdCreateCluster.Flags().StringVar(&clusterInstanceType, "instance-type", "t2.small.elasticsearch", "The instance type used in the cluster domain.")
	cmdCreateCluster.Flags().IntVar(&clusterInstanceCount, "count", 1, "The cluster instance count.")
	cmdCreateCluster.Flags().IntVar(&clusterVolumeSize, "size", 40, "The size of the disks in gigabytes.")

	cmdRoot.AddCommand(cmdCreateCluster)
}

func runCmdCreateCluster(cmd *cobra.Command, args []string) {

	// get the identify of the current user
	resp, err := iamSvc.GetUser(&iam.GetUserInput{})

	if err != nil {
		stderr("Failed to retrieve user information from IAM: %v", err)
		os.Exit(1)
	}

	accountID := getAccountIdentifier(*resp.User.Arn)

	// generate the policy document
	targs := struct {
		UserARN    string
		AccountID  string
		DomainName string
	}{
		*resp.User.Arn,
		accountID,
		clusterDomainName,
	}

	var b bytes.Buffer

	t := template.Must(template.New("policy").Parse(pubSubToAnyTopicPolicy))

	err = t.Execute(&b, targs)

	if err != nil {
		stderr("Failed to retrieve user information from IAM: %v", err)
		os.Exit(1)
	}

	// create an elastic search cluster

	params := &elasticsearchservice.CreateElasticsearchDomainInput{
		DomainName:     aws.String(clusterDomainName),
		AccessPolicies: aws.String(b.String()),
		EBSOptions: &elasticsearchservice.EBSOptions{
			EBSEnabled: aws.Bool(true),
			VolumeSize: aws.Int64(int64(clusterVolumeSize)),
			VolumeType: aws.String("gp2"),
		},
		ElasticsearchClusterConfig: &elasticsearchservice.ElasticsearchClusterConfig{
			InstanceCount:        aws.Int64(int64(clusterInstanceCount)),
			InstanceType:         aws.String(clusterInstanceType),
			ZoneAwarenessEnabled: aws.Bool(false),
		},
	}

	esresp, err := elasticSearchSvc.CreateElasticsearchDomain(params)

	if err != nil {
		stderr("Failed to create elastic search cluster: %v", err)
		os.Exit(1)
	}

	fmt.Println("DomainStatus: ", esresp.DomainStatus.String())

}

// aws es create-elasticsearch-domain --domain-name weblogs --elasticsearch-cluster-config InstanceType=m3.large.elasticsearch,InstanceCount=5,ZoneAwarenessEnabled=true --ebs-options EBSEnabled=true,VolumeType=gp2,VolumeSize=100 --access-policies ''

func getAccountIdentifier(arn string) string {

	n1 := arnRegex.SubexpNames()
	r2 := arnRegex.FindAllStringSubmatch(arn, -1)[0]

	md := map[string]string{}
	for i, n := range r2 {
		md[n1[i]] = n
	}

	return md["aws_account_id"]
}
