#!/bin/bash

ln -sf /usr/share/zoneinfo/Asia/Tokyo /etc/localtime

yum update -y

yum install -y ec2-instance-connect

# Configure SSH port
sed -i 's/^#Port 22/Port 22/' /etc/ssh/sshd_config
sed -i "/^Port 22/a Port ${ssh_port}" /etc/ssh/sshd_config
systemctl restart sshd

# Configure deploy key
mkdir -p /home/ec2-user/.ssh
chown ec2-user:ec2-user /home/ec2-user/.ssh
chmod 700 /home/ec2-user/.ssh
echo "${deploy_key}" > /home/ec2-user/.ssh/deploy_key
chown ec2-user:ec2-user /home/ec2-user/.ssh/deploy_key
chmod 600 /home/ec2-user/.ssh/deploy_key
cat <<EOF > /home/ec2-user/.ssh/config
Host github.com
  HostName github.com
  User git
  IdentityFile /home/ec2-user/.ssh/deploy_key
  StrictHostKeyChecking no
EOF
chown ec2-user:ec2-user /home/ec2-user/.ssh/config
chmod 600 /home/ec2-user/.ssh/config
