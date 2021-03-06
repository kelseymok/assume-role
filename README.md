# Assume Role
The purpose of this project is to create a development environment which comes with an easy way to assume roles. The assumption made here is that you have a repository containing code which needs an AWS Profile or assumed role to run (e.g. Terraform code, AWS Cloudformation templates). This project means to be used as a submodule as the run scripts mount your project directory into the container.

## Requirements
* An IAM User account with at least one AccessKey and SecretKeyID
* Docker
* A device with an MFA Authenticator

## CLI Access
1. Set up your AWS CLI credentials according to the [AWS CLI docs](https://docs.aws.amazon.com/cli/latest/userguide/cli-chap-configure.html). 
2. [Set up a MFA](https://docs.aws.amazon.com/IAM/latest/UserGuide/id_credentials_mfa_enable_virtual.html#enable-virt-mfa-for-iam-user). 
3. Pick a profile name that represents your user in your organization. Recommendation: `<your-organisation-name>-base`. For example: `penguin-base`.
4. Configure the AWS CLI Credentials with details from the AWS Console > My Security Credentials.
```bash
export base_profile=penguin-base

aws configure --profile ${base_profile} set aws_access_key_id <your Access Key ID>
aws configure --profile ${base_profile} set aws_secret_access_key <your Secret Access Key>
aws configure --profile ${base_profile} set region <some-region:eu-central-1>
aws configure --profile ${base_profile} set mfa_serial <the ARN of your MFA device>
```

## Assume a Role
In the root directory of your project, build the image:
```bash
export base_profile=penguin-base
export role=myOrg/developer

docker run -it \
  -v "$(cd ~; pwd)/.aws:/root/.aws" \
  --entrypoint="" \
  kb1rd/assume-role assume-role ${base_profile} ${role}
```
This puts a profile called `myOrg-developer` (note the `/` to `-` substitution -- it's a feature) in your AWS Credentials (`~/.aws/credentials`) to accommodate for roles that have a path. For roles without a path, the profile name is `developer`.

## Use it
### AWS CLI
```bash
aws sts get-caller-identity --profile myOrg-developer
```
 
### Terraform
In your provider:
```hcl-terraform
provider "aws" {
  region = "eu-central-1"
  profile = "myOrg-developer"
}
```

## Interactive Usage
This image comes packed with Terraform `0.12.24` and can be used as a development environment:
```bash
docker run -it \
  -v "$(cd ~; pwd)/.aws:/root/.aws" \
  -v "$(cd ~; pwd)/PATH-TO-PROJECT:/app" \
  kb1rd/assume-role bash
```