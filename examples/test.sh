#!/bin/bash

rm -f .terraform.lock.hcl
rm -rf .terraform
#rm -f terraform.tfstate
(cd .. && make install)
terraform init
terraform validate
