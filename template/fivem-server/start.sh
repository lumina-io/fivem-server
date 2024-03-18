#!/bin/bash
set -e
cd "$(dirname "$0")" && _base=$PWD

cd ${_base}
bash ./run.sh +set txAdminPort ${TXADMIN_PORT:-40120}
