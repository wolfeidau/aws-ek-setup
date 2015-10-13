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
	"github.com/spf13/viper"
)

var (
	cmdCreateCluster = &cobra.Command{
		Use:   "up",
		Short: "Create an elasticsearch cluster.",
		Long:  ``,
		Run:   runCmdCreateCluster,
	}

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
	cmdRoot.AddCommand(cmdCreateCluster)

	viper.SetConfigName("cluster")

	// setup defaults for the configuration
	viper.SetDefault("clusterName", "content")
	viper.SetDefault("clusterSize", 1)
	viper.SetDefault("instanceType", "t2.small.elasticsearch")
	viper.SetDefault("volumeSize", 20)
	viper.SetDefault("zoneAware", false)

}

func runCmdCreateCluster(cmd *cobra.Command, args []string) {

	err := viper.ReadInConfig()

	if err != nil {
		stderr("Failed to load cluster configuration file: %v", err)
		os.Exit(1)
	}

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
		viper.GetString("clusterName"),
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
		DomainName:     aws.String(viper.GetString("clusterName")),
		AccessPolicies: aws.String(b.String()),
		EBSOptions: &elasticsearchservice.EBSOptions{
			EBSEnabled: aws.Bool(true),
			VolumeSize: aws.Int64(int64(viper.GetInt("volumeSize"))),
			VolumeType: aws.String("gp2"),
		},
		ElasticsearchClusterConfig: &elasticsearchservice.ElasticsearchClusterConfig{
			InstanceCount:        aws.Int64(int64(viper.GetInt("clusterSize"))),
			InstanceType:         aws.String(viper.GetString("instanceType")),
			ZoneAwarenessEnabled: aws.Bool(viper.GetBool("zoneAware")),
		},
	}

	esresp, err := elasticSearchSvc.CreateElasticsearchDomain(params)

	if err != nil {
		stderr("Failed to create elastic search cluster: %v", err)
		os.Exit(1)
	}

	fmt.Println("DomainStatus: ", esresp.DomainStatus.String())
}

func getAccountIdentifier(arn string) string {

	n1 := arnRegex.SubexpNames()
	r2 := arnRegex.FindAllStringSubmatch(arn, -1)[0]

	md := map[string]string{}
	for i, n := range r2 {
		md[n1[i]] = n
	}

	return md["aws_account_id"]
}
