#!/bin/bash

cd "$(dirname "$0")"
cd ..

source config.sh
docker compose up -d mariadb
docker exec -it fivem-mariadb mariadb -u root -p${MYSQL_ROOT_PASSWORD} -h ${MYSQL_ADDRESS} -P ${MYSQL_PORT} -D ${MYSQL_DATABASE}
