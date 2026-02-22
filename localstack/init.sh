#!/bin/bash

awslocal s3api create-bucket \
  --bucket sampay-private \
  --create-bucket-configuration LocationConstraint="$AWS_DEFAULT_REGION"

awslocal s3api create-bucket \
  --bucket sampay-public \
  --create-bucket-configuration LocationConstraint="$AWS_DEFAULT_REGION"

awslocal sqs create-queue \
  --queue-name sampay-worker \
  --region ap-northeast-1

awslocal sqs create-queue \
  --queue-name sampay-worker-dlq \
  --region ap-northeast-1
