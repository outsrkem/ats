#!/bin/sh
set -e

if [ -z "$1" ];then
  set -- ats -c /usr/local/bin/ats.yaml
fi

exec "$@"
