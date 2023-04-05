# cloud-resource-collector
Collects cloud resources in a given account. Supports multiple cloud providers.

### Setup the AWS Collector

The AWS Collector requires you to provide region and credential information, by setting up shared credential and config files:

1. Create a text file with the following content (replacing the keys with your AWS keys)
```
[default]
aws_access_key_id = YOUR_AWS_ACCESS_KEY_ID
aws_secret_access_key = YOUR_AWS_SECRET_ACCESS_KEY
```
If you are using Windows save the file under `C:\Users\<yourUserName>\.aws\credentials`.
If you are using Linux, MacOS, or Unix save the file under `~/.aws/credentials`

2. Create a text file with the following content (choosing the appropriate region)
```
[default]
region = eu-north-1
output = json
```
If you are using Windows save the file under `C:\Users\<yourUserName>\.aws\config`.
If you are using Linux, MacOS, or Unix save the file under `~/.aws/config`

Note: Pagination is not yet implemented, the collector will return only the first page of resources.


### Executing

At the root directory run: 
```azure
make build
```

Then execute while supplying the requested provider and the output file name as arguments:
```azure
./bin/collect -provider aws -out b.json
```