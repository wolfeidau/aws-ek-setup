# aws-ek-setup

This tool uses the [aws-go-sdk](https://github.com/aws/aws-sdk-go) to bootstrap an [elasticsearch](https://www.elastic.co/products/elasticsearch) cluster with [cloudwatch](https://aws.amazon.com/cloudwatch/) logs configured to forward to this cluster.

# prerequisites

To use this you need an AWS account with some services logging to cloudwatch logs.

# building

```
make
```

# usage

```
Create elasticsearch cluster.

Usage:
  aws-ek-setup create [flags]

Flags:
      --count=1: The cluster instance count.
      --instance-type="t2.small.elasticsearch": The instance type used in the cluster domain.
      --name="": The name of the cluster domain.
      --size=40: The size of the disks in gigabytes.


Global Flags:
      --aws-debug[=false]: Log debug information from aws-sdk-go library

```

Brief example which creates a new cluster named `testcluster`.

```
AWS_PROFILE=XXX aws-ek-setup create --name=testcluster
```

# warning

This is a work in progress at the moment and has some pretty basic defaults right now.

# todo

* Discover all cloudwatch log groups and stream them to elastic search

# License

aws-ek-setup is Copyright (c) 2015 Mark Wolfe @wolfeidau and licensed under the MIT license. All rights not explicitly granted in the MIT license are reserved. See the included LICENSE.md file for more details.

