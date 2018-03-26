#!/bin/bash -ex

scriptdir="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"

cd ${scriptdir}/../../../../vendor/k8s.io/code-generator && ./generate-groups.sh \
  all \
  github.com/tonya11en/mysql-operator/pkg/client \
  github.com/tonya11en/mysql-operator/pkg/apis \
  "myproject:v1alpha1" \
