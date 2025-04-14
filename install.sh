#!/bin/bash

CONFIG_FILE=~/.geeksonator

check_dependencies() {
    if ! command -v docker &> /dev/null
    then
        echo "command 'docker' not found"
        exit
    fi
}

create_conf() {
    # Read user input
    read -p "Telegram bot token (Required): " BOT_TOKEN
    read -p "Timeout seconds (Default: 15): " TIMEOUT_SECONDS
    read -p "Debug mode (Default: false): " DEBUG_MODE
    read -p "Debug telegram bot token (Default: ""): " DEBUG_BOT_TOKEN

    # Set default values
    BOT_TOKEN=${BOT_TOKEN:-""}
    TIMEOUT_SECONDS=${TIMEOUT_SECONDS:-15}
    DEBUG_MODE=${DEBUG_MODE:-false}
    DEBUG_BOT_TOKEN=${DEBUG_BOT_TOKEN:-""}

    # Write configuration to CONFIG_FILE
    echo "GEEKSONATOR_TELEGRAM_BOT_TOKEN=$BOT_TOKEN" > $CONFIG_FILE
    echo "GEEKSONATOR_TELEGRAM_TIMEOUT_SECONDS=$TIMEOUT_SECONDS" >> $CONFIG_FILE
    echo "GEEKSONATOR_DEBUG_MODE=$DEBUG_MODE" >> $CONFIG_FILE
    echo "GEEKSONATOR_DEBUG_TELEGRAM_BOT_TOKEN=$DEBUG_BOT_TOKEN" >> $CONFIG_FILE
}

echo -e "\033[7m--]] Geeksonator Installer [[--\033[0m"

echo "Checking dependencies..."
check_dependencies
echo "Checking dependencies: OK"

echo "Creating configuration file ($CONFIG_FILE)..."
create_conf
echo "Creating configuration file: DONE"

echo "All done! Run next command to start the bot:"
echo -e "\033[7mdocker run -d --env-file $CONFIG_FILE --name geeksonator.app ghcr.io/phpgeeks-club/geeksonator:latest\033[0m"
