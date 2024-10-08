# AWS SSO Credentials Getter

Tools like terraform are currently unable to handle AWS auth via SSO
and rely on the access key and secret being in the credentials file.

SSOcred is a cli tool that will grab temporary AWS credentials from IAM Identity Center 
and store them in ~/.aws/credentials

This is mostly a Golang fork of https://github.com/PredictMobile/aws-sso-credentials-getter with a few changes

Notably, if you do not provide a profile argument but have your AWS_PROFILE environment variable set,
it will use that profile

## Installation

Install Go from the [official website](https://go.dev/)

clone this repository and build the executable. Then move it to your bin folder

```bash
git clone https://github.com/DaraDadachanji/go-aws-sso-credentials-getter.git
cd go-aws-sso-credentials-getter
go mod tidy
go build
mv ./go-aws-sso-credentials-getter /usr/local/bin/ssocred
```

## Configuration

You should set up your profiles in your `~/.aws/config` file

```
[profile my-profile]
sso_account_id = 123456654321
sso_role_name = my_aws_role_name
sso_start_url = https://my-organization.awsapps.com/start
sso_region = us-east-1
region = us-east-1

[profile my-second-profile]
sso_account_id = 654321123456
sso_role_name = my_aws_role_name
sso_start_url = https://my-organization.awsapps.com/start
sso_region = us-east-1
region = eu-west-2
```

The key should be the profile name generated by SSO while the value should be
your preferred alias for it.

This step is optional

## Usage

first login to identity center using `aws sso login --profile {{profile-alias}}`
then run `ssocred {profile-alias}` in your terminal

note that you only need to log into identity center once on any profile, 
then you can run sso-cred on as many profiles as you want