# Assume Role

## Using in your Repo
```
git submodule add git@github.com:kelseymok/assume-role.git
git submodule update --init
```

## Console

```
https://signin.aws.amazon.com/switchrole?account=<your-account_id_number>&roleName=<role_name>
```

## CLI

Set up your AWS CLI credentials according to the [AWS CLI docs](https://docs.aws.amazon.com/cli/latest/userguide/cli-chap-configure.html). You should have the following files with the profile `base`:

```bash
# ~/.aws/config

[profile base]                                            
region = eu-central-1                                          
mfa_serial = arn:aws:iam::<your-aws-account-number>:mfa/<your-username>
```

```bash
# ~/.aws/credentials

[base]
aws_access_key_id = <some-access-key-id>
aws_secret_access_key = <some-secret-access-key>
```

I have provided a couple of Docker Images to help with assuming a role. First build the image: `./assume-role/ build`

There are two ways you can run the Docker Image:

1. As a tool to assume a role: `./assume-role assume-role <Role Name>`
2. As a tool to assume a role and as a development environment (comes with Terraform and a few base libraries): `./assume-role/go run` (and then `assume-role <Role Name>` once dropped into a bash session)

Both methods update your Shared Credentials file (`~/.aws/credentials`) with a profile named `<Role Name>`. 

If you use Terraform, make sure to add a profile to your provider configuration:
```hcl-terraform
provider "aws" {
  region = "eu-central-1"
  profile = "Your Role Name"
}
```

