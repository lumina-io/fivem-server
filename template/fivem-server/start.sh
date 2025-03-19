#!/bin/bash
set -e
cd "$(dirname "$0")" && _base=$PWD

if [ -e "${SERVER_PATH}" ]; then
    cd ${SERVER_PATH}

    if [ -e "autorun.sh" ]; then
        echo ":: Running autorun.sh"
        bash ./autorun.sh
    fi
fi

# Start Server
cd ${_base}
export TXHOST_TXA_PORT=${TXADMIN_PORT:-40120}

if [ "$DIRECT" == "true" ]; then
    cd $DIRECT_DIR
    bash ${_base}/run.sh +exec server.cfg
else
    bash ./run.sh +set txAdminPort
fi
