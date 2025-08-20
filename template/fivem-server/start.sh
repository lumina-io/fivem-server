#!/bin/sh
set -e
cd "$(dirname "$0")" && _base=$PWD

if [ -e "${SERVER_PATH}" ]; then
    cd ${SERVER_PATH}

    if [ -e "autorun.sh" ]; then
        echo ":: Running autorun.sh"
        sh ./autorun.sh
    fi
fi

# Start Server
if [ "$DIRECT" == "true" ]; then
    cd ${SERVER_PATH}
    exec sh ${_base}/run.sh +exec server.cfg +set onesync on
else
    cd ${_base}
    exec sh ./run.sh
fi
