#!/bin/bash
set -e
cd "$(dirname "$0")" && _base=$PWD

if [ -e "${SERVER_PATH}" ]; then
    cd ${SERVER_PATH}

    if [ -e "autorun.sh" ]; then
        echo ":: Running autorun.sh"
        kontra bash ./autorun.sh
    fi
fi

# Start Server
export TXHOST_TXA_PORT=${TXADMIN_PORT:-40120}

if [ "$DIRECT" == "true" ]; then
    cd ${SERVER_PATH}
    kontra bash ${_base}/run.sh +exec server.cfg +set onesync on
else
    cd ${_base}
    kontra bash ./run.sh
fi
