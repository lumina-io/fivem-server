#!/bin/bash
set -e
cd "$(dirname "$0")"

source config.env
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
    docker compose up -d --build
}

function _stop() {
    echo ":: Stopping container..."
    docker compose down
}

if [ "$1" = "start" ]; then
    _start
elif [ "$1" = "stop" ]; then
    _stop
else
    echo "Usage: launcher <start|stop>"
    exit 1
fi
