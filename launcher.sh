#!/bin/bash
set -e
cd "$(dirname "$0")"

source $PWD/config.sh

# Set container user
export USER_ID=`id -u`
export GROUP_ID=`id -g`

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
        docker compose -f docker-compose.yaml -f docker-compose.dev.yaml up --build
    else
        docker compose up --build
    fi

    #echo ":: Attach console (exit: Ctrl+P, Ctrl+Q)"
    #docker start -ia "${CONTAINER_NAME}"
}

function _stop() {
    echo ":: Stopping container..."
    docker compose down
}

function _build() {
    echo ":: Building image..."
    docker compose build --no-cache
}

function _start_database() {
    mkdir -p $PWD/local

    _preexec

    echo ":: Starting container..."
    if [ "$USE_DEV" = "true" ]; then
        echo ":: USE_DEV enabled."
        docker compose -f docker-compose.yaml -f docker-compose.dev.yaml up mariadb --build
    else
        docker compose up mariadb --build
    fi

    #echo ":: Attach console (exit: Ctrl+P, Ctrl+Q)"
    #docker start -ia "${CONTAINER_NAME}"
}


if [ "$1" = "start" ]; then
    _start
elif [ "$1" = "start-database" ]; then
    _start_database
elif [ "$1" = "stop" ]; then
    _stop
elif [ "$1" = "restart" ]; then
    _stop
    _start
elif [ "$1" = "build" ]; then
    _build
else
    echo "Usage: launcher <start|start-database|stop|restart|build>"
    exit 1
fi
