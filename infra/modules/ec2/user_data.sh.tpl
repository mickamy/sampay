#!/bin/bash

sudo yum update -y

sudo yum install -y ec2-instance-connect

# Configure deploy key
echo {deploy_key} > /home/ec2-user/.ssh/deploy_key
chown ec2-user:ec2-user /home/ec2-user/.ssh/deploy_key
chmod 600 /home/ec2-user/.ssh/deploy_key
cat <<EOF > /home/ec2-user/.ssh/config
Host github.com
  HostName github.com
  User git
  IdentityFile /home/ec2-user/.ssh/deploy_key
EOF
chown ec2-user:ec2-user /home/ec2-user/.ssh/config
chmod 600 /home/ec2-user/.ssh/config
