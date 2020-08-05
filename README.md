# Assume Role
The purpose of this project is to create a development environment which comes with an easy way to assume roles. The assumption made here is that you have a repository containing code which needs an AWS Profile or assumed role to run (e.g. Terraform code, AWS Cloudformation templates). This project means to be used as a submodule as the run scripts mount your project directory into the container.

## Requirements
* An IAM User account with at least one AccessKey and SecretKeyID
* Docker
* A device with an MFA Authenticator

## Using in your Repo
First-time usage:
```
git submodule add git@github.com:kelseymok/assume-role.git
git submodule update --init
```

Cloning a project with the submodule:
```
git submodule update --init
git pull --recurse-submodules && git submodule update --remote
```

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
./assume-role/go build
```

### Two ways to Assume a Role
There are two ways to use this docker image:
1. As a tool to assume a role: 
```bash
./assume-role/go assume-role <base_profile:penguin-base> <role-name:my-awesome-role> 
```
This creates/updates a profile named `<role-name>` in your Shared Credentials file (`~/.aws/credentials`) which you can use on your local machine by calling AWS cli commands with the `--profile <role-name>` flag (for example: `aws sts get-caller-identity --profile penguin-base-my-awesome-role`)
 
2. As a tool to assume a role and as a development environment (comes with Terraform and a few base libraries): 
```bash
./assume-role/go run

# Drops into a bash session
bash# assume-role <base-profile:penguin-base> <role-name:my-awesome-role>
``` 

Both methods create/update your Shared Credentials file (`~/.aws/credentials`) with a profile named `<role-name>` which you can use on your local machine or docker container by calling AWS cli commands with the `--profile <role-name>` flag (for example: `aws sts get-caller-identity --profile penguin-base-my-awesome-role`)

NOTE: If your AWS Role name contains a path, such as `penguins/developer`, the resulting profile in your Shared Credentials file will be `penguins-developer` (the `/` is converted to `-` to aesthetic reasons).

## Using Your Assumed Role
If you use Terraform, make sure to add a profile to your provider configuration:
```hcl-terraform
provider "aws" {
  region = "eu-central-1"
  profile = "Your Role Name"
}
```

If you're using the AWS CLI, don't forget to add `--profile <role-name>` to your AWS CLI calls (for example: `aws sts get-caller-identity --profile penguin-base-my-awesome-role`)

## Console
Because most people also use the Console, don't forget that you can always access your role in the console the following way:
```
https://signin.aws.amazon.com/switchrole?account=<your-account_id_number>&roleName=<role_name>
```
