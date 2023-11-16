# cloud-resource-collector
Collects cloud resources in a given account. Supports multiple cloud providers.

## Prerequisites

### Setup the AWS Collector

The AWS Collector requires you to provide region and credential information. You can do this either by setting up 
shared credential and config files or by setting environment variables.

To setup the credential and config files do the following:

- Create a text file with the following content (replacing the keys with your AWS keys)
```ini
[default]
aws_access_key_id = YOUR_AWS_ACCESS_KEY_ID
aws_secret_access_key = YOUR_AWS_SECRET_ACCESS_KEY
```
If you are using Windows save the file under `C:\Users\<yourUserName>\.aws\credentials`.
If you are using Linux, MacOS, or Unix save the file under `~/.aws/credentials`

- Create a text file with the following content (choosing the appropriate region)
```ini
[default]
region = eu-north-1
output = json
```
If you are using Windows save the file under `C:\Users\<yourUserName>\.aws\config`.
If you are using Linux, MacOS, or Unix save the file under `~/.aws/config`

Alternatively, you can set the following environment files:
```shell
export AWS_ACCESS_KEY_ID=YOUR_AWS_ACCESS_KEY_ID
export AWS_SECRET_ACCESS_KEY=YOUR_AWS_SECRET_ACCESS_KEY
export AWS_REGION=eu-north-1
```

Note: Pagination is not yet implemented, the collector will return only the first page of resources.

### Setup the IBM Collector

The IBM collector requires an IBM API key to be supplied through the following environment variable:
```shell
export IBMCLOUD_API_KEY=<ibm-cloud-api-key>
```

## Usage

```
$ ./bin/collect -h
Usage of ./bin/collect:
  -out string
        file path to store results
  -provider string
        cloud provider from which to collect resources
  -region value
        cloud region from which to collect resources
```
* Value of `-provider` must be either `ibm` or `aws`
* The `-region` argument can appear multiple times. If running with no `-region` arguments, resources from all regions are collected.

## Build the project

```shell
git clone git@github.com:np-guard/cloud-resource-collector.git
cd cloud-resource-collector
make mod
make build
```

