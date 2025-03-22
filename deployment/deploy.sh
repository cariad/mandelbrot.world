#!/bin/bash

set -euo pipefail

repository_template_path=deployment/repository.cf.yml
repository_stack_name=mandelbrotworld-repository

service_template_path=deployment/service.cf.yml
service_stack_name=mandelbrotworld-service

echo "💬 Validating ${repository_template_path}…"
aws cloudformation validate-template --template-body "file://${repository_template_path}" > /dev/null

echo "💬 Validating ${service_template_path}…"
aws cloudformation validate-template --template-body "file://${service_template_path}" > /dev/null

echo "💬 Deploying ${repository_stack_name}…"
aws cloudformation deploy --template-file "${repository_template_path}" --stack-name "${repository_stack_name}"

echo "💬 Discovering repository…"
repository_uri=$(
	aws cloudformation describe-stacks \
		--output text \
		--query "Stacks[0].Outputs[?OutputKey=='RepositoryUri'].OutputValue" \
		--stack-name "${repository_stack_name}"
)

repository_server="${repository_uri%%/*}"

echo "💬 Building ${repository_uri}…"
docker build . -t "${repository_uri}"

image=$(docker inspect --format='{{index .RepoDigests 0}}' "${repository_uri}")

echo "💬 Logging into ${repository_server}…"
aws ecr get-login-password | docker login --username AWS --password-stdin "${repository_server}"

echo "💬 Pushing ${repository_uri}…"
docker push "${repository_uri}"

echo "💬 Deploying ${service_stack_name} with ${image}…"
aws cloudformation deploy \
	--capabilities        CAPABILITY_IAM \
	--parameter-overrides Image="${image}" \
	--stack-name          "${service_stack_name}" \
	--template-file       "${service_template_path}"
