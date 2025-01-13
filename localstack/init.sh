#!/bin/bash

awslocal s3api create-bucket --bucket sampay-public --create-bucket-configuration LocationConstraint=ap-northeast-1
