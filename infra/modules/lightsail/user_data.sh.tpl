#!/bin/bash

sudo yum update -y

sudo yum install -y aws-cli

sudo tee -a /etc/profile.d/aws_creds.sh > /dev/null << 'EOF'
export AWS_REGION=${aws_region}
export AWS_DEFAULT_REGION=${aws_region}
EOF

source /etc/profile

aws s3 ls
