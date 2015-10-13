# aws-ek-setup

This tool uses the [aws-go-sdk](https://github.com/aws/aws-sdk-go) to bootstrap an [elasticsearch](https://www.elastic.co/products/elasticsearch) cluster with [cloudwatch](https://aws.amazon.com/cloudwatch/) logs configured to forward to this cluster.

# prerequisites

To use this you need an AWS account with some services logging to cloudwatch logs.

# building

```
make
```

# usage

Firstly create a `cluster.yml` in the current directory using the example provided.

```
Create an elasticsearch cluster.

Usage:
  aws-ek-setup up [flags]

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

