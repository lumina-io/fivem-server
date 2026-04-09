#!/bin/bash
set -e
cd "$(dirname "$0")"

source config.env

# Set container use
export USER_ID=$(id -u)
export GROUP_ID=$(id -g)

function compose() {
    docker compose --env-file ./config.env "$@"
}

function _database_console() {
    query="$*"

    if [ "$query" = "" ]; then
        compose run --rm -it mariadb /usr/bin/mariadb -u root -p${MARIADB_ROOT_PASSWORD} -P ${MARIADB_PORT}
    else
        compose run --rm -T mariadb /usr/bin/mariadb -u root -p${MARIADB_ROOT_PASSWORD} -P ${MARIADB_PORT} -v -e "$*"
    fi
}

function _preexec() {
    if [ -e "$PWD/preexec.sh" ]; then
        bash $PWD/preexec.sh
    fi
}

function _start() {
    mkdir -p $PWD/local

    _preexec

    echo ":: Starting container..."
    if [ "$USE_DEV" = "true" ]; then
        echo ":: USE_DEV enabled."
        compose -f docker-compose.yaml -f docker-compose.dev.yaml up --build
    else
        compose up --build
    fi

    #echo ":: Attach console (exit: Ctrl+P, Ctrl+Q)"
    #docker start -ia "${CONTAINER_NAME}"
}

function _stop() {
    echo ":: Stopping container..."
    compose down
}

function _build() {
    echo ":: Building image..."
    compose build --no-cache
}

function _config() {
    compose config
}

function _start_database() {
    mkdir -p $PWD/local/mysql

    _preexec

    echo ":: Starting container..."
    if [ "$USE_DEV" = "true" ]; then
        echo ":: USE_DEV enabled."
        compose -f docker-compose.yaml -f docker-compose.dev.yaml up mariadb --build
    else
        compose up mariadb --build
    fi

    #echo ":: Attach console (exit: Ctrl+P, Ctrl+Q)"
    #docker start -ia "${CONTAINER_NAME}"
}

if [ "$1" = "start" ]; then
    _start
elif [ "$1" = "start-database" ]; then
    _start_database
elif [ "$1" = "database-console" ]; then
    shift 1
    _database_console $@
elif [ "$1" = "stop" ]; then
    _stop
elif [ "$1" = "restart" ]; then
    _stop
    _start
elif [ "$1" = "build" ]; then
    _build
elif [ "$1" = "config" ]; then
    _config
elif [ "$1" = "compose" ]; then
    shift 1
    compose $@
else
    echo "Usage: launcher <start|start-database|database-console|stop|restart|build|config|compose>"
    exit 1
fi
