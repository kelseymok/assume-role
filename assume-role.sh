#!/usr/bin/env bash

set -e

base_profile=$1
if [ -z "${base_profile}" ]; then
  echo "Base Profile not set. Usage <BASE_PROFILE> <ROLE>"
  exit 1
fi

role=$2
if [ -z "${role}" ]; then
  echo "Role not set. Usage: <BASE_PROFILE> <ROLE>"
  exit 1
fi

display_role=$(echo $role | tr / -)

mfa_serial=$(aws configure get mfa_serial --profile ${base_profile})
if [ -z "${mfa_serial}" ]; then
  echo "mfa_serial for profile config does not exist. Did you set it?"
  exit 1
fi

mfa_serial=$(aws configure get mfa_serial --profile ${base_profile})
read -p "Enter MFA (${mfa_serial}): " mfa_response
if [ -z "${mfa_response}" ]; then
  echo "MFA not provided"
  exit 1
fi

account=$(aws sts get-caller-identity --profile ${base_profile} | jq -r '.Account')

response=$(aws sts assume-role \
--role-arn "arn:aws:iam::${account}:role/${role}" \
--role-session-name session \
--profile ${base_profile} \
--serial-number "${mfa_serial}" \
--token-code "${mfa_response}")

aws configure set aws_access_key_id $(echo $response | jq -r '.Credentials.AccessKeyId') --profile "${display_role}"
aws configure set aws_secret_access_key $(echo $response | jq -r '.Credentials.SecretAccessKey') --profile "${display_role}"
aws configure set aws_session_token $(echo $response | jq -r '.Credentials.SessionToken') --profile "${display_role}"

