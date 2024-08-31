#!/bin/bash
set -e
cd "$(dirname "$0")" && _base=$PWD

rewriteGlobal() {
    local TARGET="$1"

    echo "=> Rewrite parameters..."

    for f in $(find . -mindepth 1 -maxdepth 2 -type f \
        -and -name '*.yaml' \
        -or -name '*.yml' \
        -or -name '*.json' \
        -or -name '*.txt' \
        -or -name '*.ini' \
        -or -name '*.cfg' \
        -or -name '*.conf' \
        -or -name '*.config' \
        -or -name '*.properties' \
    ); do
        sed -i $f -e "s/!!SERVER_HOSTNAME!!/${SERVER_HOSTNAME}/g"
        sed -i $f -e "s/!!SERVER_PROJECT_NAME!!/${SERVER_PROJECT_NAME}/g"
        sed -i $f -e "s/!!SERVER_PROJECT_DESC!!/${SERVER_PROJECT_DESC}/g"
    done
}


rewriteGlobal /app/fivem/txData

# Start Server
cd ${_base}
bash ./run.sh +set txAdminPort ${TXADMIN_PORT:-40120}
