#!/bin/bash

app=ats
version=${1:-build_0}
commit=`git rev-parse HEAD`

repository='swr.cn-north-1.myhuaweicloud.com/onge'

docker build \
--build-arg ATS_REVISION=${commit} \
--build-arg ATS_VERSION=${version} \
--label org.opencontainers.image.revision=${commit} \
--label org.opencontainers.image.version=${version} \
--tag ${repository}/${app}:${version} .

docker push ${repository}/${app}:${version}
