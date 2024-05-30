#!/bin/bash

cd lambda

#aws ecr-public get-login-password --region eu-west-2 | docker login --username AWS --password-stdin public.ecr.aws/s3l6t2t7
#aws ecr-public get-login-password --region us-east-1 | docker login --username AWS --password-stdin public.ecr.aws

aws ecr get-login-password --region eu-west-2 | docker login --username AWS --password-stdin 109587176248.dkr.ecr.eu-west-2.amazonaws.com/ftaplatform

sam build
sam package --output-template-file packaged.yaml --image-repository 109587176248.dkr.ecr.eu-west-2.amazonaws.com/ftaplatform
sam deploy