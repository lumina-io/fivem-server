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
    docker-compose up --build

    #echo ":: Attach console (exit: Ctrl+P, Ctrl+Q)"
    #docker start -ia "${CONTAINER_NAME}"
}

function _stop() {
    echo ":: Stopping container..."
    docker-compose down
}

if [ "$1" = "start" ]; then
    _start
elif [ "$1" = "stop" ]; then
    _stop
elif [ "$1" = "restart" ]; then
    _stop
    _start
else
    echo "Usage: launcher <start|stop|restart>"
    exit 1
fi