# cloud-resource-collector
Collects cloud resources in a given account. Supports multiple cloud providers.

## Prerequisites

### Setup the AWS Collector

The AWS Collector requires you to provide region and credential information by setting up shared credential and config files:

1. Create a text file with the following content (replacing the keys with your AWS keys)
```shell
[default]
aws_access_key_id = YOUR_AWS_ACCESS_KEY_ID
aws_secret_access_key = YOUR_AWS_SECRET_ACCESS_KEY
```
If you are using Windows save the file under `C:\Users\<yourUserName>\.aws\credentials`.
If you are using Linux, MacOS, or Unix save the file under `~/.aws/credentials`

2. Create a text file with the following content (choosing the appropriate region)
```shell
[default]
region = eu-north-1
output = json
```
If you are using Windows save the file under `C:\Users\<yourUserName>\.aws\config`.
If you are using Linux, MacOS, or Unix save the file under `~/.aws/config`

Note: Pagination is not yet implemented, the collector will return only the first page of resources.


## Usage

```shell
$ ./bin/collect -h
Usage of C:\MyStuff\Governance\np-guard\cloud-resource-collector\bin\collect:
  -out string
        file path to store results
  -provider string
        cloud provider from which to collect resources
```

## Build the project

```shell
git clone git@github.com:np-guard/cloud-resource-collector.git
cd cloud-resource-collector
go mod download
make
```

